package user

import (
	_ "embed"

	"github.com/alextanhongpin/errors"
	"github.com/alextanhongpin/errors/examples/app"
)

const (
	MinAge = 13
	MaxAge = 150
)

//go:embed errors.json
var errorCodes []byte

var (
	// Load error codes.
	_                   = app.Errors.MustLoad(errorCodes)
	ErrNotFound         = app.Errors.Code("user.notFound")                                           // For text-only errors without params.
	ErrInvalidAge       = errors.Build(app.Errors.Code("user.invalidAge"), InvalidAgeParams{MaxAge}) // For errors with constant params.
	ErrUnderAge         = errors.Build(app.Errors.Code("user.underAge"), UnderAgeParams{MinAge})     //
	ErrInvalidName      = errors.Partial[InvalidNameParams](app.Errors.Code("user.invalidName"))     // For errors with dynamic params.
	ErrValidationErrors = errors.Partial[ValidationErrors](app.Errors.Code("user.validationErrors")) //
)

type InvalidAgeParams struct {
	MaxAge int64 `json:"maxAge"`
}

type UnderAgeParams struct {
	MinAge int64 `json:"minAge"`
}

type InvalidNameParams struct {
	Name string `json:"name"`
}

type ValidationErrors []ValidationFieldError

type ValidationFieldError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}
