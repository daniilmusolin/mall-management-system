package service

import (
	"context"
	"time"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type ContractService struct {
	repo        port.ContractRepository
	premiseRepo port.PremiseRepository
	tenantRepo  port.TenantRepository
	txManager   port.TransactionManager
}

func NewContractService(r port.ContractRepository, p port.PremiseRepository, t port.TenantRepository, tx port.TransactionManager) *ContractService {
	return &ContractService{repo: r, premiseRepo: p, tenantRepo: t, txManager: tx}
}

func (s *ContractService) Create(ctx context.Context, c *domain.LeaseContract) error {
	if err := c.ValidateForCreate(); err != nil {
		return err
	}

	premise, err := s.premiseRepo.GetByID(ctx, c.PremiseID)
	if err != nil {
		return domain.NewValidationError("premise_id", "помещение не найдено")
	}
	if premise.Status == "rented" {
		return domain.NewValidationError("premise_id", "помещение уже сдано")
	}

	if _, err := s.tenantRepo.GetByID(ctx, c.TenantID); err != nil {
		return domain.NewValidationError("tenant_id", "арендатор не найден")
	}

	c.Status = "active"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.txManager.WithTransaction(ctx, func(tx interface{}) error {
		if err := s.repo.CreateTx(ctx, tx, c); err != nil {
			return err
		}
		return s.premiseRepo.UpdateStatusTx(ctx, tx, c.PremiseID, "rented")
	})
}

func (s *ContractService) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.LeaseContract, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetAll(ctx, pagination)
}

func (s *ContractService) Update(ctx context.Context, c *domain.LeaseContract) error {
	if err := c.ValidateForUpdate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, c)
}

func (s *ContractService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
