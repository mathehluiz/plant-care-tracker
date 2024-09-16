package models

import (
	"database/sql"
	"github.com/mathehluiz/plant-care-tracker/internal/db/drivers"
	"time"
)

type PGUser struct {
	Id         int64               `db:"id"`
	Email      string              `db:"email"`
	Username   string              `db:"username"`
	Password   string              `db:"password"`
	Roles      drivers.StringArray `db:"roles"`
	Active     bool                `db:"active"`
	Verified   bool                `db:"verified"`
	ExternalId string              `db:"external_id"`
	CreatedAt  time.Time           `db:"created_at"`
	UpdatedAt  time.Time           `db:"updated_at"`
	DeletedAt  sql.NullTime        `db:"deleted_at"`
}
