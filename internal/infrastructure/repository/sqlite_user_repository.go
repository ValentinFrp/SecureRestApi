package repository

import (
	"database/sql"
	"time"

	"github.com/valentinfrappart/securerestapi/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{
		db: db,
	}
}

func (r *SQLiteUserRepository) Create(email, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(query, email, passwordHash, now, now)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return user, nil
}

func (r *SQLiteUserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) FindByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
