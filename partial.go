package errors

type Partial[T any] struct {
	err          *Error
	templateFunc TemplateFunc[T]
}

type PartialOption[T any] func(*Partial[T])

func WithTemplateFunc[T any](fn TemplateFunc[T]) PartialOption[T] {
	return func(p *Partial[T]) {
		p.templateFunc = fn
	}
}

func NewPartial[T any](err *Error, options ...PartialOption[T]) *Partial[T] {
	e := err.clone()
	e.partial = true

	p := &Partial[T]{
		err:          e,
		templateFunc: templateWithLanguageTag[T],
	}

	for _, opt := range options {
		opt(p)
	}

	return p
}

func (e *Partial[T]) SetParams(params T) *Error {
	err := e.err.clone()

	for lang, msg := range e.err.translations {
		tmsg, terr := e.templateFunc(lang, msg, params)
		if terr != nil {
			err.translations[lang] = msg
		} else {
			err.translations[lang] = tmsg
		}
	}

	err.Message = err.translations[err.lang]
	err.Params = params
	err.partial = false

	return err
}

// Unwrap returns the base error. Useful when checking errors.Is without setting
// the params.
func (e *Partial[T]) Unwrap() *Error {
	return e.err
}
