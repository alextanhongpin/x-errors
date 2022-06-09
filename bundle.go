package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

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
		if !b.allowedKinds[kind] {
			return fmt.Errorf("%w: %s", ErrInvalidKind, kind)
		}

		for code, message := range messageByCode {
			if _, ok := b.errorByCode[code]; ok {
				return fmt.Errorf("%w: %s", ErrDuplicateCode, code)
			}

			b.errorByCode[code] = &Error{
				Code:    string(code),
				Kind:    string(kind),
				Message: message,
				Params:  nil,
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

func (b *Bundle) Code(code Code) *Error {
	err, ok := b.errorByCode[code]
	if !ok {
		panic(fmt.Errorf("%w: %s", ErrCodeNotFound, code))
	}

	return err.clone()
}
