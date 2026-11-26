package mysql

import (
	"context"
	"database/sql"
	"time"

	"mall-api/internal/domain"
)

type contractRepo struct{ db *sql.DB }

func NewContractRepo(db *sql.DB) *contractRepo { return &contractRepo{db: db} }

func (r *contractRepo) CreateTx(ctx context.Context, tx interface{}, c *domain.LeaseContract) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrConflict
	}
	query := `INSERT INTO LeaseContracts (PremiseId, TenantId, StartDate, EndDate, RatePerSqm, DepositAmount, Status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := sqlTx.ExecContext(ctx, query, c.PremiseID, c.TenantID, c.StartDate, c.EndDate, c.RatePerSqm, c.DepositAmount, c.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *contractRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.LeaseContract, error) {
	query := `SELECT ContractId, PremiseId, TenantId, StartDate, EndDate, RatePerSqm, DepositAmount, Status FROM LeaseContracts ORDER BY ContractId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contracts []domain.LeaseContract
	for rows.Next() {
		var c domain.LeaseContract
		var startDate, endDate time.Time
		if err := rows.Scan(&c.ID, &c.PremiseID, &c.TenantID, &startDate, &endDate, &c.RatePerSqm, &c.DepositAmount, &c.Status); err != nil {
			return nil, err
		}
		c.StartDate = startDate.Format("2006-01-02")
		c.EndDate = endDate.Format("2006-01-02")
		contracts = append(contracts, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *contractRepo) GetByID(ctx context.Context, id int64) (*domain.LeaseContract, error) {
	query := `SELECT ContractId, PremiseId, TenantId, StartDate, EndDate, RatePerSqm, DepositAmount, Status FROM LeaseContracts WHERE ContractId = ?`
	var c domain.LeaseContract
	var startDate, endDate time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.PremiseID, &c.TenantID, &startDate, &endDate, &c.RatePerSqm, &c.DepositAmount, &c.Status)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	c.StartDate = startDate.Format("2006-01-02")
	c.EndDate = endDate.Format("2006-01-02")
	return &c, nil
}

func (r *contractRepo) Update(ctx context.Context, c *domain.LeaseContract) error {
	query := `UPDATE LeaseContracts SET PremiseId = ?, TenantId = ?, StartDate = ?, EndDate = ?, RatePerSqm = ?, DepositAmount = ?, Status = ? WHERE ContractId = ?`
	res, err := r.db.ExecContext(ctx, query, c.PremiseID, c.TenantID, c.StartDate, c.EndDate, c.RatePerSqm, c.DepositAmount, c.Status, c.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *contractRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM LeaseContracts WHERE ContractId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
