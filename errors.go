package errors

import (
	"bytes"
	"errors"
	"html/template"
)

// Error represents an error.
type Error struct {
	Kind    string `json:"kind"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Params  any    `json:"params"`

	// Stores the original template message.
	template string
}

// Error fulfills the error interface.
func (e Error) Error() string {
	return exec(e.template, e.Params)
}

func (e Error) String() string {
	return e.Error()
}

// Is satisfies the error interface.
func (e Error) Is(target error) bool {
	var err *Error
	if !errors.As(target, &err) {
		return false
	}

	return err.Kind == e.Kind && err.Code == e.Code
}

func (e *Error) clone() *Error {
	clone := *e
	return &clone
}

type Partial[T any] func(t T) *Error

func NewPartial[T any](err *Error) func(t T) *Error {
	return func(t T) *Error {
		e := err.clone()
		e.Params = t

		return e
	}
}

func exec[T any](msg string, data T) string {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg
	}

	return buf.String()
}
