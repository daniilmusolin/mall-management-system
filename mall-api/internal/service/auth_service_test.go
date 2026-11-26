package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"mall-api/internal/domain"
	"mall-api/internal/port"
	"mall-api/internal/service"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	mockRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(nil, domain.ErrNotFound)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	req := domain.RegisterInput{
		Email:    "test@test.com",
		Password: "password123",
		FullName: "Test User",
	}

	resp, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "test@test.com", resp.User.Email)
	assert.Equal(t, domain.RoleSupport, resp.User.Role)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserExists(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	existingUser := &domain.User{Email: "test@test.com"}
	mockRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(existingUser, nil)

	req := domain.RegisterInput{
		Email:    "test@test.com",
		Password: "password123",
		FullName: "Test User",
	}

	_, err := svc.Register(context.Background(), req)

	assert.ErrorIs(t, err, domain.ErrConflict)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           1,
		Email:        "test@test.com",
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(user, nil)

	req := domain.LoginInput{
		Email:    "test@test.com",
		Password: "password123",
	}

	resp, err := svc.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "test@test.com", resp.User.Email)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		Email:        "test@test.com",
		PasswordHash: string(hash),
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@test.com").Return(user, nil)

	req := domain.LoginInput{
		Email:    "test@test.com",
		Password: "wrong",
	}

	_, err := svc.Login(context.Background(), req)
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthService_GetAllUsers(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	pagination := domain.Pagination{Limit: 10, Offset: 0}
	users := []domain.User{{ID: 1, Email: "admin@test.com"}}

	mockRepo.On("GetAll", mock.Anything, pagination).Return(users, nil)

	res, err := svc.GetAllUsers(context.Background(), pagination)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateUser_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	user := &domain.User{ID: 1, Email: "user@test.com", FullName: "Old Name", Role: domain.RoleSupport}
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(user, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.FullName == "New Name" && u.Role == domain.RoleManager
	})).Return(nil)

	in := domain.UpdateUserInput{FullName: "New Name", Role: domain.RoleManager}
	err := svc.UpdateUser(context.Background(), 1, in)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_DeleteUser_Success(t *testing.T) {
	mockRepo := new(port.MockUserRepo)
	svc := service.NewAuthService(mockRepo, "secret")

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	err := svc.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
