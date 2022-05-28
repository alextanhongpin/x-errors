package errors

type Partial[T any] struct {
	err *Error
}

func NewPartial[T any](err *Error) *Partial[T] {
	e := err.clone()
	e.partial = true

	return &Partial[T]{
		err: e,
	}
}

func (e *Partial[T]) SetParams(params T) *Error {
	err := e.err.clone()

	for lang, msg := range e.err.translations {
		tmsg, terr := makeTemplate(msg, params)
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
