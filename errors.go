package errors

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"golang.org/x/text/language"
)

var (
	ErrTranslationUndefined = errors.New("translation is not defined")
	ErrInvalidKind          = errors.New("kind is invalid")
	ErrCodeNotFound         = errors.New("error code not found")
	ErrDuplicateCode        = errors.New("duplicate error code")
	ErrPartial              = errors.New("rendering incomplete error")
)

type TemplateFunc[T any] func(lang language.Tag, msg string, data T) (string, error)

func Template[T any](msg string, data T) (string, error) {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func templateWithLanguageTag[T any](lang language.Tag, msg string, data T) (string, error) {
	return Template(msg, data)
}
