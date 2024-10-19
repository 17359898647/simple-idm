// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const emailVerify = `-- name: EmailVerify :exec
UPDATE users
SET verified_at = NOW()
WHERE email = $1
`

func (q *Queries) EmailVerify(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, emailVerify, email)
	return err
}

const findUser = `-- name: FindUser :one
SELECT uuid, name, email, password
FROM users
WHERE email = $1
`

type FindUserRow struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Name     sql.NullString `json:"name"`
	Email    string         `json:"email"`
	Password []byte         `json:"password"`
}

func (q *Queries) FindUser(ctx context.Context, email string) (FindUserRow, error) {
	row := q.db.QueryRow(ctx, findUser, email)
	var i FindUserRow
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :many
SELECT users.uuid, name, username, email, password
FROM users
WHERE username = $1
`

type FindUserByUsernameRow struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Name     sql.NullString `json:"name"`
	Username sql.NullString `json:"username"`
	Email    string         `json:"email"`
	Password []byte         `json:"password"`
}

func (q *Queries) FindUserByUsername(ctx context.Context, username sql.NullString) ([]FindUserByUsernameRow, error) {
	rows, err := q.db.Query(ctx, findUserByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUserByUsernameRow
	for rows.Next() {
		var i FindUserByUsernameRow
		if err := rows.Scan(
			&i.Uuid,
			&i.Name,
			&i.Username,
			&i.Email,
			&i.Password,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUserInfoWithRoles = `-- name: FindUserInfoWithRoles :many
SELECT u.email, u.username, u.name, COALESCE(array_agg(r.role_name), '{}') AS roles
FROM public.users u
LEFT JOIN public.user_roles ur ON u.uuid = ur.user_uuid
LEFT JOIN public.roles r ON ur.role_uuid = r.uuid
WHERE u.uuid = $1
GROUP BY u.email, u.username, u.name
`

type FindUserInfoWithRolesRow struct {
	Email    string         `json:"email"`
	Username sql.NullString `json:"username"`
	Name     sql.NullString `json:"name"`
	Roles    interface{}    `json:"roles"`
}

func (q *Queries) FindUserInfoWithRoles(ctx context.Context, argUuid uuid.UUID) ([]FindUserInfoWithRolesRow, error) {
	rows, err := q.db.Query(ctx, findUserInfoWithRoles, argUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUserInfoWithRolesRow
	for rows.Next() {
		var i FindUserInfoWithRolesRow
		if err := rows.Scan(
			&i.Email,
			&i.Username,
			&i.Name,
			&i.Roles,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUserRolesByUserUuid = `-- name: FindUserRolesByUserUuid :many
SELECT role_name
FROM user_roles ur
LEFT JOIN roles ON ur.role_uuid = roles.uuid
WHERE ur.user_uuid = $1
`

func (q *Queries) FindUserRolesByUserUuid(ctx context.Context, userUuid uuid.UUID) ([]sql.NullString, error) {
	rows, err := q.db.Query(ctx, findUserRolesByUserUuid, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var role_name sql.NullString
		if err := rows.Scan(&role_name); err != nil {
			return nil, err
		}
		items = append(items, role_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUsers = `-- name: FindUsers :many
SELECT uuid, created_at, last_modified_at, deleted_at, created_by, email, name
FROM users
limit 20
`

type FindUsersRow struct {
	Uuid           uuid.UUID      `json:"uuid"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	CreatedBy      sql.NullString `json:"created_by"`
	Email          string         `json:"email"`
	Name           sql.NullString `json:"name"`
}

func (q *Queries) FindUsers(ctx context.Context) ([]FindUsersRow, error) {
	rows, err := q.db.Query(ctx, findUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUsersRow
	for rows.Next() {
		var i FindUsersRow
		if err := rows.Scan(
			&i.Uuid,
			&i.CreatedAt,
			&i.LastModifiedAt,
			&i.DeletedAt,
			&i.CreatedBy,
			&i.Email,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const initPassword = `-- name: InitPassword :one
SELECT uuid
FROM users
WHERE email = $1
`

func (q *Queries) InitPassword(ctx context.Context, email string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, initPassword, email)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const registerUser = `-- name: RegisterUser :one
INSERT INTO users (email, name, password, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING uuid, created_at, last_modified_at, deleted_at, created_by, email, name, password, verified_at, username
`

type RegisterUserParams struct {
	Email    string         `json:"email"`
	Name     sql.NullString `json:"name"`
	Password []byte         `json:"password"`
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (User, error) {
	row := q.db.QueryRow(ctx, registerUser, arg.Email, arg.Name, arg.Password)
	var i User
	err := row.Scan(
		&i.Uuid,
		&i.CreatedAt,
		&i.LastModifiedAt,
		&i.DeletedAt,
		&i.CreatedBy,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.VerifiedAt,
		&i.Username,
	)
	return i, err
}

const resetPassword = `-- name: ResetPassword :exec
UPDATE users
SET password = $1, 
    last_modified_at = NOW()
WHERE email = $2
`

type ResetPasswordParams struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

func (q *Queries) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	_, err := q.db.Exec(ctx, resetPassword, arg.Password, arg.Email)
	return err
}
