//go:build integration

package mysql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLContractRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	roomNumber := "ROOM-FOR-CONTRACT"
	inn := "888888888888"
	db.ExecContext(ctx, `DELETE FROM LeaseContracts WHERE PremiseId IN (SELECT PremiseId FROM Premises WHERE RoomNumber = ?)`, roomNumber)
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)
	db.ExecContext(ctx, `DELETE FROM Tenants WHERE INN = ?`, inn)

	premiseRepo := mysql.NewPremiseRepo(db)
	premise := &domain.Premise{RoomNumber: roomNumber, Floor: 1, Area: 10, Type: "Test", Zone: "A", Status: "free"}
	assert.NoError(t, premiseRepo.Create(ctx, premise))

	tenantRepo := mysql.NewTenantRepo(db)
	tenant := &domain.Tenant{Name: "Tenant For Contract", INN: inn, LegalAddress: "addr", ContactPhone: "+70000000000", Email: "t@test.ru"}
	assert.NoError(t, tenantRepo.Create(ctx, tenant))

	repo := mysql.NewContractRepo(db)
	txManager := mysql.NewTransactionManager(db)

	contract := &domain.LeaseContract{
		PremiseID:     premise.ID,
		TenantID:      tenant.ID,
		StartDate:     "2026-01-01",
		EndDate:       "2026-12-31",
		RatePerSqm:    1000,
		DepositAmount: 5000,
		Status:        "active",
	}

	err := txManager.WithTransaction(ctx, func(tx interface{}) error {
		return repo.CreateTx(ctx, tx, contract)
	})
	assert.NoError(t, err)
	assert.NotZero(t, contract.ID)

	found, err := repo.GetByID(ctx, contract.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, premise.ID, found.PremiseID)
	assert.Equal(t, tenant.ID, found.TenantID)
	assert.Equal(t, "2026-01-01", found.StartDate)
	assert.Equal(t, "2026-12-31", found.EndDate)
	assert.Equal(t, 1000.0, found.RatePerSqm)
	assert.Equal(t, 5000.0, found.DepositAmount)
	assert.Equal(t, "active", found.Status)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	list, err := repo.GetAll(ctx, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	contract.Status = "terminated"
	contract.RatePerSqm = 1200
	err = repo.Update(ctx, contract)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, contract.ID)
	assert.NoError(t, err)
	assert.Equal(t, "terminated", updated.Status)
	assert.Equal(t, 1200.0, updated.RatePerSqm)

	err = repo.Delete(ctx, contract.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, contract.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.Delete(ctx, contract.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
