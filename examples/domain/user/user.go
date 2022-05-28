package user

import (
	_ "embed"

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
	_ = app.MustLoadError(errorCodes)

	// All errors.
	ErrNotFound         = app.NewError("user.notFound")                                        // For text-only errors without params.
	ErrInvalidAge       = app.NewFullErrorCustom("user.invalidAge", InvalidAgeParams{MaxAge})  // For errors with constant params.
	ErrUnderAge         = app.NewFullError("user.underAge", UnderAgeParams{MinAge})            //
	ErrInvalidName      = app.NewPartialError[InvalidNameParams]("user.invalidName")           // For errors with dynamic params.
	ErrValidationErrors = app.NewPartialError[ValidationErrorsParams]("user.validationErrors") //
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

type ValidationErrorsParams struct {
	Count       int64   `json:"count"`
	PluralError string  `json:"-"`
	Errors      []error `json:"errors"`
}
