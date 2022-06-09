package main

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/examples/domain/user"
	"github.com/alextanhongpin/errors/examples/errors"
)

func main() {
	debug(user.ErrValidationErrors.WithParams(
		user.ValidationErrorsParams{
			Count: 2,

			// NOTE: Not locale sensitive ...
			PluralError: "errors",
			Errors: []error{
				user.ErrInvalidName.WithParams(user.InvalidNameParams{
					Name: "john appleseed",
				}),
				user.ErrInvalidAge,
			},
		},
	))
}

func debug(err error) {
	fmt.Println("is ErrValidationErrors?", errors.Is(err, user.ErrValidationErrors.WithParams(user.ValidationErrorsParams{})))

	var custom *errors.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)
	fmt.Println("error?", custom.Error())

	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrValidationErrors.Unwrap())

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))
}
