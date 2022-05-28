package errors

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"
)

var ErrPartial = errors.New("rendering incomplete error")

// Error represents an error.
type Error struct {
	Kind    string `json:"kind"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Params  any    `json:"params"`

	// Error inherits the language and translations from the bundle.
	lang         language.Tag
	translations map[language.Tag]string
	partial      bool
}

func (e *Error) Localize(lang language.Tag) *Error {
	if e.partial {
		panic(fmt.Errorf("%w: %q.%q", ErrPartial, e.Kind, e.Code))
	}

	err := e.clone()

	msg, ok := err.translations[lang]
	if !ok {
		panic(fmt.Errorf("%w: %q.%q", ErrTranslationUndefined, err.Code, lang))
	}

	err.lang = lang
	err.Message = msg

	return err
}

// Error fulfills the error interface.
func (e Error) Error() string {
	return e.Message
}

// Is satisfies the error interface.
func (e Error) Is(target error) bool {
	var err *Error
	if !errors.As(target, &err) {
		return false
	}

	return err.Kind == e.Kind && err.Code == e.Code
}

func (e *Error) IsPartial() bool {
	return e.partial
}

func (e *Error) clone() *Error {
	clone := *e
	clone.translations = make(map[language.Tag]string)

	for k, v := range e.translations {
		clone.translations[k] = v
	}

	return &clone
}
