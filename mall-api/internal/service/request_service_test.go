package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mall-api/internal/domain"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func newRequestService(mr *port.MockRequestRepo, mp *port.MockPremiseRepo) *service.RequestService {
	return service.NewRequestService(mr, mp)
}

func TestRequestService_Create_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{
		PremiseID:   1,
		RequestType: "Сантехника",
		Description: "Протекает кран",
		Priority:    "high",
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil)
	mr.On("Create", mock.Anything, sr).Return(nil)

	err := svc.Create(context.Background(), sr)

	assert.NoError(t, err)
	assert.Equal(t, "open", sr.Status)
	mr.AssertExpectations(t)
	mp.AssertExpectations(t)
}

func TestRequestService_Create_InvalidPriority(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{
		PremiseID:   1,
		RequestType: "Сантехника",
		Description: "Протекает кран",
		Priority:    "unknown",
	}

	err := svc.Create(context.Background(), sr)

	assert.Error(t, err)
	mr.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRequestService_Create_ShortDescription(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{
		PremiseID:   1,
		RequestType: "Сантехника",
		Description: "ok",
		Priority:    "high",
	}

	err := svc.Create(context.Background(), sr)

	assert.Error(t, err)
	mr.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRequestService_Create_PremiseNotFound(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{
		PremiseID:   99,
		RequestType: "Сантехника",
		Description: "Протекает кран",
		Priority:    "high",
	}

	mp.On("GetByID", mock.Anything, int64(99)).Return(nil, domain.ErrNotFound)

	err := svc.Create(context.Background(), sr)

	assert.Error(t, err)
	mr.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestRequestService_GetAll(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	requests := []domain.ServiceRequest{{ID: 1}}

	mr.On("GetAll", mock.Anything, pagination).Return(requests, nil)

	res, err := svc.GetAll(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestRequestService_Update_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{ID: 1, Status: "closed"}

	mr.On("Update", mock.Anything, sr).Return(nil)

	err := svc.Update(context.Background(), sr)

	assert.NoError(t, err)
	mr.AssertExpectations(t)
}

func TestRequestService_Update_InvalidStatus(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	sr := &domain.ServiceRequest{ID: 1, Status: "unknown"}

	err := svc.Update(context.Background(), sr)

	assert.Error(t, err)
	mr.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestRequestService_Delete_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := newRequestService(mr, mp)

	mr.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mr.AssertExpectations(t)
}
