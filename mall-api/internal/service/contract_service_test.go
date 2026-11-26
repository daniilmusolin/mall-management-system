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

func newContractService(mc *port.MockContractRepo, mp *port.MockPremiseRepo, mt *port.MockTenantRepo, mtx *port.MockTxManager) *service.ContractService {
	return service.NewContractService(mc, mp, mt, mtx)
}

func TestContractService_Create_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  1,
		TenantID:   2,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "free"}, nil)
	mt.On("GetByID", mock.Anything, int64(2)).Return(&domain.Tenant{ID: 2}, nil)
	mc.On("CreateTx", mock.Anything, mock.Anything, contract).Return(nil)
	mp.On("UpdateStatusTx", mock.Anything, mock.Anything, int64(1), "rented").Return(nil)

	err := svc.Create(context.Background(), contract)

	assert.NoError(t, err)
	assert.Equal(t, "active", contract.Status)
	mc.AssertExpectations(t)
	mp.AssertExpectations(t)
	mt.AssertExpectations(t)
}

func TestContractService_Create_PremiseNotFound(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  99,
		TenantID:   2,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(99)).Return(nil, domain.ErrNotFound)

	err := svc.Create(context.Background(), contract)

	assert.Error(t, err)
	mc.AssertNotCalled(t, "CreateTx", mock.Anything, mock.Anything, mock.Anything)
}

func TestContractService_Create_PremiseRented(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  1,
		TenantID:   2,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "rented"}, nil)

	err := svc.Create(context.Background(), contract)

	assert.Error(t, err)
	mc.AssertNotCalled(t, "CreateTx", mock.Anything, mock.Anything, mock.Anything)
}

func TestContractService_Create_TenantNotFound(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  1,
		TenantID:   99,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "free"}, nil)
	mt.On("GetByID", mock.Anything, int64(99)).Return(nil, domain.ErrNotFound)

	err := svc.Create(context.Background(), contract)

	assert.Error(t, err)
	mc.AssertNotCalled(t, "CreateTx", mock.Anything, mock.Anything, mock.Anything)
}

func TestContractService_Create_InvalidDates(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  1,
		TenantID:   2,
		StartDate:  "2026-12-31",
		EndDate:    "2026-01-01",
		RatePerSqm: 1000,
	}

	err := svc.Create(context.Background(), contract)

	assert.Error(t, err)
	mc.AssertNotCalled(t, "CreateTx", mock.Anything, mock.Anything, mock.Anything)
}

func TestContractService_Create_Rollback(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		PremiseID:  1,
		TenantID:   2,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "free"}, nil)
	mt.On("GetByID", mock.Anything, int64(2)).Return(&domain.Tenant{ID: 2}, nil)
	mc.On("CreateTx", mock.Anything, mock.Anything, contract).Return(nil)
	mp.On("UpdateStatusTx", mock.Anything, mock.Anything, int64(1), "rented").Return(errors.New("db error"))

	err := svc.Create(context.Background(), contract)

	assert.Error(t, err)
}

func TestContractService_GetAll(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	contracts := []domain.LeaseContract{{ID: 1}}

	mc.On("GetAll", mock.Anything, pagination).Return(contracts, nil)

	res, err := svc.GetAll(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestContractService_Update_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		ID:         1,
		PremiseID:  1,
		TenantID:   2,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
		Status:     "terminated",
	}

	mc.On("Update", mock.Anything, contract).Return(nil)

	err := svc.Update(context.Background(), contract)

	assert.NoError(t, err)
	mc.AssertExpectations(t)
}

func TestContractService_Update_InvalidStatus(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	contract := &domain.LeaseContract{
		ID:     1,
		Status: "unknown",
	}

	err := svc.Update(context.Background(), contract)

	assert.Error(t, err)
	mc.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestContractService_Delete_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	mc.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1}, nil)
	mc.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mc.AssertExpectations(t)
}

func TestContractService_Delete_NotFound(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := newContractService(mc, mp, mt, mtx)

	mc.On("GetByID", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	err := svc.Delete(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mc.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}
