package domain

import (
	"time"

	"github.com/mathehluiz/plant-care-tracker/internal/errs"
)

type Plant struct {
	Id               int64               `json:"id"`
	Name             string              `json:"name"`
	AcquisitionDate  time.Time           `json:"acquisitionDate"`
	Location         string              `json:"location"`
	CareFrequency    int                 `json:"careFrequency"`
	UserId           int64               `json:"userId"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}

func NewPlant(name, location string, acquisitionDate time.Time, careFrequency int, userId int64) (*Plant, error) {
	if len(name) < 3 || len(name) > 100 {
		return nil, errs.ErrInvalidPlantName
	}

	if len(location) < 3 || len(location) > 100 {
		return nil, errs.ErrInvalidPlantLocation
	}

	if careFrequency < 1 || careFrequency > 365 {
		return nil, errs.ErrInvalidPlantCareFrequency
	}

	return &Plant{
		Name:            name,
		Location:        location,
		AcquisitionDate: acquisitionDate,
		CareFrequency:   careFrequency,
		UserId:          userId,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func (p *Plant) Update(name, location string, acquisitionDate time.Time, careFrequency int) error {
	if len(name) < 3 || len(name) > 100 {
		return errs.ErrInvalidPlantName
	}

	if len(location) < 3 || len(location) > 100 {
		return errs.ErrInvalidPlantLocation
	}

	if careFrequency < 1 || careFrequency > 365 {
		return errs.ErrInvalidPlantCareFrequency
	}

	p.Name = name
	p.Location = location
	p.AcquisitionDate = acquisitionDate
	p.CareFrequency = careFrequency
	p.UpdatedAt = time.Now()

	return nil
}