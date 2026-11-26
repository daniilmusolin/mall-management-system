package port

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"mall-api/internal/domain"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.User, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserIDByRefreshToken(ctx context.Context, token string) (int64, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

type MockPremiseRepo struct{ mock.Mock }

func (m *MockPremiseRepo) Create(ctx context.Context, premise *domain.Premise) error {
	args := m.Called(ctx, premise)
	return args.Error(0)
}

func (m *MockPremiseRepo) CreateTx(ctx context.Context, tx interface{}, premise *domain.Premise) error {
	args := m.Called(ctx, tx, premise)
	return args.Error(0)
}

func (m *MockPremiseRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Premise, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Premise), args.Error(1)
}

func (m *MockPremiseRepo) GetByID(ctx context.Context, id int64) (*domain.Premise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Premise), args.Error(1)
}

func (m *MockPremiseRepo) ExistsByRoomNumber(ctx context.Context, roomNumber string) (bool, error) {
	args := m.Called(ctx, roomNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockPremiseRepo) UpdateStatusTx(ctx context.Context, tx interface{}, id int64, status string) error {
	args := m.Called(ctx, tx, id, status)
	return args.Error(0)
}

func (m *MockPremiseRepo) Update(ctx context.Context, premise *domain.Premise) error {
	args := m.Called(ctx, premise)
	return args.Error(0)
}

func (m *MockPremiseRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTenantRepo struct{ mock.Mock }

func (m *MockTenantRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Tenant, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) GetByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tenant), args.Error(1)
}

func (m *MockTenantRepo) ExistsByINN(ctx context.Context, inn string) (bool, error) {
	args := m.Called(ctx, inn)
	return args.Bool(0), args.Error(1)
}

func (m *MockTenantRepo) Update(ctx context.Context, tenant *domain.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockContractRepo struct{ mock.Mock }

func (m *MockContractRepo) CreateTx(ctx context.Context, tx interface{}, contract *domain.LeaseContract) error {
	args := m.Called(ctx, tx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.LeaseContract, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LeaseContract), args.Error(1)
}

func (m *MockContractRepo) GetByID(ctx context.Context, id int64) (*domain.LeaseContract, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LeaseContract), args.Error(1)
}

func (m *MockContractRepo) Update(ctx context.Context, contract *domain.LeaseContract) error {
	args := m.Called(ctx, contract)
	return args.Error(0)
}

func (m *MockContractRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPaymentRepo struct{ mock.Mock }

func (m *MockPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Payment, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) Update(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRequestRepo struct{ mock.Mock }

func (m *MockRequestRepo) Create(ctx context.Context, request *domain.ServiceRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockRequestRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.ServiceRequest, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ServiceRequest), args.Error(1)
}

func (m *MockRequestRepo) GetByID(ctx context.Context, id int64) (*domain.ServiceRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ServiceRequest), args.Error(1)
}

func (m *MockRequestRepo) Update(ctx context.Context, request *domain.ServiceRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockRequestRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTxManager struct{ mock.Mock }

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(tx interface{}) error) error {
	return fn(mock.Anything)
}