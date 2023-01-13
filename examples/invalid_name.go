package main

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/examples/domain/user"
	"github.com/alextanhongpin/errors/examples/errors"
)

type ErrorWithoutParams struct {
	*errors.Error
}

func (e *ErrorWithoutParams) MarshalJSON() ([]byte, error) {
	if e == nil || e.Error == nil {
		return []byte("null"), nil
	}

	type response struct {
		Kind    string `json:"kind"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	return json.Marshal(response{
		Kind:    e.Kind,
		Code:    e.Code,
		Message: e.Error.Error(),
	})
}

func main() {
	debug(user.ErrInvalidName.WithParams(map[string]string{
		"Name": "alice",
	}))
}

func debug(err error) {
	{
		err := createUser()
		var custom *errors.Error
		if errors.As(err, &custom) {
			fmt.Println("errors.As?", true)
		}
		fmt.Printf("%#v\n", custom)

		fmt.Println(errors.Is(err, user.ErrInvalidName))
	}

	fmt.Println("is ErrInvalidName?", errors.Is(err, user.ErrInvalidName))

	var custom *errors.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)

	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrInvalidName)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))

	b, err = json.MarshalIndent(&ErrorWithoutParams{Error: custom}, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("stripped params?", string(b))
	fmt.Println()
}

func createUser() error {
	return user.InvalidNameError("john", "op:create_user")
}
