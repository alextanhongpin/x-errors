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
		*errors.Error
		Params *bool `json:"params,omitempty"`
	}

	return json.Marshal(response{Error: e.Error})
}

func main() {
	debug(user.ErrInvalidName.WithParams(user.InvalidNameParams{
		Name: "john appleseed",
	}))
}

func debug(err error) {
	fmt.Println("is ErrInvalidName?", errors.Is(err, user.ErrInvalidName.WithParams(user.InvalidNameParams{})))

	var custom *errors.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)

	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrInvalidName.Unwrap())

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
}
