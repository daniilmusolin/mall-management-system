//go:build integration

package mysql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLRequestRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	roomNumber := "ROOM-FOR-REQUEST"
	db.ExecContext(ctx, `DELETE FROM ServiceRequests WHERE PremiseId IN (SELECT PremiseId FROM Premises WHERE RoomNumber = ?)`, roomNumber)
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)

	premiseRepo := mysql.NewPremiseRepo(db)
	premise := &domain.Premise{RoomNumber: roomNumber, Floor: 1, Area: 10, Type: "Test", Zone: "A", Status: "free"}
	assert.NoError(t, premiseRepo.Create(ctx, premise))

	repo := mysql.NewRequestRepo(db)
	request := &domain.ServiceRequest{
		PremiseID:   premise.ID,
		RequestType: "Сантехника",
		Description: "Протекает кран",
		Priority:    "high",
		Status:      "open",
	}

	err := repo.Create(ctx, request)
	assert.NoError(t, err)
	assert.NotZero(t, request.ID)

	found, err := repo.GetByID(ctx, request.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, premise.ID, found.PremiseID)
	assert.Equal(t, "Сантехника", found.RequestType)
	assert.Equal(t, "Протекает кран", found.Description)
	assert.Equal(t, "high", found.Priority)
	assert.Equal(t, "open", found.Status)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	list, err := repo.GetAll(ctx, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	request.Status = "closed"
	err = repo.Update(ctx, request)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, request.ID)
	assert.NoError(t, err)
	assert.Equal(t, "closed", updated.Status)

	err = repo.Delete(ctx, request.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, request.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.Delete(ctx, request.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
