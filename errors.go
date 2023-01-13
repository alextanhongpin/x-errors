package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
)

// Error represents an error.
type Error struct {
	Kind    string
	Code    string
	Message string
	Params  any
	Tags    []string
	Cause   error
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
func (e *Error) Is(target error) bool {
	var err *Error
	if !errors.As(target, &err) {
		return false
	}

	return err.Code == e.Code
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Wrap(cause error) *Error {
	if cause == nil {
		return e
	}

	err := e.Copy()
	err.Cause = cause

	return err
}

func (e *Error) WithParams(params any) *Error {
	err := e.Copy()
	err.Params = params

	return err
}

func (e *Error) WithTag(tags ...string) *Error {
	err := e.Copy()
	err.Tags = append(err.Tags, tags...)

	return err
}

func (e *Error) Copy() *Error {
	err := *e
	tags := make([]string, len(err.Tags))
	copy(tags, e.Tags)
	err.Tags = tags

	return &err
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind    string   `json:"kind"`
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Params  any      `json:"params,omitempty"`
		Tags    []string `json:"tags,omitempty"`
	}{
		Kind:    e.Kind,
		Code:    e.Code,
		Message: e.Error(), // The message returned should be formatted
		Params:  e.Params,
		Tags:    e.Tags,
	})
}

type PartialError[T any] struct {
	err *Error
}

func ToPartial[T any](err *Error) *PartialError[T] {
	return &PartialError[T]{
		err: err,
	}
}

func (p *PartialError[T]) Unwrap() *Error {
	return p.err
}

func (p *PartialError[T]) WithParams(t T) *Error {
	return p.err.WithParams(t)
}

func (p *PartialError[T]) WithTag(tags ...string) *PartialError[T] {
	return ToPartial[T](p.err.WithTag(tags...))
}

func exec[T any](msg string, data T) string {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg
	}

	return buf.String()
}
