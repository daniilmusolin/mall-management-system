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

func TestPremiseHandler_GetAll_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	premises := []domain.Premise{{ID: 1, RoomNumber: "A-101"}}
	mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(premises, nil)

	req := httptest.NewRequest("GET", "/api/premises?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "A-101")
}

func TestPremiseHandler_Create_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	newPremise := domain.Premise{
		RoomNumber: "B-201",
		Floor:      2,
		Area:       50.0,
		Type:       "Shop",
		Zone:       "B",
	}

	mockRepo.On("ExistsByRoomNumber", mock.Anything, "B-201").Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Premise")).Return(nil)

	body, _ := json.Marshal(newPremise)
	req := httptest.NewRequest("POST", "/api/premises", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPremiseHandler_Create_ValidationError(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	body := []byte(`{"room_number":"","floor":-1,"area":0,"type":""}`)
	req := httptest.NewRequest("POST", "/api/premises", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestPremiseHandler_Update_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	updatePremise := domain.Premise{
		RoomNumber: "B-201-UPDATED",
		Floor:      2,
		Area:       55.0,
		Type:       "Shop",
		Zone:       "B",
		Status:     "free",
	}

	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Premise")).Return(nil)

	body, _ := json.Marshal(updatePremise)
	req := httptest.NewRequest("PUT", "/api/premises/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPremiseHandler_Update_InvalidID(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	body, _ := json.Marshal(domain.Premise{RoomNumber: "X", Area: 1, Type: "Shop"})
	req := httptest.NewRequest("PUT", "/api/premises/abc", bytes.NewBuffer(body))
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPremiseHandler_Delete_Success(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/premises/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPremiseHandler_Delete_NotFound(t *testing.T) {
	mockRepo := new(port.MockPremiseRepo)
	svc := service.NewPremiseService(mockRepo)
	h := handler.NewPremiseHandler(svc)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	req := httptest.NewRequest("DELETE", "/api/premises/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
