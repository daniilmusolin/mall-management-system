package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mall-api/internal/domain"
	"mall-api/internal/handler"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func TestPaymentHandler_GetAll_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := service.NewPaymentService(mp, mc)
	h := handler.NewPaymentHandler(svc)

	payments := []domain.Payment{{ID: 1, Amount: 10000}}
	mp.On("GetAll", mock.Anything, mock.Anything).Return(payments, nil)

	req := httptest.NewRequest("GET", "/api/payments", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPaymentHandler_Create_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := service.NewPaymentService(mp, mc)
	h := handler.NewPaymentHandler(svc)

	newPayment := domain.Payment{
		ContractID:  1,
		Amount:      10000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	mc.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1, Status: "active"}, nil)
	mp.On("Create", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	body, _ := json.Marshal(newPayment)
	req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPaymentHandler_Create_NegativeAmount(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := service.NewPaymentService(mp, mc)
	h := handler.NewPaymentHandler(svc)

	newPayment := domain.Payment{
		ContractID:  1,
		Amount:      -500,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
	}

	body, _ := json.Marshal(newPayment)
	req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestPaymentHandler_Update_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := service.NewPaymentService(mp, mc)
	h := handler.NewPaymentHandler(svc)

	updatePayment := domain.Payment{
		Amount:      15000,
		PaymentDate: "2026-02-15",
		PaymentType: "rent",
		Status:      "completed",
	}

	mp.On("Update", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)

	body, _ := json.Marshal(updatePayment)
	req := httptest.NewRequest("PUT", "/api/payments/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPaymentHandler_Delete_Success(t *testing.T) {
	mp := new(port.MockPaymentRepo)
	mc := new(port.MockContractRepo)
	svc := service.NewPaymentService(mp, mc)
	h := handler.NewPaymentHandler(svc)

	mp.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/payments/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
