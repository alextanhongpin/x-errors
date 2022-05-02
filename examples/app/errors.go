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
