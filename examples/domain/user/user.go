package user

import (
	_ "embed"

	"github.com/alextanhongpin/errors/examples/errors"
)

const (
	MinAge = 13
	MaxAge = 150
)

//go:embed errors.json
var errorCodes []byte

var (
	// Load error codes.
	_ = errors.MustLoad(errorCodes)

	// All errors.
	ErrNotFound         = errors.New("user.notFound") // For text-only errors without params.
	ErrAlreadyExists    = errors.New("user.alreadyExists")
	ErrInvalidAge       = errors.NewFull("user.invalidAge", InvalidAgeParams{MaxAge})        // For errors with constant params.
	ErrUnderAge         = errors.NewFull("user.underAge", UnderAgeParams{MinAge})            //
	ErrInvalidName      = errors.NewPartial[InvalidNameParams]("user.invalidName")           // For errors with dynamic params.
	ErrValidationErrors = errors.NewPartial[ValidationErrorsParams]("user.validationErrors") //

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
