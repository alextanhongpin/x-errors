package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Errors
var (
	ErrInvalidKind   = errors.New("errors: kind is invalid")
	ErrCodeNotFound  = errors.New("errors: error code not found")
	ErrDuplicateCode = errors.New("errors: duplicate error code")
)

type Code string
type Kind string

type Bundle struct {
	errorByCode  map[Code]*Error
	allowedKinds map[Kind]bool
	unmarshalFn  func(raw []byte, v any) error
}

type Options struct {
	AllowedKinds []Kind
	UnmarshalFn  func([]byte, any) error
}

func NewBundle(opt *Options) *Bundle {
	if opt == nil {
		opt = &Options{
			UnmarshalFn:  json.Unmarshal,
			AllowedKinds: make([]Kind, 0),
		}
	}

	if opt.UnmarshalFn == nil {
		opt.UnmarshalFn = json.Unmarshal
	}

	allowedKinds := make(map[Kind]bool)
	for _, kind := range opt.AllowedKinds {
		allowedKinds[kind] = true
	}

	return &Bundle{
		allowedKinds: allowedKinds,
		errorByCode:  make(map[Code]*Error),
		unmarshalFn:  opt.UnmarshalFn,
	}
}

func (b *Bundle) Load(errorBytes []byte) error {
	var data map[Kind]map[Code]string
	if err := b.unmarshalFn(errorBytes, &data); err != nil {
		return err
	}

	for kind, messageByCode := range data {
		if err := b.validateKindExists(kind); err != nil {
			return err
		}

		for code, message := range messageByCode {
			if err := b.validateCodeUnique(code); err != nil {
				return err
			}

			b.errorByCode[code] = &Error{
				Code:    string(code),
				Kind:    string(kind),
				Message: message,
				Params:  nil,
				Tags:    nil,
			}
		}
	}

	return nil
}

func (b *Bundle) MustLoad(errorBytes []byte) bool {
	if err := b.Load(errorBytes); err != nil {
		panic(err)
	}

	return true
}

func (b *Bundle) Add(kind Kind, code Code, message string, tags ...Tag) *Error {
	if err := b.validateCodeUnique(code); err != nil {
		panic(fmt.Errorf("%w: %s", ErrDuplicateCode, code))
	}

	err := &Error{
		Code:    string(code),
		Kind:    string(kind),
		Message: message,
		Params:  nil,
		Tags:    tags,
	}

	b.errorByCode[code] = err

	return err
}

func (b *Bundle) Get(code Code) *Error {
	if err := b.validateCodeExists(code); err != nil {
		panic(err)
	}

	return b.errorByCode[code]
}

func (b *Bundle) validateCodeExists(code Code) error {
	if _, found := b.errorByCode[code]; !found {
		return fmt.Errorf("%w: %s", ErrCodeNotFound, code)
	}

	return nil
}

func (b *Bundle) validateCodeUnique(code Code) error {
	if _, found := b.errorByCode[code]; found {
		return fmt.Errorf("%w: %s", ErrDuplicateCode, code)
	}

	return nil
}

func (b *Bundle) validateKindExists(kind Kind) error {
	if b.allowedKinds[kind] {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidKind, kind)
}
