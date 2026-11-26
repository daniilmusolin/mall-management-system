package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type AuthService struct {
	repo   port.UserRepository
	secret []byte
}

func NewAuthService(repo port.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, secret: []byte(secret)}
}

func (s *AuthService) Register(ctx context.Context, in domain.RegisterInput) (*domain.AuthResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByEmail(ctx, in.Email)
	if err == nil && existing != nil {
		return nil, domain.ErrConflict
	}
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        in.Email,
		PasswordHash: string(hash),
		FullName:     in.FullName,
		Role:         domain.RoleSupport,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.generateTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, in domain.LoginInput) (*domain.AuthResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	user, err := s.repo.GetByEmail(ctx, in.Email)
	if err != nil || user == nil {
		return nil, domain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}
	return s.generateTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	userID, err := s.repo.GetUserIDByRefreshToken(ctx, refreshToken)
	if err != nil || userID == 0 {
		return nil, domain.ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, domain.ErrUnauthorized
	}

	_ = s.repo.DeleteRefreshToken(ctx, refreshToken)
	return s.generateTokens(ctx, user)
}

func (s *AuthService) GetAllUsers(ctx context.Context, pagination domain.Pagination) ([]domain.User, error) {
	return s.repo.GetAll(ctx, pagination)
}

func (s *AuthService) UpdateUser(ctx context.Context, id int64, in domain.UpdateUserInput) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	v := domain.NewValidator()
	user.FullName = domain.RequireString(v, "full_name", in.FullName, 2, 255)
	user.Role = domain.Role(domain.RequireEnum(v, "role", string(in.Role), "admin", "manager", "accountant", "technician", "support"))

	if err := v.Err(); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *AuthService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *AuthService) generateTokens(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	accessClaims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	refreshBytes := make([]byte, 32)
	_, _ = rand.Read(refreshBytes)
	refreshStr := hex.EncodeToString(refreshBytes)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	if err := s.repo.SaveRefreshToken(ctx, user.ID, refreshStr, expiresAt); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:        accessStr,
		RefreshToken: refreshStr,
		User:         *user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
}
