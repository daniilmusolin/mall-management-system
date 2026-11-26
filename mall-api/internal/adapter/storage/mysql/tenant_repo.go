package mysql

import (
	"context"
	"database/sql"

	"mall-api/internal/domain"
)

type tenantRepo struct{ db *sql.DB }

func NewTenantRepo(db *sql.DB) *tenantRepo { return &tenantRepo{db: db} }

func (r *tenantRepo) Create(ctx context.Context, t *domain.Tenant) error {
	query := `INSERT INTO Tenants (Name, INN, LegalAddress, ContactPhone, Email) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.Name, t.INN, t.LegalAddress, t.ContactPhone, t.Email)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (r *tenantRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Tenant, error) {
	query := `SELECT TenantId, Name, INN, LegalAddress, ContactPhone, Email FROM Tenants ORDER BY TenantId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tenants []domain.Tenant
	for rows.Next() {
		var t domain.Tenant
		if err := rows.Scan(&t.ID, &t.Name, &t.INN, &t.LegalAddress, &t.ContactPhone, &t.Email); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *tenantRepo) GetByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	query := `SELECT TenantId, Name, INN, LegalAddress, ContactPhone, Email FROM Tenants WHERE TenantId = ?`
	var t domain.Tenant
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Name, &t.INN, &t.LegalAddress, &t.ContactPhone, &t.Email)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepo) ExistsByINN(ctx context.Context, inn string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Tenants WHERE INN = ?`, inn).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tenantRepo) Update(ctx context.Context, t *domain.Tenant) error {
	query := `UPDATE Tenants SET Name = ?, INN = ?, LegalAddress = ?, ContactPhone = ?, Email = ? WHERE TenantId = ?`
	res, err := r.db.ExecContext(ctx, query, t.Name, t.INN, t.LegalAddress, t.ContactPhone, t.Email, t.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *tenantRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM Tenants WHERE TenantId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
