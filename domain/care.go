package domain

import (
	"fmt"
	"time"

	"github.com/mathehluiz/plant-care-tracker/internal/errs"
)

type Care struct {
	Id          int64     `json:"id"`
	PlantId     int64     `json:"plantId"`
	UserId      int64     `json:"-"`
	LastCare    time.Time `json:"lastCare"`
	NextCare    time.Time `json:"nextCare"`
	Name        string    `json:"name"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewCare(plantId, userId int64, nextCare time.Time, name, notes string) (*Care, error) {
	if len(name) < 3 || len(name) > 100 {
		return nil, errs.ErrInvalidCareName
	}

	if len(notes) < 3 || len(notes) > 1000 {
		return nil, errs.ErrInvalidCareNotes
	}
	lastCare := time.Now()
	if lastCare.After(nextCare) {
		fmt.Println("nextCare", nextCare)
		fmt.Println("lastCare", lastCare)
		return nil, errs.ErrInvalidCareDate
	}

	return &Care{
		PlantId:   plantId,
		UserId:    userId,
		LastCare:  lastCare,
		NextCare:  nextCare,
		Name:      name,
		Notes:     notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (c *Care) Update(plantId, userId int64, lastCare, nextCare time.Time, name, notes string) error {
	if len(name) < 3 || len(name) > 100 {
		return errs.ErrInvalidCareName
	}

	if len(notes) < 3 || len(notes) > 1000 {
		return errs.ErrInvalidCareNotes
	}

	if lastCare.After(nextCare) {
		return errs.ErrInvalidCareDate
	}

	c.PlantId = plantId
	c.UserId = userId
	c.LastCare = lastCare
	c.NextCare = nextCare
	c.Name = name
	c.Notes = notes
	c.UpdatedAt = time.Now()

	return nil
}