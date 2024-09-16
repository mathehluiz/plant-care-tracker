package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mathehluiz/plant-care-tracker/domain"
	"github.com/mathehluiz/plant-care-tracker/internal/db/drivers"
	"github.com/mathehluiz/plant-care-tracker/internal/db/models"
	"github.com/mathehluiz/plant-care-tracker/internal/errs"
)

var _ = (domain.UserStorer)((*userRepository)(nil))

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db}
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	selectQuery := `SELECT id, external_id, email, username, password, roles, active, verified
		FROM users WHERE username = $1;`

	return u.getUser(ctx, selectQuery, username)
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	selectQuery := `SELECT id, external_id, email, username, password, roles, active, verified
		FROM users WHERE email = $1;`

	return u.getUser(ctx, selectQuery, email)
}

func (u *userRepository) getUser(ctx context.Context, query string, args string) (*domain.User, error) {
	var user []models.PGUser

	if err := u.db.SelectContext(ctx, &user, query, args); err != nil {
		return nil, err
	}

	if len(user) == 0 {
		return nil, errs.ErrSelectNotMatch
	}

	if len(user) > 1 {
		return nil, errs.ErtSelectMultipleMatch
	}

	return &domain.User{
		Id:         user[0].Id,
		ExternalId: user[0].ExternalId,
		Email:      user[0].Email,
		Username:   user[0].Username,
		Password:   user[0].Password,
		Roles:      user[0].Roles,
		Active:     user[0].Active,
		Verified:   user[0].Verified,
	}, nil
}

func (u *userRepository) GetUserByExternalId(ctx context.Context, id string) (*domain.User, error) {
	selectQuery := `SELECT id, external_id, email, username, password, roles, active, verified
		FROM users WHERE external_id = $1;`

	return u.getUser(ctx, selectQuery, id)
}

func (u *userRepository) AddRolesToUser(ctx context.Context, id string, roles []string) error {
	user, err := u.GetUserByExternalId(ctx, id)
	if err != nil {
		return err
	}

	newRoles := make([]string, 0)

	for _, role := range roles {
		exists := false
		for _, r := range user.Roles {
			if role == r {
				exists = true
			}
		}
		if !exists {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) == 0 {
		return nil
	}

	updateQuery := `UPDATE users SET roles = $1 WHERE external_id = $2;`

	return RunUpdateExec(ctx, u.db, updateQuery, drivers.StringArray(append(user.Roles, newRoles...)), id)
}

func (u *userRepository) RemoveRolesFromUser(ctx context.Context, id string, roles []string) error {
	user, err := u.GetUserByExternalId(ctx, id)
	if err != nil {
		return err
	}

	newRoles := make([]string, 0)

	for _, role := range user.Roles {
		for _, r := range roles {
			if role != r {
				newRoles = append(newRoles, role)
			}
		}
	}

	updateQuery := `UPDATE users SET roles = $1 WHERE external_id = $2;`

	return RunUpdateExec(ctx, u.db, updateQuery, drivers.StringArray(newRoles), id)
}

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	deleteQuery := `UPDATE users SET deleted_at = NOW(), active = false WHERE external_id = $1;`

	return RunUpdateExec(ctx, u.db, deleteQuery, id)
}

func (u *userRepository) UpdatePassword(ctx context.Context, id, password string) error {
	updateQuery := `UPDATE users SET password = $1 WHERE external_id = $2;`

	return RunUpdateExec(ctx, u.db, updateQuery, password, id)
}

func (u *userRepository) VerifyUser(ctx context.Context, id string) error {
	updateQuery := `UPDATE users SET verified = true WHERE external_id = $1;`

	return RunUpdateExec(ctx, u.db, updateQuery, id)
}

func (u *userRepository) UpdateActiveUserStatus(ctx context.Context, id string, active bool) error {
	updateQuery := `UPDATE users SET active = $1 WHERE external_id = $2;`

	return RunUpdateExec(ctx, u.db, updateQuery, active, id)
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	insertQuery := `
		INSERT INTO users (email, username, password, roles)
		VALUES ($1, $2, $3, $4) RETURNING external_id;
	`

	var id string
	if err := u.db.QueryRowContext(ctx, insertQuery,
		user.Email, user.Username, user.Password, drivers.StringArray(user.Roles)).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
