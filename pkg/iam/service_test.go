package iam

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendant/simple-idm/pkg/iam/iamdb"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func containerLog(ctx context.Context, container testcontainers.Container) {
	// Retrieve logs
	logs, err := container.Logs(ctx)
	if err != nil {
		slog.Error("Failed to get container logs:", "err", err)
	}
	defer logs.Close()

	// Process and display logs
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Print each log line
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		slog.Error("Error reading logs", "err", err)
	}
}

func setupTestDatabase(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()

	// Create PostgreSQL container
	dbName := "idm_db"
	dbUser := "idm"
	dbPassword := "pwd"

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("../../migrations", "idm_db.sql")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	if err != nil {
		slog.Error("Failed to start container:", "err", err)
	}

	// containerLog(ctx, container)

	// Generate the connection string
	connString, err := container.ConnectionString(ctx)
	fmt.Println("Connection string:", connString)
	require.NoError(t, err)

	// Create connection pool
	poolConfig, err := pgxpool.ParseConfig(connString)
	require.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return pool, cleanup
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	pool, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create test dependencies
	queries := iamdb.New(pool)
	service := NewIamServiceWithQueries(queries)

	// Create a test role first
	role, err := queries.CreateRole(ctx, "test-role")
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name      string
		email     string
		username  string
		fullName  string
		roleUuids []uuid.UUID
		wantErr   bool
	}{
		{
			name:      "valid user with role",
			email:     "test@example.com",
			username:  "testuser",
			fullName:  "Test User",
			roleUuids: []uuid.UUID{role},
			wantErr:   false,
		},
		{
			name:      "valid user without role",
			email:     "test2@example.com",
			username:  "testuser2",
			fullName:  "Test User 2",
			roleUuids: []uuid.UUID{},
			wantErr:   false,
		},
		{
			name:     "missing email",
			username: "testuser3",
			fullName: "Test User 3",
			wantErr:  true,
		},
		{
			name:     "missing username",
			email:    "test3@example.com",
			fullName: "Test User 3",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.CreateUser(ctx, tt.email, tt.username, tt.fullName, tt.roleUuids, "")
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tt.email, user.Email)
			// Username is populated from login, which is not set in this test
			// assert.Equal(t, tt.username, user.Username)
			if tt.fullName != "" {
				assert.Equal(t, tt.fullName, user.Name)
			}

			// Verify roles if any were assigned
			if len(tt.roleUuids) > 0 {

				assert.Len(t, user.Roles, len(tt.roleUuids))
				if len(user.Roles) > 0 {
					assert.Equal(t, tt.roleUuids[0], user.Roles[0].ID)
					assert.Equal(t, "test-role", user.Roles[0].Name)
				}
			} else {
				// For users without roles, the roles array should be empty
				assert.Empty(t, user.Roles, "Expected empty roles array")
			}
		})
	}
}

