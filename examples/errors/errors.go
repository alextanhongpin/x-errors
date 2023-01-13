package errors

import (
	stderrors "errors"

	"github.com/alextanhongpin/errors"
)

// Alias to avoid referencing the original package.
type (
	Error = errors.Error
)

var (
	// Export std errors functionality.
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

var (
	// Export custom errors functionality.
	_ = errors.MustAddKinds(
		Unknown,
		Internal,
		BadInput,
		NotFound,
		AlreadyExists,
		FailedPreconditions,
		Unauthorized,
		Forbidden,
	)
	Get      = errors.Get
	MustLoad = errors.MustLoad
)

func ToPartial[T any](err *Error) *errors.PartialError[T] {
	return errors.ToPartial[T](err)
}
