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
}

// Error fulfills the error interface.
func (e Error) Error() string {
	if e.Params != nil {
		return exec(e.Message, e.Params)
	}

	return e.Message
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

type Partial[T any] struct {
	err *Error
}

func NewPartial[T any](err *Error) *Partial[T] {
	return &Partial[T]{
		err: err,
	}
}

func (p *Partial[T]) Unwrap() *Error {
	return p.err
}

func (p *Partial[T]) WithParams(t T) *Error {
	err := p.err.clone()
	err.Params = t

	return err
}

func exec[T any](msg string, data T) string {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg
	}

	return buf.String()
}
