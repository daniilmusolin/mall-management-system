package mysql

import (
	"context"
	"database/sql"
	"mall-api/internal/port"
)

type txManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) port.TransactionManager {
	return &txManager{db: db}
}

func (m *txManager) WithTransaction(ctx context.Context, fn func(tx interface{}) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}