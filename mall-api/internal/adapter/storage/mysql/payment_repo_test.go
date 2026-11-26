//go:build integration

package mysql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLPaymentRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	roomNumber := "ROOM-FOR-PAYMENT"
	inn := "777777777777"
	refNumber := "REF-UNIQUE-1"

	db.ExecContext(ctx, `DELETE FROM Payments WHERE ReferenceNumber = ?`, refNumber)
	db.ExecContext(ctx, `DELETE FROM LeaseContracts WHERE PremiseId IN (SELECT PremiseId FROM Premises WHERE RoomNumber = ?)`, roomNumber)
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)
	db.ExecContext(ctx, `DELETE FROM Tenants WHERE INN = ?`, inn)

	premiseRepo := mysql.NewPremiseRepo(db)
	premise := &domain.Premise{RoomNumber: roomNumber, Floor: 1, Area: 10, Type: "Test", Zone: "A", Status: "free"}
	assert.NoError(t, premiseRepo.Create(ctx, premise))

	tenantRepo := mysql.NewTenantRepo(db)
	tenant := &domain.Tenant{Name: "Tenant For Payment", INN: inn, LegalAddress: "addr", ContactPhone: "+70000000000", Email: "t@test.ru"}
	assert.NoError(t, tenantRepo.Create(ctx, tenant))

	contractRepo := mysql.NewContractRepo(db)
	txManager := mysql.NewTransactionManager(db)
	contract := &domain.LeaseContract{
		PremiseID:     premise.ID,
		TenantID:      tenant.ID,
		StartDate:     "2026-01-01",
		EndDate:       "2026-12-31",
		RatePerSqm:    1000,
		DepositAmount: 0,
		Status:        "active",
	}
	err := txManager.WithTransaction(ctx, func(tx interface{}) error {
		return contractRepo.CreateTx(ctx, tx, contract)
	})
	assert.NoError(t, err)
	assert.NotZero(t, contract.ID)

	repo := mysql.NewPaymentRepo(db)
	payment := &domain.Payment{
		ContractID:      contract.ID,
		Amount:          15000.50,
		PaymentDate:     "2026-02-15",
		PaymentType:     "rent",
		Status:          "completed",
		ReferenceNumber: refNumber,
	}

	err = repo.Create(ctx, payment)
	assert.NoError(t, err)
	assert.NotZero(t, payment.ID)

	found, err := repo.GetByID(ctx, payment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, contract.ID, found.ContractID)
	assert.Equal(t, 15000.50, found.Amount)
	assert.Equal(t, "2026-02-15", found.PaymentDate)
	assert.Equal(t, "rent", found.PaymentType)
	assert.Equal(t, "completed", found.Status)
	assert.Equal(t, refNumber, found.ReferenceNumber)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	list, err := repo.GetAll(ctx, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	payment.Amount = 16000.00
	payment.Status = "pending"
	err = repo.Update(ctx, payment)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, payment.ID)
	assert.NoError(t, err)
	assert.Equal(t, 16000.00, updated.Amount)
	assert.Equal(t, "pending", updated.Status)

	err = repo.Delete(ctx, payment.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, payment.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.Delete(ctx, payment.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
