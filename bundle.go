package errors

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
)

// Errors
var (
	ErrInvalidKind   = errors.New("errors: kind is invalid")
	ErrCodeNotFound  = errors.New("errors: error code not found")
	ErrDuplicateCode = errors.New("errors: duplicate error code")
)

var (
	// Default bundle.
	bundle         = NewBundle(nil)
	Get            = bundle.Get
	Add            = bundle.Add
	AddKinds       = bundle.AddKinds
	MustAddKinds   = bundle.MustAddKinds
	SetUnmarshalFn = bundle.SetUnmarshalFn
	Len            = bundle.Len
	Load           = bundle.Load
	Iter           = bundle.Iter
	MustLoad       = bundle.MustLoad
	LoadFS         = bundle.LoadFS
	MustLoadFS     = bundle.MustLoadFS
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

func (b *Bundle) AddKinds(kinds ...Kind) error {
	for _, k := range kinds {
		b.allowedKinds[k] = true
	}

	return b.Iter(func(code Code, err *Error) error {
		if e := b.validateKindExists(Kind(err.Kind)); e != nil {
			return e
		}

		return nil
	})
}

func (b *Bundle) MustAddKinds(kinds ...Kind) bool {
	if err := b.AddKinds(kinds...); err != nil {
		panic(err)
	}

	return true
}

// SetUnmarshalFn sets the unmarshal function and returns true to allow
// variable initialization.
func (b *Bundle) SetUnmarshalFn(unmarshalFn func(raw []byte, v any) error) bool {
	b.unmarshalFn = unmarshalFn
	return true
}

func (b *Bundle) Len() int {
	return len(b.errorByCode)
}

func (b *Bundle) Iter(fn func(code Code, err *Error) error) error {
	for code, err := range b.errorByCode {
		if e := fn(code, err); e != nil {
			return e
		}
	}

	return nil
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

func (b *Bundle) LoadFS(fs embed.FS) error {
	dirs := []string{"."}
	for len(dirs) > 0 {
		var dir string
		dir, dirs = dirs[0], dirs[1:]

		entries, err := fs.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				return err
			}

			dir := filepath.Join(dir, info.Name())
			if info.IsDir() {
				dirs = append(dirs, dir)
				continue
			}

			by, err := fs.ReadFile(dir)
			if err != nil {
				return err
			}
			if err := b.Load(by); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Bundle) MustLoadFS(fs embed.FS) bool {
	if err := b.LoadFS(fs); err != nil {
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
	// NOTE: If we do not set allowed kinds, by default all
	// is allowed.
	if len(b.allowedKinds) == 0 || b.allowedKinds[kind] {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidKind, kind)
}
