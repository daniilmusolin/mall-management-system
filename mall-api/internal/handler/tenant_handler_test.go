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

func TestTenantHandler_GetAll_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	tenants := []domain.Tenant{{ID: 1, Name: "ООО Ромашка"}}
	mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(tenants, nil)

	req := httptest.NewRequest("GET", "/api/tenants", nil)
	w := httptest.NewRecorder()

	h.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTenantHandler_Create_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	newTenant := domain.Tenant{
		Name:         "ИП Иванов",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Ленина, д. 1",
		ContactPhone: "+79001234567",
		Email:        "ivanov@test.ru",
	}

	mockRepo.On("ExistsByINN", mock.Anything, "123456789012").Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Tenant")).Return(nil)

	body, _ := json.Marshal(newTenant)
	req := httptest.NewRequest("POST", "/api/tenants", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestTenantHandler_Create_InvalidJSON(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	req := httptest.NewRequest("POST", "/api/tenants", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTenantHandler_Create_ValidationError(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	body := []byte(`{"name":"","inn":"abc"}`)
	req := httptest.NewRequest("POST", "/api/tenants", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestTenantHandler_Update_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	updateTenant := domain.Tenant{
		Name:         "ИП Иванов Обновленный",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Ленина, д. 2",
		ContactPhone: "+79001234567",
		Email:        "ivanov@test.ru",
	}

	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Tenant")).Return(nil)

	body, _ := json.Marshal(updateTenant)
	req := httptest.NewRequest("PUT", "/api/tenants/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTenantHandler_Delete_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)
	h := handler.NewTenantHandler(svc)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/tenants/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
