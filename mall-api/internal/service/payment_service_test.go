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

func newPaymentService(mp *port.MockPaymentRepo, mc *port.MockContractRepo) *service.PaymentService {
	return service.NewPaymentService(mp, mc)
}

func TestPaymentService_Create_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ContractID:  1,
		Amount:      10000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	mc.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1, Status: "active"}, nil)
	mp.On("Create", mock.Anything, p).Return(nil)

	err := svc.Create(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, "completed", p.Status)
	mp.AssertExpectations(t)
	mc.AssertExpectations(t)
}

func TestPaymentService_Create_InvalidAmount(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ContractID:  1,
		Amount:      -100,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mp.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPaymentService_Create_InvalidType(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ContractID:  1,
		Amount:      1000,
		PaymentDate: "2026-02-15",
		PaymentType: "unknown",
	}

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mp.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPaymentService_Create_ContractNotFound(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ContractID:  99,
		Amount:      1000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	mc.On("GetByID", mock.Anything, int64(99)).Return(nil, domain.ErrNotFound)

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mp.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPaymentService_Create_InactiveContract(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ContractID:  1,
		Amount:      1000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	mc.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1, Status: "terminated"}, nil)

	err := svc.Create(context.Background(), p)

	assert.Error(t, err)
	mp.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPaymentService_GetAll(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	payments := []domain.Payment{{ID: 1}}

	mp.On("GetAll", mock.Anything, pagination).Return(payments, nil)

	res, err := svc.GetAll(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestPaymentService_Update_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{
		ID:          1,
		Amount:      15000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
		Status:      "completed",
	}

	mp.On("Update", mock.Anything, p).Return(nil)

	err := svc.Update(context.Background(), p)

	assert.NoError(t, err)
	mp.AssertExpectations(t)
}

func TestPaymentService_Update_InvalidAmount(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	p := &domain.Payment{ID: 1, Amount: -5}

	err := svc.Update(context.Background(), p)

	assert.Error(t, err)
	mp.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPaymentService_Delete_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := newPaymentService(mp, mc)

	mp.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mp.AssertExpectations(t)
}
