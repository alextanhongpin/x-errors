package errors

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

var (
	ErrTranslationUndefined = errors.New("translation is not defined")
	ErrInvalidKind          = errors.New("kind is invalid")
	ErrCodeNotFound         = errors.New("error code not found")
	ErrDuplicateCode        = errors.New("duplicate error code")
)

func makeTemplate(msg string, data any) (string, error) {
	t := template.Must(template.New("").Parse(msg))

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return msg, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
