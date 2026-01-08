package store

import (
	"backend-go/internals/tokens"
	"database/sql"
	"time"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteALLTokenForUser(userId int, scope string) error
}

func (ts *PostgresTokenStore) CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO tokens (user_id, hash, expiry, scope) 
		VALUES ($1, $2, $3, $4)
	`

	_, err = ts.db.Exec(query, token.UserID, token.Hash, token.Expiry, token.Scope)
	return token, err
}

func (ts *PostgresTokenStore) DeleteALLTokenForUser(userId int, scope string) error {
	query := `
		DELETE FROM tokens 
		WHERE scope = $1 AND user_id = $2
	`

	_, err := ts.db.Exec(query, scope, userId)
	return err
}
