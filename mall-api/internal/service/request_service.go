package service

import (
	"context"
	"time"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type RequestService struct {
	repo        port.RequestRepository
	premiseRepo port.PremiseRepository
}

func NewRequestService(r port.RequestRepository, p port.PremiseRepository) *RequestService {
	return &RequestService{repo: r, premiseRepo: p}
}

func (s *RequestService) Create(ctx context.Context, sr *domain.ServiceRequest) error {
	if err := sr.ValidateForCreate(); err != nil {
		return err
	}
	if _, err := s.premiseRepo.GetByID(ctx, sr.PremiseID); err != nil {
		return domain.NewValidationError("premise_id", "помещение не найдено")
	}
	sr.Status = "open"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, sr)
}

func (s *RequestService) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.ServiceRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetAll(ctx, pagination)
}

func (s *RequestService) Update(ctx context.Context, sr *domain.ServiceRequest) error {
	if err := sr.ValidateForUpdate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, sr)
}

func (s *RequestService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
