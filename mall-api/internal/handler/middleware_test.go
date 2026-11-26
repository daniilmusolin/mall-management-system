package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"mall-api/internal/domain"
	"mall-api/internal/handler"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func generateTestToken(role domain.Role, secret string) string {
	claims := jwt.MapClaims{
		"sub":  float64(1),
		"role": string(role),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, _ := token.SignedString([]byte(secret))
	return str
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	authSvc := service.NewAuthService(mockRepo, "secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := handler.AuthMiddleware(authSvc)
	protectedHandler := authMiddleware(dummyHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	protectedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	authSvc := service.NewAuthService(mockRepo, "secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := handler.AuthMiddleware(authSvc)
	protectedHandler := authMiddleware(dummyHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	w := httptest.NewRecorder()

	protectedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	authSvc := service.NewAuthService(mockRepo, "secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authMiddleware := handler.AuthMiddleware(authSvc)
	protectedHandler := authMiddleware(dummyHandler)

	token := generateTestToken(domain.RoleAdmin, "secret")

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	protectedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	authSvc := service.NewAuthService(mockRepo, "secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rbacMiddleware := handler.RequireRole(domain.RoleAdmin)
	protectedHandler := handler.AuthMiddleware(authSvc)(rbacMiddleware(dummyHandler))

	token := generateTestToken(domain.RoleTechnician, "secret")

	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	protectedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRole_Allowed(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	authSvc := service.NewAuthService(mockRepo, "secret")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rbacMiddleware := handler.RequireRole(domain.RoleAdmin, domain.RoleManager)
	protectedHandler := handler.AuthMiddleware(authSvc)(rbacMiddleware(dummyHandler))

	token := generateTestToken(domain.RoleManager, "secret")

	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	protectedHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
