//go:build integration

package mysql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/domain"
)

func TestMySQLPremiseRepo_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := mysql.NewPremiseRepo(db)
	ctx := context.Background()

	roomNumber := "UNIQUE-ROOM-101"
	db.ExecContext(ctx, `DELETE FROM Premises WHERE RoomNumber = ?`, roomNumber)

	premise := &domain.Premise{
		RoomNumber: roomNumber,
		Floor:      1,
		Area:       45.5,
		Type:       "Shop",
		Zone:       "A",
		Status:     "free",
	}

	err := repo.Create(ctx, premise)
	assert.NoError(t, err)
	assert.NotZero(t, premise.ID)

	found, err := repo.GetByID(ctx, premise.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, roomNumber, found.RoomNumber)
	assert.Equal(t, 1, found.Floor)
	assert.Equal(t, 45.5, found.Area)
	assert.Equal(t, "Shop", found.Type)
	assert.Equal(t, "A", found.Zone)
	assert.Equal(t, "free", found.Status)

	exists, err := repo.ExistsByRoomNumber(ctx, roomNumber)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByRoomNumber(ctx, "NO-SUCH-ROOM")
	assert.NoError(t, err)
	assert.False(t, exists)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	list, err := repo.GetAll(ctx, pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	premise.Status = "maintenance"
	err = repo.Update(ctx, premise)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, premise.ID)
	assert.NoError(t, err)
	assert.Equal(t, "maintenance", updated.Status)

	err = repo.Delete(ctx, premise.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, premise.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	err = repo.Delete(ctx, premise.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
