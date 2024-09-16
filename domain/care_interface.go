package domain

import (
	"context"
)

type CareStorer interface {
	CreateCare(ctx context.Context, care *Care) (int64, error)
	GetPlantCares(ctx context.Context, plantId int64) ([]*Care, error)
	GetCareByID(ctx context.Context, id int64) (*Care, error)
	UpdateCare(ctx context.Context, care *Care) error
	DeleteCare(ctx context.Context, id int64) error
}