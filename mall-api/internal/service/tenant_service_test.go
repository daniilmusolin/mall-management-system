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

func TestTenantService_Create_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	tn := &domain.Tenant{
		Name:         "ООО Тест",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Тестовая, д. 1",
		ContactPhone: "+79001234567",
		Email:        "test@test.ru",
	}

	mockRepo.On("ExistsByINN", mock.Anything, "123456789012").Return(false, nil)
	mockRepo.On("Create", mock.Anything, tn).Return(nil)

	err := svc.Create(context.Background(), tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantService_Create_InvalidINN(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	tn := &domain.Tenant{
		Name:         "ООО Тест",
		INN:          "123",
		LegalAddress: "г. Москва, ул. Тестовая, д. 1",
		Email:        "test@test.ru",
	}

	err := svc.Create(context.Background(), tn)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTenantService_Create_InvalidEmail(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	tn := &domain.Tenant{
		Name:         "ООО Тест",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Тестовая, д. 1",
		Email:        "not-email",
	}

	err := svc.Create(context.Background(), tn)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTenantService_Create_DuplicateINN(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	tn := &domain.Tenant{
		Name:         "ООО Тест",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Тестовая, д. 1",
		Email:        "test@test.ru",
	}

	mockRepo.On("ExistsByINN", mock.Anything, "123456789012").Return(true, nil)

	err := svc.Create(context.Background(), tn)

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestTenantService_GetAll(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	tenants := []domain.Tenant{{ID: 1, Name: "ООО Ромашка"}}

	mockRepo.On("GetAll", mock.Anything, pagination).Return(tenants, nil)

	res, err := svc.GetAll(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestTenantService_Update_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	tn := &domain.Tenant{
		ID:           1,
		Name:         "ООО Тест Обновлено",
		INN:          "123456789012",
		LegalAddress: "г. Москва, ул. Тестовая, д. 2",
		ContactPhone: "+79001234567",
		Email:        "test@test.ru",
	}

	mockRepo.On("Update", mock.Anything, tn).Return(nil)

	err := svc.Update(context.Background(), tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantService_Delete_Success(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantService_Delete_NotFound(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, domain.ErrNotFound)

	err := svc.Delete(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestTenantService_Delete_Error(t *testing.T) {
	mockRepo := new(port.MockTenantRepo)
	svc := service.NewTenantService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("db error"))

	err := svc.Delete(context.Background(), 1)

	assert.Error(t, err)
}
