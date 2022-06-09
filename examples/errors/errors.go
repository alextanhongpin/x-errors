package errors

import (
	"encoding/json"
	stderrors "errors"

	"github.com/alextanhongpin/errors"
)

// Alias to avoid referencing the original package.
type Error = errors.Error

var (
	Is     = stderrors.Is
	As     = stderrors.As
	Unwrap = stderrors.Unwrap
)

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

func MustLoad(errorCodes []byte) bool {
	return bundle.MustLoad(errorCodes)
}

func New(code string) *Error {
	return bundle.Code(errors.Code(code))
}

func NewPartial[T any](code string) *errors.Partial[T] {
	return errors.NewPartial[T](New(code))
}

func NewFull[T any](code string, params T) *errors.Error {
	return NewPartial[T](code).WithParams(params)
}
