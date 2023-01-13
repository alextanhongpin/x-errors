package user

import (
	_ "embed"
	"encoding/json"

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
	_ = errors.MustLoad(errorCodes, json.Unmarshal)

	// All errors.
	ErrNotFound         = errors.Get("user.notFound") // For text-only errors without params.
	ErrAlreadyExists    = errors.Get("user.alreadyExists")
	ErrInvalidAge       = errors.Get("user.invalidAge")       // For errors with constant params.
	ErrUnderAge         = errors.Get("user.underAge")         //
	ErrInvalidName      = errors.Get("user.invalidName")      // For errors with dynamic params.
	ErrValidationErrors = errors.Get("user.validationErrors") //

)

func InvalidAgeError() error {
	type params struct {
		MaxAge int64 `json:"maxAge"`
	}
	return errors.ToPartial[params](ErrInvalidAge).WithParams(params{
		MaxAge: MaxAge,
	})
}

func UnderAgeError() error {
	type params struct {
		MinAge int64 `json:"minAge"`
	}
	return errors.ToPartial[params](ErrUnderAge).WithParams(params{
		MinAge: MinAge,
	})
}

func InvalidNameError(name string, tags ...string) error {
	type params struct {
		Name string `json:"name"`
	}
	return errors.ToPartial[params](ErrInvalidName).
		WithParams(params{
			Name: name,
		}).
		WithTag(tags...)
}

func ValidationErrors(errs []error) error {
	type params struct {
		Count       int     `json:"count"`
		PluralError string  `json:"-"`
		Errors      []error `json:"errors"`
	}
	count := len(errs)
	var pluralError string
	if count == 1 {
		pluralError = "error"
	} else {
		pluralError = "errors"
	}
	return errors.ToPartial[params](ErrValidationErrors).WithParams(params{
		Count:       count,
		PluralError: pluralError,
		Errors:      errs,
	})

}
