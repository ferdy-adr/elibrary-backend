package users

import (
	"database/sql"

	"github.com/ferdy-adr/elibrary-backend/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, full_name) 
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.FullName)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *Repository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password, full_name, created_at, updated_at 
		FROM users 
		WHERE username = ?
	`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password, full_name, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.FullName, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
