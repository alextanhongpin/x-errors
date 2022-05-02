package errors

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"
)

type Code string
type Kind string

type Bundle struct {
	errorByCode      map[Code]*Error
	allowedKinds     map[Kind]bool
	allowedLanguages map[language.Tag]bool
	defaultLang      language.Tag
	unmarshalFn      func(raw []byte, v any) error
}

type Options struct {
	DefaultLanguage  language.Tag
	AllowedLanguages []language.Tag
	AllowedKinds     []Kind
	UnmarshalFn      func([]byte, any) error
}

func NewBundle(opt *Options) *Bundle {
	if opt == nil {
		opt = &Options{
			DefaultLanguage:  language.English,
			UnmarshalFn:      json.Unmarshal,
			AllowedKinds:     make([]Kind, 0),
			AllowedLanguages: make([]language.Tag, 0),
		}
	}

	if opt.DefaultLanguage == language.Und {
		opt.DefaultLanguage = language.English
	}

	if opt.UnmarshalFn == nil {
		opt.UnmarshalFn = json.Unmarshal
	}

	allowedLanguages := make(map[language.Tag]bool)
	allowedLanguages[opt.DefaultLanguage] = true

	for _, lang := range opt.AllowedLanguages {
		allowedLanguages[lang] = true
	}

	allowedKinds := make(map[Kind]bool)
	for _, kind := range opt.AllowedKinds {
		allowedKinds[kind] = true
	}

	return &Bundle{
		defaultLang:      opt.DefaultLanguage,
		allowedLanguages: allowedLanguages,
		allowedKinds:     allowedKinds,
		errorByCode:      make(map[Code]*Error),
		unmarshalFn:      opt.UnmarshalFn,
	}
}

func (b *Bundle) Load(errorBytes []byte) error {
	var data map[Kind]map[Code]map[string]string
	if err := b.unmarshalFn(errorBytes, &data); err != nil {
		return err
	}

	for kind, translationsByCode := range data {
		if !b.allowedKinds[kind] {
			return fmt.Errorf("%w: %s", ErrInvalidKind, kind)
		}

		for code, messageByLanguage := range translationsByCode {
			if _, ok := b.errorByCode[code]; ok {
				return fmt.Errorf("%w: %s", ErrDuplicateCode, code)
			}

			translations := make(map[language.Tag]string)
			for lang, message := range messageByLanguage {
				translations[language.MustParse(lang)] = message
			}

			for lang := range b.allowedLanguages {
				if _, ok := translations[lang]; !ok {
					return fmt.Errorf("%w: %q.%q", ErrTranslationUndefined, code, lang)
				}
			}

			b.errorByCode[code] = &Error{
				Code:         string(code),
				Kind:         string(kind),
				Message:      translations[b.defaultLang],
				Params:       nil,
				lang:         b.defaultLang,
				translations: translations,
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

	return err.Clone()
}
