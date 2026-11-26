package mysql

import (
	"context"
	"database/sql"

	"mall-api/internal/domain"
)

type requestRepo struct{ db *sql.DB }

func NewRequestRepo(db *sql.DB) *requestRepo { return &requestRepo{db: db} }

func (r *requestRepo) Create(ctx context.Context, s *domain.ServiceRequest) error {
	query := `INSERT INTO ServiceRequests (PremiseId, RequestType, Description, Priority, Status) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.PremiseID, s.RequestType, s.Description, s.Priority, s.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (r *requestRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.ServiceRequest, error) {
	query := `SELECT RequestId, PremiseId, RequestType, Description, Priority, Status, AssignedTo, CreatedDate, CompletedDate FROM ServiceRequests ORDER BY RequestId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []domain.ServiceRequest
	for rows.Next() {
		var s domain.ServiceRequest
		if err := rows.Scan(&s.ID, &s.PremiseID, &s.RequestType, &s.Description, &s.Priority, &s.Status, &s.AssignedTo, &s.CreatedDate, &s.CompletedDate); err != nil {
			return nil, err
		}
		requests = append(requests, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *requestRepo) GetByID(ctx context.Context, id int64) (*domain.ServiceRequest, error) {
	query := `SELECT RequestId, PremiseId, RequestType, Description, Priority, Status, AssignedTo, CreatedDate, CompletedDate FROM ServiceRequests WHERE RequestId = ?`
	var s domain.ServiceRequest
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.PremiseID, &s.RequestType, &s.Description, &s.Priority, &s.Status, &s.AssignedTo, &s.CreatedDate, &s.CompletedDate)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *requestRepo) Update(ctx context.Context, s *domain.ServiceRequest) error {
	query := `UPDATE ServiceRequests SET PremiseId = ?, RequestType = ?, Description = ?, Priority = ?, Status = ?, AssignedTo = ?, CompletedDate = ? WHERE RequestId = ?`
	res, err := r.db.ExecContext(ctx, query, s.PremiseID, s.RequestType, s.Description, s.Priority, s.Status, s.AssignedTo, s.CompletedDate, s.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *requestRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM ServiceRequests WHERE RequestId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
