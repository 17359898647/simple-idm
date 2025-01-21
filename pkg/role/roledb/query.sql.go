// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package roledb

import (
	"context"

	"github.com/google/uuid"
)

const createRole = `-- name: CreateRole :one
INSERT INTO roles (name) VALUES ($1) RETURNING uuid
`

func (q *Queries) CreateRole(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createRole, name)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM roles WHERE uuid = $1
`

func (q *Queries) DeleteRole(ctx context.Context, argUuid uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRole, argUuid)
	return err
}

const findRoles = `-- name: FindRoles :many
SELECT uuid, name
FROM roles
ORDER BY name ASC
`

type FindRolesRow struct {
	Uuid uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (q *Queries) FindRoles(ctx context.Context) ([]FindRolesRow, error) {
	rows, err := q.db.Query(ctx, findRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindRolesRow
	for rows.Next() {
		var i FindRolesRow
		if err := rows.Scan(&i.Uuid, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleUUID = `-- name: GetRoleUUID :one
SELECT uuid, name
FROM roles
WHERE uuid = $1
`

type GetRoleUUIDRow struct {
	Uuid uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (q *Queries) GetRoleUUID(ctx context.Context, argUuid uuid.UUID) (GetRoleUUIDRow, error) {
	row := q.db.QueryRow(ctx, getRoleUUID, argUuid)
	var i GetRoleUUIDRow
	err := row.Scan(&i.Uuid, &i.Name)
	return i, err
}

const hasUsers = `-- name: HasUsers :one
SELECT EXISTS (
    SELECT 1 FROM user_roles WHERE role_uuid = $1
) as has_users
`

func (q *Queries) HasUsers(ctx context.Context, roleUuid uuid.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, hasUsers, roleUuid)
	var has_users bool
	err := row.Scan(&has_users)
	return has_users, err
}

const updateRole = `-- name: UpdateRole :exec
UPDATE roles SET name = $2 WHERE uuid = $1
`

type UpdateRoleParams struct {
	Uuid uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) error {
	_, err := q.db.Exec(ctx, updateRole, arg.Uuid, arg.Name)
	return err
}
