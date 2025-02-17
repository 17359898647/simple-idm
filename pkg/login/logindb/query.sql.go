// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package logindb

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const findUser = `-- name: FindUser :one




SELECT uuid, username, password
FROM login
WHERE username = $1
AND deleted_at IS NULL
`

type FindUserRow struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Username sql.NullString `json:"username"`
	Password []byte         `json:"password"`
}

// -- name: FindUsers :many
// SELECT uuid, created_at, last_modified_at, deleted_at, created_by, email, name
// FROM users
// limit 20;
// -- name: RegisterUser :one
// INSERT INTO users (email, name, password, created_at)
// VALUES ($1, $2, $3, NOW())
// RETURNING *;
// -- name: EmailVerify :exec
// UPDATE users
// SET verified_at = NOW()
// WHERE email = $1;
// -- name: InitPassword :one
// SELECT uuid
// FROM users
// WHERE email = $1;
func (q *Queries) FindUser(ctx context.Context, username sql.NullString) (FindUserRow, error) {
	row := q.db.QueryRow(ctx, findUser, username)
	var i FindUserRow
	err := row.Scan(&i.Uuid, &i.Username, &i.Password)
	return i, err
}

const findUserByUsername = `-- name: FindUserByUsername :many
SELECT u.uuid, u.username, u.password, u.email, u.name, u.created_at, u.last_modified_at,
       array_agg(r.name) as roles
FROM users u
LEFT JOIN user_roles ur ON u.uuid = ur.user_uuid
LEFT JOIN roles r ON ur.role_uuid = r.uuid
WHERE u.username = $1
GROUP BY u.uuid, u.username, u.password, u.email, u.name, u.created_at, u.last_modified_at
`

type FindUserByUsernameRow struct {
	Uuid           uuid.UUID      `json:"uuid"`
	Username       sql.NullString `json:"username"`
	Password       []byte         `json:"password"`
	Email          string         `json:"email"`
	Name           sql.NullString `json:"name"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
	Roles          interface{}    `json:"roles"`
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
			&i.Username,
			&i.Password,
			&i.Email,
			&i.Name,
			&i.CreatedAt,
			&i.LastModifiedAt,
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

const findUserInfoWithRoles = `-- name: FindUserInfoWithRoles :one
SELECT u.email, u.username, u.name, COALESCE(array_agg(r.name), '{}') AS roles
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

func (q *Queries) FindUserInfoWithRoles(ctx context.Context, argUuid uuid.UUID) (FindUserInfoWithRolesRow, error) {
	row := q.db.QueryRow(ctx, findUserInfoWithRoles, argUuid)
	var i FindUserInfoWithRolesRow
	err := row.Scan(
		&i.Email,
		&i.Username,
		&i.Name,
		&i.Roles,
	)
	return i, err
}

const findUserRolesByUserUuid = `-- name: FindUserRolesByUserUuid :many
SELECT name
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
		var name sql.NullString
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findUsernameByEmail = `-- name: FindUsernameByEmail :one
SELECT username
FROM users
WHERE email = $1
`

func (q *Queries) FindUsernameByEmail(ctx context.Context, email string) (sql.NullString, error) {
	row := q.db.QueryRow(ctx, findUsernameByEmail, email)
	var username sql.NullString
	err := row.Scan(&username)
	return username, err
}

const initPasswordByUsername = `-- name: InitPasswordByUsername :one
SELECT uuid
FROM users
WHERE username = $1
`

func (q *Queries) InitPasswordByUsername(ctx context.Context, username sql.NullString) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, initPasswordByUsername, username)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const initPasswordResetToken = `-- name: InitPasswordResetToken :exec
INSERT INTO password_reset_tokens (user_uuid, token, expire_at)
VALUES ($1, $2, $3)
`

type InitPasswordResetTokenParams struct {
	UserUuid uuid.UUID          `json:"user_uuid"`
	Token    string             `json:"token"`
	ExpireAt pgtype.Timestamptz `json:"expire_at"`
}

func (q *Queries) InitPasswordResetToken(ctx context.Context, arg InitPasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, initPasswordResetToken, arg.UserUuid, arg.Token, arg.ExpireAt)
	return err
}

const markPasswordResetTokenUsed = `-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens
SET used_at = NOW()
WHERE token = $1
`

func (q *Queries) MarkPasswordResetTokenUsed(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, markPasswordResetTokenUsed, token)
	return err
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

const resetPasswordByUuid = `-- name: ResetPasswordByUuid :exec
UPDATE users
SET password = $1,
    last_modified_at = NOW()
WHERE uuid = $2
`

type ResetPasswordByUuidParams struct {
	Password []byte    `json:"password"`
	Uuid     uuid.UUID `json:"uuid"`
}

func (q *Queries) ResetPasswordByUuid(ctx context.Context, arg ResetPasswordByUuidParams) error {
	_, err := q.db.Exec(ctx, resetPasswordByUuid, arg.Password, arg.Uuid)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1,
    last_modified_at = NOW()
WHERE uuid = $2
`

type UpdateUserPasswordParams struct {
	Password []byte    `json:"password"`
	Uuid     uuid.UUID `json:"uuid"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Password, arg.Uuid)
	return err
}

const validatePasswordResetToken = `-- name: ValidatePasswordResetToken :one
SELECT prt.uuid as uuid, prt.user_uuid as user_uuid, u.email as email
FROM password_reset_tokens prt
JOIN users u ON u.uuid = prt.user_uuid
WHERE prt.token = $1
  AND prt.expire_at > NOW()
  AND prt.used_at IS NULL
LIMIT 1
`

type ValidatePasswordResetTokenRow struct {
	Uuid     uuid.UUID `json:"uuid"`
	UserUuid uuid.UUID `json:"user_uuid"`
	Email    string    `json:"email"`
}

func (q *Queries) ValidatePasswordResetToken(ctx context.Context, token string) (ValidatePasswordResetTokenRow, error) {
	row := q.db.QueryRow(ctx, validatePasswordResetToken, token)
	var i ValidatePasswordResetTokenRow
	err := row.Scan(&i.Uuid, &i.UserUuid, &i.Email)
	return i, err
}
