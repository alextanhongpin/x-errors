package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
)

type Tag string

type Tags []Tag

func (t Tags) IsZero() bool {
	return len(t) == 0
}

func (tags Tags) Has(tag string) bool {
	for _, t := range tags {
		if t == Tag(tag) {
			return true
		}
	}

	return false
}

func (t Tags) Copy() Tags {
	res := make(Tags, len(t))
	copy(res, t)

	return res
}

// Error represents an error.
type Error struct {
	Kind    string `json:"kind"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Params  any    `json:"params,omitempty"`
	Tags    Tags   `json:"tags,omitempty"`
	Cause   error  `json:"cause,omitempty"`
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

func (e *Error) WithTag(tags ...Tag) *Error {
	err := e.Copy()
	err.Tags = append(err.Tags, tags...)

	return err
}

func (e *Error) Copy() *Error {
	err := *e
	err.Tags = e.Tags.Copy()

	return &err
}

func (e *Error) MarshalJSON() ([]byte, error) {
	type errorResponse struct {
		Kind    string `json:"kind"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Params  any    `json:"params,omitempty"`
		Tags    Tags   `json:"tags,omitempty"`
	}

	res := errorResponse{
		Kind:    e.Kind,
		Code:    e.Code,
		Message: e.Error(), // The message returned should be formatted
		Params:  e.Params,
		Tags:    e.Tags,
	}

	return json.Marshal(res)
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

func (p *PartialError[T]) WithTag(tags ...Tag) *PartialError[T] {
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
