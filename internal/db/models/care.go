package models

import (
	"database/sql"
	"time"

	"github.com/mathehluiz/plant-care-tracker/domain"
)

type PGCare struct {
	Id          int64     `db:"id"`
	PlantId     int64     `db:"plant_id"`
	UserId      int64     `db:"user_id"`
	LastCare    time.Time `db:"last_care"`
	NextCare    time.Time `db:"next_care"`
	Name        string    `db:"name"`
	Notes       string    `db:"notes"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	Plant       *PGPlant `db:"-"`
	User 	  	*PGUser `db:"-"`
}

func PGCaresToDomainCares(cares []*PGCare) []*domain.Care {
	var domainCares []*domain.Care
	if len(cares) == 0 {
		return []*domain.Care{}
	}
	for _, care := range cares {
		domainCares = append(domainCares, &domain.Care{
			Id:        care.Id,
			PlantId:   care.PlantId,
			UserId:    care.UserId,
			LastCare:  care.LastCare,
			NextCare:  care.NextCare,
			Name:      care.Name,
			Notes:     care.Notes,
			CreatedAt: care.CreatedAt,
			UpdatedAt: care.UpdatedAt,
		})
	}
	return domainCares
}