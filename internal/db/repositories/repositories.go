package repositories

import (
	"context"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
	"github.com/jmoiron/sqlx"
)

func RunUpdateExec(ctx context.Context, db *sqlx.DB, query string, args ...any) error {
	changes, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := changes.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrNoRowsAffected
	}

	return nil
}
