//go:build integration

package mysql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLTenantRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := mysql.NewTenantRepo(db)
	ctx := context.Background()

	inn := "999999999999"
	db.ExecContext(ctx, `DELETE FROM Tenants WHERE INN = ?`, inn)

	tenant := &domain.Tenant{
		Name:         "ООО Тест",
		INN:          inn,
		LegalAddress: "г. Тест, ул. Тестовая, д. 1",
		ContactPhone: "+70000000000",
		Email:        "test@test.ru",
	}

	err := repo.Create(ctx, tenant)
	assert.NoError(t, err)
	assert.NotZero(t, tenant.ID)

	found, err := repo.GetByID(ctx, tenant.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, tenant.Name, found.Name)
	assert.Equal(t, inn, found.INN)
	assert.Equal(t, tenant.LegalAddress, found.LegalAddress)
	assert.Equal(t, tenant.ContactPhone, found.ContactPhone)
	assert.Equal(t, tenant.Email, found.Email)

	exists, err := repo.ExistsByINN(ctx, inn)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByINN(ctx, "000000000000")
	assert.NoError(t, err)
	assert.False(t, exists)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	list, err := repo.GetAll(ctx, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	tenant.Name = "ООО Тест Обновлено"
	err = repo.Update(ctx, tenant)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, tenant.ID)
	assert.NoError(t, err)
	assert.Equal(t, "ООО Тест Обновлено", updated.Name)

	err = repo.Delete(ctx, tenant.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, tenant.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.Delete(ctx, tenant.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
