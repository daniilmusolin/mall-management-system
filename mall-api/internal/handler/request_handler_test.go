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

func TestRequestHandler_GetAll_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	requests := []domain.ServiceRequest{{ID: 1, RequestType: "Electric"}}
	mr.On("GetAll", mock.Anything, mock.Anything).Return(requests, nil)

	req := httptest.NewRequest("GET", "/api/requests", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestHandler_Create_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	newRequest := domain.ServiceRequest{
		PremiseID:   1,
		RequestType: "Plumbing",
		Description: "Протекает кран в подсобке",
		Priority:    "high",
	}

	mp.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil)
	mr.On("Create", mock.Anything, mock.AnythingOfType("*domain.ServiceRequest")).Return(nil)

	body, _ := json.Marshal(newRequest)
	req := httptest.NewRequest("POST", "/api/requests", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRequestHandler_Create_InvalidJSON(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	req := httptest.NewRequest("POST", "/api/requests", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRequestHandler_Create_ValidationError(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	body := []byte(`{"premise_id":0,"request_type":"","description":"","priority":"unknown"}`)
	req := httptest.NewRequest("POST", "/api/requests", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRequestHandler_Update_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	updateRequest := domain.ServiceRequest{Status: "closed"}

	mr.On("Update", mock.Anything, mock.AnythingOfType("*domain.ServiceRequest")).Return(nil)

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest("PUT", "/api/requests/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestHandler_Delete_Success(t *testing.T) {
	mr := new(port.MockRequestRepo)
	mp := new(port.MockPremiseRepo)
	svc := service.NewRequestService(mr, mp)
	h := handler.NewRequestHandler(svc)

	mr.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/requests/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
