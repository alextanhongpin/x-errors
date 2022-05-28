package app

import (
	"encoding/json"

	"github.com/alextanhongpin/errors"
	"golang.org/x/text/language"
)

var (
	// Supported app locales.
	EN = language.English
	MS = language.Malay
)

// Alias to avoid referencing the original package.
type Error = errors.Error

var Errors = errors.NewBundle(&errors.Options{
	DefaultLanguage:  EN,
	AllowedLanguages: []language.Tag{MS},
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
	return Errors.MustLoad(errorCodes)
}

func NewError(code string) *Error {
	return Errors.Code(errors.Code(code))
}

func NewPartialError[T any](code string) *errors.Partial[T] {
	return errors.NewPartial[T](NewError(code))
}

func NewFullError[T any](code string, params T) *errors.Error {
	return NewPartialError[T](code).SetParams(params)
}