func TestFindUsers(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	pool, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create test dependencies
	queries := iamdb.New(pool)
	service := NewIamServiceWithQueries(queries)

	// Create a test role
	role, err := queries.CreateRole(ctx, "test-role")
	require.NoError(t, err)

	// Create test users
	testUsers := []struct {
		email     string
		username  string
		name      string
		roleUuids []uuid.UUID
	}{
		{
			email:     "test1@example.com",
			username:  "testuser1",
			name:      "Test User 1",
			roleUuids: []uuid.UUID{role},
		},
		{
			email:     "test2@example.com",
			username:  "testuser2",
			name:      "Test User 2",
			roleUuids: []uuid.UUID{},
		},
	}

	for _, u := range testUsers {
		_, err := service.CreateUser(ctx, u.email, u.username, u.name, u.roleUuids, "")
		require.NoError(t, err)
	}

	// Test finding users
	users, err := service.FindUsers(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, len(testUsers))

	// Verify each user
	for i, u := range users {
		assert.Equal(t, testUsers[i].email, u.Email)
		// Username is populated from login, which is not set in this test
		// assert.Equal(t, testUsers[i].username, u.Username)
		assert.Equal(t, testUsers[i].name, u.Name)

		// Verify roles
		if len(testUsers[i].roleUuids) > 0 {
			assert.Len(t, u.Roles, len(testUsers[i].roleUuids))

			// Verify each role UUID
			roleIDs := make(map[uuid.UUID]bool)
			for _, role := range u.Roles {
				roleIDs[role.ID] = true
			}
			for _, expectedUUID := range testUsers[i].roleUuids {
				assert.True(t, roleIDs[expectedUUID], "Expected role UUID not found: %s", expectedUUID)
			}
		} else {
			// For users without roles, the roles array should be empty
			assert.Empty(t, u.Roles, "Expected empty roles array")
		}
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	pool, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create test dependencies
	queries := iamdb.New(pool)
	service := NewIamServiceWithQueries(queries)

	// Create a test role
	roleUUID := uuid.New()
	_, err := pool.Exec(ctx, `
		INSERT INTO roles (id, name, description)
		VALUES ($1, $2, $3)
	`, roleUUID, "TestRole", "A test role")
	require.NoError(t, err)

	// Create a test user
	user, err := service.CreateUser(ctx, "test@example.com", "testuser", "Test User", []uuid.UUID{roleUUID}, "")
	require.NoError(t, err)

	// Test cases
	t.Run("existing user", func(t *testing.T) {
		// Get the user
		fetchedUser, err := service.GetUser(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, fetchedUser.Email)
		assert.Equal(t, user.Name, fetchedUser.Name)
	})

	t.Run("non-existent user", func(t *testing.T) {
		// Try to get a non-existent user
		_, err := service.GetUser(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	pool, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create test dependencies
	queries := iamdb.New(pool)
	service := NewIamServiceWithQueries(queries)

	// Create test roles
	role1, err := queries.CreateRole(ctx, "role-1")
	require.NoError(t, err)

	role2, err := queries.CreateRole(ctx, "role-2")
	require.NoError(t, err)

	require.NoError(t, err)

	// Create initial user with role1
	initialUser, err := service.CreateUser(ctx, "test@example.com", "testuser", "Test User", []uuid.UUID{role1}, "")
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name      string
		newName   string
		roleUuids []uuid.UUID
		wantErr   bool
	}{
		{
			name:      "update name and roles",
			newName:   "Updated User",
			roleUuids: []uuid.UUID{role2},
			wantErr:   false,
		},
		{
			name:      "update to multiple roles",
			newName:   "Multi Role User",
			roleUuids: []uuid.UUID{role1, role2},
			wantErr:   false,
		},
		{
			name:      "remove all roles",
			newName:   "No Role User",
			roleUuids: []uuid.UUID{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update user
			updatedUser, err := service.UpdateUser(ctx, initialUser.ID, tt.newName, tt.roleUuids, nil)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, updatedUser)
			assert.Equal(t, tt.newName, updatedUser.Name)

			// Verify roles
			if len(tt.roleUuids) > 0 {
				assert.Len(t, updatedUser.Roles, len(tt.roleUuids))
				roleIDs := make(map[uuid.UUID]bool)
				for _, role := range updatedUser.Roles {
					roleIDs[role.ID] = true
				}
				for _, expectedUUID := range tt.roleUuids {
					assert.True(t, roleIDs[expectedUUID], "Expected role UUID not found: %s", expectedUUID)
				}
			} else {
				// When there are no roles, the array should be empty
				assert.Empty(t, updatedUser.Roles, "Expected empty roles array")
			}

			// Verify the roles in database
			var roleCount int
			err = pool.QueryRow(ctx, `
				SELECT COUNT(*)
				FROM user_roles
				WHERE user_id = $1
			`, updatedUser.ID).Scan(&roleCount)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.roleUuids), roleCount)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	// Setup test database
	pool, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Create test dependencies
	queries := iamdb.New(pool)
	service := NewIamServiceWithQueries(queries)

	// Create a test role
	roleUUID := uuid.New()
	_, err := pool.Exec(ctx, `
		INSERT INTO roles (id, name, description)
		VALUES ($1, $2, $3)
	`, roleUUID, "TestRole", "A test role")
	require.NoError(t, err)

	// Create a test user
	user, err := service.CreateUser(ctx, "test@example.com", "testuser", "Test User", []uuid.UUID{roleUUID}, "")
	require.NoError(t, err)

	// Test cases
	t.Run("existing user", func(t *testing.T) {
		// Delete the user
		err := service.DeleteUser(ctx, user.ID)
		require.NoError(t, err)

		// Verify user is marked as deleted
		var deletedAt sql.NullTime
		err = pool.QueryRow(ctx, `
			SELECT deleted_at
			FROM users
			WHERE id = $1
		`, user.ID).Scan(&deletedAt)
		assert.NoError(t, err)
		assert.True(t, deletedAt.Valid)
	})

	t.Run("non-existent user", func(t *testing.T) {
		// Try to delete a non-existent user
		err := service.DeleteUser(ctx, uuid.New())
		assert.Error(t, err)
	})
}
