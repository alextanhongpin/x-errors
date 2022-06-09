package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/examples/app"
	"github.com/alextanhongpin/errors/examples/domain/user"
)

func main() {
	debug(user.ErrValidationErrors(
		user.ValidationErrorsParams{
			Count: 2,

			// NOTE: Not locale sensitive ...
			PluralError: "errors",
			Errors: []error{
				user.ErrInvalidName(user.InvalidNameParams{
					Name: "john appleseed",
				}),
				user.ErrInvalidAge,
			},
		},
	))
}

func debug(err error) {
	fmt.Println("is ErrValidationErrors?", errors.Is(err, user.ErrValidationErrors(user.ValidationErrorsParams{})))

	var custom *app.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)
	fmt.Println("error?", custom.Error())

	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrValidationErrors(user.ValidationErrorsParams{}))

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))
}
