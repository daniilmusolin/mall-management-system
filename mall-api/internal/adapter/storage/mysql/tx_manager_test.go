//go:build integration

package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLTransactionManager_Commit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	roomNumber := "TX-COMMIT-ROOM"
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)

	txManager := mysql.NewTransactionManager(db)
	premiseRepo := mysql.NewPremiseRepo(db)

	premise := &domain.Premise{
		RoomNumber: roomNumber,
		Floor:      1,
		Area:       10,
		Type:       "Test",
		Zone:       "A",
		Status:     "free",
	}

	err := txManager.WithTransaction(ctx, func(tx interface{}) error {
		return premiseRepo.CreateTx(ctx, tx, premise)
	})
	assert.NoError(t, err)
	assert.NotZero(t, premise.ID)

	found, err := premiseRepo.GetByID(ctx, premise.ID)
	assert.NoError(t, err)
	assert.Equal(t, roomNumber, found.RoomNumber)

	db.ExecContext(ctx, `DELETE FROM Premises WHERE PremiseId = ?`, premise.ID)
}

func TestMySQLTransactionManager_Rollback(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	roomNumber := "TX-ROLLBACK-ROOM"
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)

	txManager := mysql.NewTransactionManager(db)
	premiseRepo := mysql.NewPremiseRepo(db)

	premise := &domain.Premise{
		RoomNumber: roomNumber,
		Floor:      1,
		Area:       10,
		Type:       "Test",
		Zone:       "A",
		Status:     "free",
	}

	expectedErr := errors.New("forced rollback")
	err := txManager.WithTransaction(ctx, func(tx interface{}) error {
		if err := premiseRepo.CreateTx(ctx, tx, premise); err != nil {
			return err
		}
		return expectedErr
	})
	assert.ErrorIs(t, err, expectedErr)

	var count int
	db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Premises WHERE RoomNumber = ?`, roomNumber).Scan(&count)
	assert.Equal(t, 0, count)
}
