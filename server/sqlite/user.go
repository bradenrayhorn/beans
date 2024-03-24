package sqlite

import (
	"context"

	"github.com/bradenrayhorn/beans/server/beans"
	"zombiezen.com/go/sqlite"
)

type userRepository struct{ repository }

var _ beans.UserRepository = (*userRepository)(nil)

const userCreateSQL = `
INSERT INTO users
	(id, username, password)
	VALUES (:id,:username,:password)
`

func (r *userRepository) Create(ctx context.Context, id beans.ID, username beans.Username, passwordHash beans.PasswordHash) error {
	return db[any](r.pool).execute(ctx, userCreateSQL, map[string]any{
		":id":       id.String(),
		":username": string(username),
		":password": string(passwordHash),
	})
}

const userExistsSQL = `
SELECT EXISTS (SELECT id FROM users WHERE username = :username)
`

func (r *userRepository) Exists(ctx context.Context, username beans.Username) (bool, error) {
	return db[bool](r.pool).
		mapWith(func(stmt *sqlite.Stmt) (bool, error) { return stmt.ColumnBool(0), nil }).
		one(ctx, userExistsSQL, map[string]any{
			":username": string(username),
		})
}

const userGetOneSQL = `
SELECT * FROM users
	WHERE id = :id
`

func (r *userRepository) Get(ctx context.Context, id beans.ID) (beans.User, error) {
	return db[beans.User](r.pool).
		mapWith(mapUser).
		one(ctx, userGetOneSQL, map[string]any{
			":id": id.String(),
		})
}

const userGetByUsernameSQL = `
SELECT * FROM users
	WHERE username = :username
`

func (r *userRepository) GetByUsername(ctx context.Context, username beans.Username) (beans.User, error) {
	return db[beans.User](r.pool).
		mapWith(mapUser).
		one(ctx, userGetByUsernameSQL, map[string]any{
			":username": string(username),
		})
}

// mappers

func mapUser(stmt *sqlite.Stmt) (beans.User, error) {
	id, err := mapID(stmt, "id")
	if err != nil {
		return beans.User{}, err
	}

	return beans.User{
		ID:           id,
		Username:     beans.Username(stmt.GetText("username")),
		PasswordHash: beans.PasswordHash(stmt.GetText("password")),
	}, nil
}
