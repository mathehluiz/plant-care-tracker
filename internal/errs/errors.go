package errs

import "errors"

var (
	ErrNoRowsAffected      = errors.New("no rows affected")
	ErrSelectNotMatch      = errors.New("select query did not match any rows")
	ErtSelectMultipleMatch = errors.New("select query matched multiple rows")

	ErrInvalidBody     = errors.New("invalid body provided")
	ErrInvalidUsername = errors.New("invalid username provided")
	ErrInvalidPassword = errors.New("invalid password provided")
	ErrInvalidEmail    = errors.New("invalid email provided")

	ErrInvalidCode = errors.New("invalid code provided")

	ErrNotFound        = errors.New("not found")
	ErrAlreadyVerified = errors.New("user is already verified")
	ErrCodeExpired     = errors.New("code expired! New code sent to email")

	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrInvalidPlantName         = errors.New("invalid plant name provided")
	ErrInvalidPlantLocation     = errors.New("invalid plant location provided")
	ErrInvalidPlantCareFrequency = errors.New("invalid plant care frequency provided")

	ErrInvalidCareName  = errors.New("invalid care name provided")
	ErrInvalidCareNotes = errors.New("invalid care notes provided")
	ErrInvalidCareDate  = errors.New("invalid care date provided")
)
