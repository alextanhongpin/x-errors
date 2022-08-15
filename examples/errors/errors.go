package errors

import (
	"encoding/json"
	stderrors "errors"

	"github.com/alextanhongpin/errors"
)

// Alias to avoid referencing the original package.
type (
	Error = errors.Error
	Tag   = errors.Tag
)

var (
	Is     = stderrors.Is
	As     = stderrors.As
	Unwrap = stderrors.Unwrap
)

const (
	Unknown             errors.Kind = "unknown"
	Internal            errors.Kind = "internal"
	BadInput            errors.Kind = "bad_input"
	NotFound            errors.Kind = "not_found"
	AlreadyExists       errors.Kind = "already_exists"
	FailedPreconditions errors.Kind = "failed_preconditions"
	Unauthorized        errors.Kind = "unauthorized"
	Forbidden           errors.Kind = "forbidden"
)

var bundle = errors.NewBundle(&errors.Options{
	AllowedKinds: []errors.Kind{
		Unknown,
		Internal,
		BadInput,
		NotFound,
		AlreadyExists,
		FailedPreconditions,
		Unauthorized,
		Forbidden,
	},
	UnmarshalFn: json.Unmarshal,
})

func MustLoad(errorCodes []byte) bool {
	return bundle.MustLoad(errorCodes)
}

func New(code string) *Error {
	return bundle.Get(errors.Code(code))
}

func NewPartial[T any](code string) *errors.PartialError[T] {
	return errors.ToPartial[T](New(code))
}

func NewFull[T any](code string, params T) *errors.Error {
	return NewPartial[T](code).WithParams(params)
}
