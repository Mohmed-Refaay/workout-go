package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserFromToken(token, scope string) (*User, error)
}

var ErrUsernameExisted = errors.New("username already exists")
var ErrEmailExisted = errors.New("email already exists")

func (pgStore *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	err := pgStore.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "users_username_key") {
				return ErrUsernameExisted
			}
			if strings.Contains(err.Error(), "users_email_key") {
				return ErrEmailExisted
			}
		}
		return err
	}

	return nil
}

func (pgStore *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		Username: username,
	}
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	err := pgStore.db.QueryRow(query, username).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pgStore *PostgresUserStore) GetUserFromToken(token, scope string) (*User, error) {
	hash := sha256.Sum256([]byte(token))

	user := &User{}

	query := `
	SELECT u.id, u.email, u.username, u.password_hash, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t
	ON u.id = t.user_id
	WHERE t.hash = $1 AND t.scope = $2 AND expiry > $3
	`

	err := pgStore.db.QueryRow(query, hash[:], scope, time.Now()).Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
