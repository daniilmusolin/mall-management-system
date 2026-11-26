package handler_test

import (
	"bytes"
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

func setupFullRouter() (http.Handler, *port.MockPremiseRepo, *port.MockTenantRepo, *port.MockContractRepo, *port.MockPaymentRepo, *port.MockRequestRepo, *port.MockUserRepo) {
	premiseRepo := new(port.MockPremiseRepo)
	tenantRepo := new(port.MockTenantRepo)
	contractRepo := new(port.MockContractRepo)
	paymentRepo := new(port.MockPaymentRepo)
	requestRepo := new(port.MockRequestRepo)
	userRepo := new(port.MockUserRepo)
	txManager := new(port.MockTxManager)

	authSvc := service.NewAuthService(userRepo, "secret")
	premiseSvc := service.NewPremiseService(premiseRepo)
	tenantSvc := service.NewTenantService(tenantRepo)
	contractSvc := service.NewContractService(contractRepo, premiseRepo, tenantRepo, txManager)
	paymentSvc := service.NewPaymentService(paymentRepo, contractRepo)
	requestSvc := service.NewRequestService(requestRepo, premiseRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	premiseHandler := handler.NewPremiseHandler(premiseSvc)
	tenantHandler := handler.NewTenantHandler(tenantSvc)
	contractHandler := handler.NewContractHandler(contractSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)
	requestHandler := handler.NewRequestHandler(requestSvc)

	authMid := handler.AuthMiddleware(authSvc)
	mux := http.NewServeMux()

	authHandler.RegisterRoutes(mux, authMid)
	premiseHandler.RegisterRoutes(mux, authMid)
	tenantHandler.RegisterRoutes(mux, authMid)
	contractHandler.RegisterRoutes(mux, authMid)
	paymentHandler.RegisterRoutes(mux, authMid)
	requestHandler.RegisterRoutes(mux, authMid)

	return mux, premiseRepo, tenantRepo, contractRepo, paymentRepo, requestRepo, userRepo
}

func TestRBAC_Users(t *testing.T) {
	router, _, _, _, _, _, userRepo := setupFullRouter()

	t.Run("Admin can view users", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		userRepo.On("GetAll", mock.Anything, mock.Anything).Return([]domain.User{{ID: 1}}, nil).Once()

		req := httptest.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Manager cannot view users", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Admin can delete user", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		userRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Manager cannot delete user", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("DELETE", "/api/users/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRBAC_Premises(t *testing.T) {
	router, premiseRepo, _, _, _, _, _ := setupFullRouter()
	validJson := []byte(`{"room_number":"A-1","floor":1,"area":10,"type":"Shop"}`)

	t.Run("Admin can create premise", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		premiseRepo.On("ExistsByRoomNumber", mock.Anything, "A-1").Return(false, nil).Once()
		premiseRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Premise")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/api/premises", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Manager cannot create premise", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("POST", "/api/premises", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Admin can delete premise", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		premiseRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil).Once()
		premiseRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/premises/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Manager cannot delete premise", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("DELETE", "/api/premises/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRBAC_Tenants(t *testing.T) {
	router, _, tenantRepo, _, _, _, _ := setupFullRouter()
	validJson := []byte(`{"name":"ООО Ромашка","inn":"123456789012","legal_address":"г. Москва, ул. Ленина, 1","contact_phone":"+79001234567","email":"r@test.ru"}`)

	t.Run("Manager can create tenant", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		tenantRepo.On("ExistsByINN", mock.Anything, "123456789012").Return(false, nil).Once()
		tenantRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Tenant")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/api/tenants", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Accountant cannot create tenant", func(t *testing.T) {
		token := generateTestToken(domain.RoleAccountant, "secret")
		req := httptest.NewRequest("POST", "/api/tenants", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Admin can delete tenant", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		tenantRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil).Once()
		tenantRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/tenants/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Manager cannot delete tenant", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("DELETE", "/api/tenants/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRBAC_Contracts(t *testing.T) {
	router, premiseRepo, tenantRepo, contractRepo, _, _, _ := setupFullRouter()
	validJson := []byte(`{"premise_id":1,"tenant_id":1,"start_date":"2026-01-01","end_date":"2026-12-31","rate_per_sqm":1000}`)

	t.Run("Manager can create contract", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")

		premiseRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1, Status: "free"}, nil).Once()
		tenantRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Tenant{ID: 1}, nil).Once()
		contractRepo.On("CreateTx", mock.Anything, mock.Anything, mock.AnythingOfType("*domain.LeaseContract")).Return(nil).Once()
		premiseRepo.On("UpdateStatusTx", mock.Anything, mock.Anything, int64(1), "rented").Return(nil).Once()

		req := httptest.NewRequest("POST", "/api/contracts", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Accountant cannot create contract", func(t *testing.T) {
		token := generateTestToken(domain.RoleAccountant, "secret")
		req := httptest.NewRequest("POST", "/api/contracts", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Admin can delete contract", func(t *testing.T) {
		token := generateTestToken(domain.RoleAdmin, "secret")
		contractRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1}, nil).Once()
		contractRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/api/contracts/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestRBAC_Payments(t *testing.T) {
	router, _, _, contractRepo, paymentRepo, _, _ := setupFullRouter()
	validJson := []byte(`{"contract_id":1,"amount":1000,"payment_date":"2026-01-01","payment_type":"rent"}`)

	t.Run("Accountant can create payment", func(t *testing.T) {
		token := generateTestToken(domain.RoleAccountant, "secret")
		contractRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.LeaseContract{ID: 1, Status: "active"}, nil).Once()
		paymentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Manager cannot create payment", func(t *testing.T) {
		token := generateTestToken(domain.RoleManager, "secret")
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Accountant cannot delete payment", func(t *testing.T) {
		token := generateTestToken(domain.RoleAccountant, "secret")
		req := httptest.NewRequest("DELETE", "/api/payments/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRBAC_Requests(t *testing.T) {
	router, premiseRepo, _, _, _, requestRepo, _ := setupFullRouter()
	validJson := []byte(`{"premise_id":1,"request_type":"Plumbing","description":"Протекает кран","priority":"high"}`)

	t.Run("Support can create request", func(t *testing.T) {
		token := generateTestToken(domain.RoleSupport, "secret")
		premiseRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Premise{ID: 1}, nil).Once()
		requestRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.ServiceRequest")).Return(nil).Once()

		req := httptest.NewRequest("POST", "/api/requests", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Technician cannot create request", func(t *testing.T) {
		token := generateTestToken(domain.RoleTechnician, "secret")
		req := httptest.NewRequest("POST", "/api/requests", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Technician can update request", func(t *testing.T) {
		token := generateTestToken(domain.RoleTechnician, "secret")
		requestRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.ServiceRequest")).Return(nil).Once()

		req := httptest.NewRequest("PUT", "/api/requests/1", bytes.NewBuffer(validJson))
		req.SetPathValue("id", "1")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Support cannot update request", func(t *testing.T) {
		token := generateTestToken(domain.RoleSupport, "secret")
		req := httptest.NewRequest("PUT", "/api/requests/1", bytes.NewBuffer(validJson))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
