package domain

import "time"

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleManager    Role = "manager"
	RoleAccountant Role = "accountant"
	RoleTechnician Role = "technician"
	RoleSupport    Role = "support"
)

type Pagination struct {
	Limit  int
	Offset int
}

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	FullName string `json:"full_name"`
	Role     Role   `json:"role"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Premise struct {
	ID         int64   `json:"id"`
	RoomNumber string  `json:"room_number"`
	Floor      int     `json:"floor"`
	Area       float64 `json:"area"`
	Type       string  `json:"type"`
	Zone       string  `json:"zone"`
	Status     string  `json:"status"`
}

type Tenant struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	INN          string `json:"inn"`
	LegalAddress string `json:"legal_address"`
	ContactPhone string `json:"contact_phone"`
	Email        string `json:"email"`
}

type LeaseContract struct {
	ID            int64   `json:"id"`
	PremiseID     int64   `json:"premise_id"`
	TenantID      int64   `json:"tenant_id"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	RatePerSqm    float64 `json:"rate_per_sqm"`
	DepositAmount float64 `json:"deposit_amount"`
	Status        string  `json:"status"`
}

type Payment struct {
	ID              int64   `json:"id"`
	ContractID      int64   `json:"contract_id"`
	Amount          float64 `json:"amount"`
	PaymentDate     string  `json:"payment_date"`
	PaymentType     string  `json:"payment_type"`
	Status          string  `json:"status"`
	ReferenceNumber string  `json:"reference_number"`
}

type ServiceRequest struct {
	ID            int64   `json:"id"`
	PremiseID     int64   `json:"premise_id"`
	RequestType   string  `json:"request_type"`
	Description   string  `json:"description"`
	Priority      string  `json:"priority"`
	Status        string  `json:"status"`
	AssignedTo    *int64  `json:"assigned_to"`
	CreatedDate   string  `json:"created_date"`
	CompletedDate *string `json:"completed_date"`
}
