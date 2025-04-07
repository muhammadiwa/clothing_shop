package repository

import (
	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/internal/repository"
	"database/sql"
	"errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := "SELECT id, username, email, password FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	query := "UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?"
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepository) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
