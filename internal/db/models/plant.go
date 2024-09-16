package models

import (
	"database/sql"
	"time"
)

type PGPlant struct {
	Id               int64               `db:"id"`
	Name             string              `db:"name"`
	AcquisitionDate  time.Time           `db:"acquisition_date"`
	Location         string              `db:"location"`
	CareFrequency    int                 `db:"care_frequency"`
	UserId           int64               `db:"user_id"`
	CreatedAt        time.Time           `db:"created_at"`
	UpdatedAt        time.Time           `db:"updated_at"`
	DeletedAt        sql.NullTime        `db:"deleted_at"`
	User             *PGUser             `db:"-"`
}