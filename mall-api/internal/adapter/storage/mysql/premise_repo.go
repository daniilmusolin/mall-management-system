package mysql

import (
	"context"
	"database/sql"

	"mall-api/internal/domain"
)

type premiseRepo struct{ db *sql.DB }

func NewPremiseRepo(db *sql.DB) *premiseRepo { return &premiseRepo{db: db} }

func (r *premiseRepo) Create(ctx context.Context, p *domain.Premise) error {
	query := `INSERT INTO Premises (RoomNumber, Floor, Area, Type, Zone, Status) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.RoomNumber, p.Floor, p.Area, p.Type, p.Zone, p.Status)
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

func (r *premiseRepo) CreateTx(ctx context.Context, tx interface{}, p *domain.Premise) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrConflict
	}
	query := `INSERT INTO Premises (RoomNumber, Floor, Area, Type, Zone, Status) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := sqlTx.ExecContext(ctx, query, p.RoomNumber, p.Floor, p.Area, p.Type, p.Zone, p.Status)
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

func (r *premiseRepo) GetAll(ctx context.Context, pagination domain.Pagination) ([]domain.Premise, error) {
	query := `SELECT PremiseId, RoomNumber, Floor, Area, Type, Zone, Status FROM Premises ORDER BY PremiseId DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var premises []domain.Premise
	for rows.Next() {
		var p domain.Premise
		if err := rows.Scan(&p.ID, &p.RoomNumber, &p.Floor, &p.Area, &p.Type, &p.Zone, &p.Status); err != nil {
			return nil, err
		}
		premises = append(premises, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return premises, nil
}

func (r *premiseRepo) GetByID(ctx context.Context, id int64) (*domain.Premise, error) {
	query := `SELECT PremiseId, RoomNumber, Floor, Area, Type, Zone, Status FROM Premises WHERE PremiseId = ?`
	var p domain.Premise
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.RoomNumber, &p.Floor, &p.Area, &p.Type, &p.Zone, &p.Status)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *premiseRepo) ExistsByRoomNumber(ctx context.Context, roomNumber string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM Premises WHERE RoomNumber = ?`, roomNumber).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *premiseRepo) UpdateStatusTx(ctx context.Context, tx interface{}, id int64, status string) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrConflict
	}
	_, err := sqlTx.ExecContext(ctx, `UPDATE Premises SET Status = ? WHERE PremiseId = ?`, status, id)
	return err
}

func (r *premiseRepo) Update(ctx context.Context, p *domain.Premise) error {
	query := `UPDATE Premises SET RoomNumber = ?, Floor = ?, Area = ?, Type = ?, Zone = ?, Status = ? WHERE PremiseId = ?`
	res, err := r.db.ExecContext(ctx, query, p.RoomNumber, p.Floor, p.Area, p.Type, p.Zone, p.Status, p.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *premiseRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM Premises WHERE PremiseId = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
