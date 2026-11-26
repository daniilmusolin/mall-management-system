package service

import (
	"context"
	"time"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type PaymentService struct {
	repo         port.PaymentRepository
	contractRepo port.ContractRepository
}

func NewPaymentService(r port.PaymentRepository, c port.ContractRepository) *PaymentService {
	return &PaymentService{repo: r, contractRepo: c}
}

func (s *PaymentService) Create(ctx context.Context, p *domain.Payment) error {
	if err := p.ValidateForCreate(); err != nil {
		return err
	}
	contract, err := s.contractRepo.GetByID(ctx, p.ContractID)
	if err != nil {
		return domain.NewValidationError("contract_id", "договор не найден")
	}
	if contract.Status != "active" {
		return domain.NewValidationError("contract_id", "нельзя создать платёж по неактивному договору")
	}
	p.Status = "completed"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, p)
}

func (s *PaymentService) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetAll(ctx, pagination)
}

func (s *PaymentService) Update(ctx context.Context, p *domain.Payment) error {
	if err := p.ValidateForUpdate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, p)
}

func (s *PaymentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
