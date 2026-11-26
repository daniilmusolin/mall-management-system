package mysql

import (
	"context"
	"database/sql"
	"time"

	"mall-api/internal/domain"
)

type paymentRepo struct{ db *sql.DB }

func NewPaymentRepo(db *sql.DB) *paymentRepo { return &paymentRepo{db: db} }

func (r *paymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO Payments (ContractId, Amount, PaymentDate, PaymentType, Status, ReferenceNumber) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.ContractID, p.Amount, p.PaymentDate, p.PaymentType, p.Status, p.ReferenceNumber)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

func (r *paymentRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Payment, error) {
	query := `SELECT PaymentId, ContractId, Amount, PaymentDate, PaymentType, Status, ReferenceNumber FROM Payments ORDER BY PaymentId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		var paymentDate time.Time
		if err := rows.Scan(&p.ID, &p.ContractID, &p.Amount, &paymentDate, &p.PaymentType, &p.Status, &p.ReferenceNumber); err != nil {
			return nil, err
		}
		p.PaymentDate = paymentDate.Format("2006-01-02")
		payments = append(payments, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepo) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `SELECT PaymentId, ContractId, Amount, PaymentDate, PaymentType, Status, ReferenceNumber FROM Payments WHERE PaymentId = ?`
	var p domain.Payment
	var paymentDate time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.ContractID, &p.Amount, &paymentDate, &p.PaymentType, &p.Status, &p.ReferenceNumber)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.PaymentDate = paymentDate.Format("2006-01-02")
	return &p, nil
}

func (r *paymentRepo) Update(ctx context.Context, p *domain.Payment) error {
	query := `UPDATE Payments SET ContractId = ?, Amount = ?, PaymentDate = ?, PaymentType = ?, Status = ?, ReferenceNumber = ? WHERE PaymentId = ?`
	res, err := r.db.ExecContext(ctx, query, p.ContractID, p.Amount, p.PaymentDate, p.PaymentType, p.Status, p.ReferenceNumber, p.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *paymentRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM Payments WHERE PaymentId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
