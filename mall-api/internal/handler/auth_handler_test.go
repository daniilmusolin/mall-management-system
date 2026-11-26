package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"mall-api/internal/domain"
	"mall-api/internal/handler"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")
	h := handler.NewAuthHandler(svc)

	mockRepo.On("GetByEmail", mock.Anything, "new@test.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	body, _ := json.Marshal(domain.RegisterInput{
		Email:    "new@test.com",
		Password: "password123",
		FullName: "New User",
	})

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")
	h := handler.NewAuthHandler(svc)

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           1,
		Email:        "test@test.com",
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
	}
	mockRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(user, nil)

	body, _ := json.Marshal(domain.LoginInput{
		Email:    "test@test.com",
		Password: "password123",
	})

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestAuthHandler_GetAllUsers(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")
	h := handler.NewAuthHandler(svc)

	users := []domain.User{{ID: 1, Email: "admin@test.com"}}
	mockRepo.On("GetAll", mock.Anything, mock.Anything).Return(users, nil)

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_UpdateUser(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")
	h := handler.NewAuthHandler(svc)

	user := &domain.User{ID: 1, FullName: "Old", Role: domain.RoleSupport}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(user, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	body, _ := json.Marshal(domain.UpdateUserInput{FullName: "New", Role: domain.RoleManager})
	req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBuffer(body))
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.UpdateUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_DeleteUser(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")
	h := handler.NewAuthHandler(svc)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/users/1", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	h.DeleteUser(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
