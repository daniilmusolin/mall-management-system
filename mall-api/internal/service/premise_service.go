package service

import (
	"context"
	"time"

	"mall-api/internal/domain"
	"mall-api/internal/port"
)

type PremiseService struct {
	repo port.PremiseRepository
}

func NewPremiseService(r port.PremiseRepository) *PremiseService {
	return &PremiseService{repo: r}
}

func (s *PremiseService) Create(ctx context.Context, p *domain.Premise) error {
	if err := p.ValidateForCreate(); err != nil {
		return err
	}
	exists, err := s.repo.ExistsByRoomNumber(ctx, p.RoomNumber)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewValidationError("room_number", "помещение с таким номером уже существует")
	}
	p.Status = "free"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.Create(ctx, p)
}

func (s *PremiseService) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Premise, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.repo.GetAll(ctx, pagination)
}

func (s *PremiseService) Update(ctx context.Context, p *domain.Premise) error {
	if err := p.ValidateForUpdate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, p)
}

func (s *PremiseService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
