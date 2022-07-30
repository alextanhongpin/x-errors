package errors

import (
	"bytes"
	"errors"
	"html/template"
)

type Tags map[string]string

func (t Tags) IsZero() bool {
	return t != nil
}

func (t Tags) Has(key string) bool {
	if t.IsZero() {
		return false
	}

	_, ok := t[key]

	return ok
}

func (t Tags) Clone() Tags {
	if t.IsZero() {
		return nil
	}

	res := make(Tags)

	for k, v := range t {
		res[k] = v
	}

	return res
}

// Error represents an error.
type Error struct {
	Kind    string `json:"kind"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Params  any    `json:"params"`
	Tags    Tags   `json:"tags,omitempty"`
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

func (e *Error) WithParams(params any) *Error {
	err := e.Clone()
	err.Params = params

	return err
}

func (e *Error) WithTag(tags ...Tag) *Error {
	err := e.Clone()
	if err.Tags == nil {
		err.Tags = make(Tags)
	}

	for _, tag := range tags {
		err.Tags[tag[0]] = tag[1]
	}

	return err
}

func (e *Error) WithTags(tags Tags) *Error {
	err := e.Clone()
	if err.Tags == nil {
		err.Tags = make(Tags)
	}

	for k, v := range tags {
		err.Tags[k] = v
	}

	return err
}

func (e *Error) Clone() *Error {
	clone := *e
	clone.Tags = e.Tags.Clone()

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
	return p.err.WithParams(t)
}

func (p *Partial[T]) WithTag(tags ...Tag) *Error {
	return p.err.WithTag(tags...)
}

func (p *Partial[T]) WithTags(tags Tags) *Error {
	return p.err.WithTags(tags)
}

func exec[T any](msg string, data T) string {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg
	}

	return buf.String()
}
