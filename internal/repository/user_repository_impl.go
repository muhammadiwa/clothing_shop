package repository

import (
	"database/sql"
	"time"

	"clothing-shop-api/internal/domain/models"
	"clothing-shop-api/pkg/utils"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *models.User) error {
	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO users (
            username, email, password, role, is_verified, verification_token, 
            phone_number
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	result, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		hashedPassword,
		user.Role,
		user.IsVerified,
		user.VerificationToken,
		user.PhoneNumber,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

func (r *userRepositoryImpl) FindByID(id uint) (*models.User, error) {
	query := `
        SELECT 
            id, username, email, password, role, is_verified, verification_token, 
            reset_token, reset_token_expiry, phone_number, created_at, updated_at 
        FROM users 
        WHERE id = ? AND deleted_at IS NULL
    `

	row := r.db.QueryRow(query, id)

	user := &models.User{}
	var resetTokenExpiry sql.NullTime

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.IsVerified, &user.VerificationToken, &user.ResetToken,
		&resetTokenExpiry, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if resetTokenExpiry.Valid {
		user.ResetTokenExpiry = &resetTokenExpiry.Time
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	query := `
        SELECT 
            id, username, email, password, role, is_verified, verification_token, 
            reset_token, reset_token_expiry, phone_number, created_at, updated_at 
        FROM users 
        WHERE email = ? AND deleted_at IS NULL
    `

	row := r.db.QueryRow(query, email)

	user := &models.User{}
	var resetTokenExpiry sql.NullTime

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.IsVerified, &user.VerificationToken, &user.ResetToken,
		&resetTokenExpiry, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if resetTokenExpiry.Valid {
		user.ResetTokenExpiry = &resetTokenExpiry.Time
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByVerificationToken(token string) (*models.User, error) {
	query := `
        SELECT 
            id, username, email, password, role, is_verified, verification_token, 
            reset_token, reset_token_expiry, phone_number, created_at, updated_at 
        FROM users 
        WHERE verification_token = ? AND deleted_at IS NULL
    `

	row := r.db.QueryRow(query, token)

	user := &models.User{}
	var resetTokenExpiry sql.NullTime

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.IsVerified, &user.VerificationToken, &user.ResetToken,
		&resetTokenExpiry, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if resetTokenExpiry.Valid {
		user.ResetTokenExpiry = &resetTokenExpiry.Time
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByResetToken(token string) (*models.User, error) {
	query := `
        SELECT 
            id, username, email, password, role, is_verified, verification_token, 
            reset_token, reset_token_expiry, phone_number, created_at, updated_at 
        FROM users 
        WHERE reset_token = ? AND reset_token_expiry > NOW() AND deleted_at IS NULL
    `

	row := r.db.QueryRow(query, token)

	user := &models.User{}
	var resetTokenExpiry sql.NullTime

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role,
		&user.IsVerified, &user.VerificationToken, &user.ResetToken,
		&resetTokenExpiry, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if resetTokenExpiry.Valid {
		user.ResetTokenExpiry = &resetTokenExpiry.Time
	}

	return user, nil
}

func (r *userRepositoryImpl) Update(user *models.User) error {
	query := `
        UPDATE users 
        SET 
            username = ?,
            email = ?,
            is_verified = ?,
            verification_token = ?,
            reset_token = ?,
            reset_token_expiry = ?,
            phone_number = ?
        WHERE id = ? AND deleted_at IS NULL
    `

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.IsVerified,
		user.VerificationToken,
		user.ResetToken,
		user.ResetTokenExpiry,
		user.PhoneNumber,
		user.ID,
	)

	return err
}

func (r *userRepositoryImpl) UpdatePassword(id uint, password string) error {
	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = ? WHERE id = ? AND deleted_at IS NULL`
	_, err = r.db.Exec(query, hashedPassword, id)
	return err
}

func (r *userRepositoryImpl) Delete(id uint) error {
	// Soft delete
	query := `UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
