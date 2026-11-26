package service

import (
	"context"
	"time"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type TenantService struct {
	repo port.TenantRepository
}

func NewTenantService(r port.TenantRepository) *TenantService {
	return &TenantService{repo: r}
}

func (s *TenantService) Create(ctx context.Context, t *domain.Tenant) error {
	if err := t.ValidateForCreate(); err != nil {
		return err
	}
	exists, err := s.repo.ExistsByINN(ctx, t.INN)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewValidationError("inn", "арендатор с таким ИНН уже существует")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, t)
}

func (s *TenantService) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Tenant, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetAll(ctx, pagination)
}

func (s *TenantService) Update(ctx context.Context, t *domain.Tenant) error {
	if err := t.ValidateForUpdate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, t)
}

func (s *TenantService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
