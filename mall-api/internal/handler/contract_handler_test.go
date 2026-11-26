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

func TestContractHandler_GetAll_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	contracts := []domain.LeaseContract{{ID: 1, PremiseID: 1, TenantID: 1}}
	mc.On("GetAll", mock.Anything, mock.Anything).Return(contracts, nil)

	req := httptest.NewRequest("GET", "/api/contracts", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestContractHandler_Create_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	newContract := domain.LeaseContract{
		PremiseID:  1,
		TenantID:   1,
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "free"}, nil)
	mt.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil)
	mc.On("CreateTx", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.LeaseContract")).Return(nil)
	mp.On("UpdateStatusTx", mock.Anything, mock.Anything, int64(1), "rented").Return(nil)

	body, _ := json.Marshal(newContract)
	req := httptest.NewRequest("POST", "/api/contracts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestContractHandler_Create_InvalidJSON(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	req := httptest.NewRequest("POST", "/api/contracts", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestContractHandler_Create_ValidationError(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	body := []byte(`{"premise_id":0,"tenant_id":0,"start_date":"","end_date":"","rate_per_sqm":-1}`)
	req := httptest.NewRequest("POST", "/api/contracts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestContractHandler_Update_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	updateContract := domain.LeaseContract{
		StartDate:  "2026-01-01",
		EndDate:    "2026-12-31",
		RatePerSqm: 1000,
		Status:     "terminated",
	}

	mc.On("Update", mock.Anything, mock.AnythingOfType("*domain.LeaseContract")).Return(nil)

	body, _ := json.Marshal(updateContract)
	req := httptest.NewRequest("PUT", "/api/contracts/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestContractHandler_Delete_Success(t *testing.T) {
	mc := new(port.MockContractRepo)
	mp := new(port.MockPremiseRepo)
	mt := new(port.MockTenantRepo)
	mtx := new(port.MockTxManager)
	svc := service.NewContractService(mc, mp, mt, mtx)
	h := handler.NewContractHandler(svc)

	mc.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1}, nil)
	mc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/contracts/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
