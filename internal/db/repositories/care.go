package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/db/models"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
)

var _ = (domain.CareStorer)((*careRepository)(nil))

type careRepository struct {
	db *sqlx.DB
}

func NewCareRepository(db *sqlx.DB) *careRepository {
	return &careRepository{db}
}

func (r *careRepository) CreateCare(ctx context.Context, care *domain.Care) (int64, error) {
	query := `INSERT INTO cares (plant_id, user_id, last_care, next_care, name, notes, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, care.PlantId, care.UserId, care.LastCare, care.NextCare, care.Name, care.Notes, care.CreatedAt, care.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *careRepository) GetPlantCares(ctx context.Context, plantId int64) ([]*domain.Care, error) {
	query := `SELECT id, plant_id, user_id, last_care, next_care, name, notes, created_at, updated_at
	FROM cares WHERE plant_id = $1`

	var cares []*models.PGCare
	err := r.db.SelectContext(ctx, &cares, query, plantId)
	if err != nil {
		return nil, err
	}

	return models.PGCaresToDomainCares(cares), nil
}

func (r *careRepository) GetCareByID(ctx context.Context, id int64) (*domain.Care, error) {
	query := `SELECT id, plant_id, user_id, last_care, next_care, name, notes, created_at, updated_at
	FROM cares WHERE id = $1`

	var care[] models.PGCare
	if err := r.db.SelectContext(ctx, &care, query, id); err != nil {
		return nil, err
	}

	if len(care) == 0 {
		return nil, errs.ErrSelectNotMatch
	}

	if len(care) > 1 {
		return nil, errs.ErtSelectMultipleMatch
	}

	return &domain.Care{
		Id:       care[0].Id,
		PlantId:  care[0].PlantId,
		UserId:   care[0].UserId,
		LastCare: care[0].LastCare,
		NextCare: care[0].NextCare,
		Name:     care[0].Name,
		Notes:    care[0].Notes,
		CreatedAt: care[0].CreatedAt,
		UpdatedAt: care[0].UpdatedAt,
	}, nil
}

func (r *careRepository) UpdateCare(ctx context.Context, care *domain.Care) error {
	query := `UPDATE cares SET plant_id = $1, user_id = $2, last_care = $3, next_care = $4, name = $5, notes = $6, updated_at = $7
	WHERE id = $8`

	return RunUpdateExec(ctx, r.db, query, care.PlantId, care.UserId, care.LastCare, care.NextCare, care.Name, care.Notes, care.UpdatedAt, care.Id)
}

func (r *careRepository) DeleteCare(ctx context.Context, id int64) error {
	query := `DELETE FROM cares WHERE id = $1`

	return RunUpdateExec(ctx, r.db, query, id)
}