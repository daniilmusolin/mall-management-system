package port

import (
	"context"
	"time"

	"mall-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id int64) error

	SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	GetUserIDByRefreshToken(ctx context.Context, token string) (int64, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}

type PremiseRepository interface {
	Create(ctx context.Context, premise *domain.Premise) error
	CreateTx(ctx context.Context, tx interface{}, premise *domain.Premise) error
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Premise, error)
	GetByID(ctx context.Context, id int64) (*domain.Premise, error)
	ExistsByRoomNumber(ctx context.Context, roomNumber string) (bool, error)
	UpdateStatusTx(ctx context.Context, tx interface{}, id int64, status string) error
	Update(ctx context.Context, premise *domain.Premise) error
	Delete(ctx context.Context, id int64) error
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *domain.Tenant) error
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Tenant, error)
	GetByID(ctx context.Context, id int64) (*domain.Tenant, error)
	ExistsByINN(ctx context.Context, inn string) (bool, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	Delete(ctx context.Context, id int64) error
}

type ContractRepository interface {
	CreateTx(ctx context.Context, tx interface{}, contract *domain.LeaseContract) error
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.LeaseContract, error)
	GetByID(ctx context.Context, id int64) (*domain.LeaseContract, error)
	Update(ctx context.Context, contract *domain.LeaseContract) error
	Delete(ctx context.Context, id int64) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Payment, error)
	GetByID(ctx context.Context, id int64) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
	Delete(ctx context.Context, id int64) error
}

type RequestRepository interface {
	Create(ctx context.Context, request *domain.ServiceRequest) error
	GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.ServiceRequest, error)
	GetByID(ctx context.Context, id int64) (*domain.ServiceRequest, error)
	Update(ctx context.Context, request *domain.ServiceRequest) error
	Delete(ctx context.Context, id int64) error
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx interface{}) error) error
}
