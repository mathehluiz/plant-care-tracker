package domain

import "context"

type PlantStorer interface {
	CreatePlant(ctx context.Context, plant *Plant) (int64, error)
	GetPlantByID(ctx context.Context, id int64) (*Plant, error)
	GetPlantsByUserID(ctx context.Context, userID int64) ([]*Plant, error)
	UpdatePlant(ctx context.Context, plant *Plant) error
	DeletePlant(ctx context.Context, id int64) error
}