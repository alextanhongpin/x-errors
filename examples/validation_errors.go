package main

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/examples/domain/user"
	"github.com/alextanhongpin/errors/examples/errors"
)

func main() {
	debug(user.ValidationErrors([]error{user.InvalidNameError("john appleseed"), user.InvalidAgeError()}))
}

func debug(err error) {
	fmt.Println("is ErrValidationErrors?", errors.Is(err, user.ErrValidationErrors))

	var custom *errors.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)
	fmt.Println("error?", custom.Error())

	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrValidationErrors)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))
}
