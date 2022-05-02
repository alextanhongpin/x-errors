package errors

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"
)

// Error represents an error.
type Error struct {
	Kind         string `json:"kind"`
	Code         string `json:"code"`
	Message      string `json:"message"`
	Params       any    `json:"params"`
	lang         language.Tag
	translations map[language.Tag]string
}

func (e *Error) SetLanguage(lang language.Tag) *Error {
	err := e.Clone()

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
	if errors.As(target, &err) {
		return err.Kind == e.Kind && err.Code == e.Code
	}

	return false
}

func (e *Error) Clone() *Error {
	cerr := *e
	cerr.translations = make(map[language.Tag]string)

	for k, v := range e.translations {
		cerr.translations[k] = v
	}

	return &cerr
}
