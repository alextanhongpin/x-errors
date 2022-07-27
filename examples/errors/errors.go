package errors

import (
	"encoding/json"
	stderrors "errors"

	"github.com/alextanhongpin/errors"
)

// Alias to avoid referencing the original package.
type (
	Error = errors.Error
	Tags  = errors.Tags
	Tag   = errors.Tag
)

func T(key, value string) Tag {
	return Tag{Key: key, Value: value}
}

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

func New(code string, tags ...Tag) *Error {
	return bundle.Code(errors.Code(code), tags...)
}

func NewPartial[T any](code string, tags ...Tag) *errors.Partial[T] {
	return errors.NewPartial[T](New(code, tags...))
}

func NewFull[T any](code string, params T, tags ...Tag) *errors.Error {
	return NewPartial[T](code, tags...).WithParams(params)
}
