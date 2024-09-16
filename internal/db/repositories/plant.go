package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/db/models"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
)

var _ = (domain.PlantStorer)((*plantRepository)(nil))

type plantRepository struct {
	db *sqlx.DB
}

func NewPlantRepository(db *sqlx.DB) *plantRepository {
	return &plantRepository{db}
}

func (p *plantRepository) CreatePlant(ctx context.Context, plant *domain.Plant) (int64, error) {
	insertQuery := `INSERT INTO plants (name, acquisition_date, location, care_frequency, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	var id int64
	err := p.db.QueryRowxContext(ctx, insertQuery, plant.Name, plant.AcquisitionDate, plant.Location, plant.CareFrequency, plant.UserId, plant.CreatedAt, plant.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *plantRepository) GetPlantByID(ctx context.Context, id int64) (*domain.Plant, error) {
	selectQuery := `SELECT id, name, acquisition_date, location, care_frequency, user_id, created_at, updated_at, deleted_at
		FROM plants WHERE id = $1;`

	var plant []models.PGPlant
	if err := p.db.SelectContext(ctx, &plant, selectQuery, id); err != nil {
		return nil, err
	}

	if len(plant) == 0 {
		return nil, errs.ErrSelectNotMatch
	}

	if len(plant) > 1 {
		return nil, errs.ErtSelectMultipleMatch
	}

	return &domain.Plant{
		Id:              plant[0].Id,
		Name:            plant[0].Name,
		AcquisitionDate: plant[0].AcquisitionDate,
		Location:        plant[0].Location,
		CareFrequency:   plant[0].CareFrequency,
		UserId:          plant[0].UserId,
		CreatedAt:       plant[0].CreatedAt,
		UpdatedAt:       plant[0].UpdatedAt,
	}, nil
}

func (p *plantRepository) GetPlantsByUserID(ctx context.Context, userID int64) ([]*domain.Plant, error) {
	selectQuery := `SELECT id, name, acquisition_date, location, care_frequency, user_id, created_at, updated_at, deleted_at
		FROM plants WHERE user_id = $1;`

	var plants []models.PGPlant
	if err := p.db.SelectContext(ctx, &plants, selectQuery, userID); err != nil {
		return nil, err
	}

	var domainPlants []*domain.Plant
	for _, plant := range plants {
		domainPlants = append(domainPlants, &domain.Plant{
			Id:              plant.Id,
			Name:            plant.Name,
			AcquisitionDate: plant.AcquisitionDate,
			Location:        plant.Location,
			CareFrequency:   plant.CareFrequency,
			UserId:          plant.UserId,
			CreatedAt:       plant.CreatedAt,
			UpdatedAt:       plant.UpdatedAt,
		})
	}

	return domainPlants, nil
}

func (p *plantRepository) UpdatePlant(ctx context.Context, plant *domain.Plant) error {
	updateQuery := `UPDATE plants SET name = $1, acquisition_date = $2, location = $3, care_frequency = $4, updated_at = $5
		WHERE id = $6;`

	return RunUpdateExec(ctx, p.db, updateQuery, plant.Name, plant.AcquisitionDate, plant.Location, plant.CareFrequency, plant.UpdatedAt, plant.Id)
}

func (p *plantRepository) DeletePlant(ctx context.Context, id int64) error {
	deleteQuery := `UPDATE plants SET deleted_at = $1 WHERE id = $2;`
    now := time.Now()
	
	return RunUpdateExec(ctx, p.db, deleteQuery, now, id)
}