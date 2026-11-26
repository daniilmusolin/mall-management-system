//go:build integration

package mysql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/shopping_center_test?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestMySQLUserRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := mysql.NewUserRepo(db)
	ctx := context.Background()

	email := "unique_user@test.com"
	db.ExecContext(ctx, `DELETE FROM Users WHERE Email = ?`, email)

	user := &domain.User{
		Email:        email,
		PasswordHash: "hash",
		FullName:     "Test User",
		Role:         domain.RoleAdmin,
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	found, err := repo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.FullName, found.FullName)
	assert.Equal(t, user.Role, found.Role)

	_, err = repo.GetByEmail(ctx, "nonexistent@test.com")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
