package main

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/examples/domain/user"
	"github.com/alextanhongpin/errors/examples/errors"
)

func main() {
	debug(user.ErrInvalidAge)
}

func debug(err error) {
	fmt.Println("is ErrInvalidAge?", errors.Is(err, user.ErrInvalidAge))

	var custom *errors.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("message?", custom)
	fmt.Println("error?", custom.Error())
	fmt.Println("is original modified?", user.ErrInvalidAge)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))
}
