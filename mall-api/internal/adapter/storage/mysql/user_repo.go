package mysql

import (
	"context"
	"database/sql"
	"time"

	"mall-api/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO Users (Email, PasswordHash, FullName, Role) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, u.Email, u.PasswordHash, u.FullName, u.Role)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT UserId, Email, PasswordHash, FullName, Role, CreatedAt FROM Users WHERE Email = ?`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.User, error) {
	query := `SELECT UserId, Email, PasswordHash, FullName, Role, CreatedAt FROM Users ORDER BY UserId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT UserId, Email, PasswordHash, FullName, Role, CreatedAt FROM Users WHERE UserId = ?`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	query := `UPDATE Users SET FullName = ?, Role = ? WHERE UserId = ?`
	res, err := r.db.ExecContext(ctx, query, u.FullName, u.Role, u.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM Users WHERE UserId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *userRepo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `INSERT INTO RefreshTokens (UserId, Token, ExpiresAt) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (r *userRepo) GetUserIDByRefreshToken(ctx context.Context, token string) (int64, error) {
	query := `SELECT UserId FROM RefreshTokens WHERE Token = ? AND ExpiresAt > NOW()`
	var userID int64
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, domain.ErrNotFound
	}
	return userID, err
}

func (r *userRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM RefreshTokens WHERE Token = ?`, token)
	return err
}