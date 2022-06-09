package app

import (
	"encoding/json"

	"github.com/alextanhongpin/errors"
)

// Alias to avoid referencing the original package.
type Error = errors.Error

var bundle = errors.NewBundle(&errors.Options{
	AllowedKinds: []errors.Kind{
		"unknown",
		"internal",
		"bad_input",
		"not_found",
		"already_exists",
		"failed_preconditions",
		"unauthorized",
		"forbidden",
	},
	UnmarshalFn: json.Unmarshal,
})

func MustLoadError(errorCodes []byte) bool {
	return bundle.MustLoad(errorCodes)
}

func NewError(code string) *Error {
	return bundle.Code(errors.Code(code))
}

func NewPartialError[T any](code string) errors.Partial[T] {
	return errors.NewPartial[T](NewError(code))
}

func NewFullError[T any](code string, params T) *errors.Error {
	return NewPartialError[T](code)(params)
}
