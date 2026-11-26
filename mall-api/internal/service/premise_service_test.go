package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mall-api/internal/domain"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func TestPremiseService_Create_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		RoomNumber: "A-101",
		Floor:      1,
		Area:       45.5,
		Type:       "Shop",
		Zone:       "A",
	}

	mockRepo.On("ExistsByRoomNumber", mock.Anything, "A-101").Return(false, nil)
	mockRepo.On("Create", mock.Anything, p).Return(nil)

	err := svc.Create(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, "free", p.Status)
	mockRepo.AssertExpectations(t)
}

func TestPremiseService_Create_InvalidArea(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		RoomNumber: "A-101",
		Floor:      1,
		Area:       -10.0,
		Type:       "Shop",
	}

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPremiseService_Create_InvalidRoomNumber(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		RoomNumber: "",
		Floor:      1,
		Area:       10,
		Type:       "Shop",
	}

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPremiseService_Create_DuplicateRoomNumber(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		RoomNumber: "A-101",
		Floor:      1,
		Area:       10,
		Type:       "Shop",
	}

	mockRepo.On("ExistsByRoomNumber", mock.Anything, "A-101").Return(true, nil)

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPremiseService_GetAll(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	premises := []domain.Premise{{ID: 1, RoomNumber: "A-101"}}

	mockRepo.On("GetAll", mock.Anything, pagination).Return(premises, nil)

	res, err := svc.GetAll(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestPremiseService_Update_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		ID:         1,
		RoomNumber: "A-101",
		Floor:      1,
		Area:       10,
		Type:       "Shop",
		Status:     "free",
	}

	mockRepo.On("Update", mock.Anything, p).Return(nil)

	err := svc.Update(context.Background(), p)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPremiseService_Update_InvalidStatus(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	p := &domain.Premise{
		ID:         1,
		RoomNumber: "A-101",
		Floor:      1,
		Area:       10,
		Type:       "Shop",
		Status:     "unknown_status",
	}

	err := svc.Update(context.Background(), p)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPremiseService_Delete_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPremiseService_Delete_NotFound(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	err := svc.Delete(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestPremiseService_Delete_Error(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("db error"))

	err := svc.Delete(context.Background(), 1)

	assert.Error(t, err)
}
