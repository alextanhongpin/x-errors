package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/examples/app"
	"github.com/alextanhongpin/errors/examples/domain/user"
)

func main() {
	debug(user.ErrInvalidAge)
}

func debug(err error) {
	fmt.Println("is ErrInvalidAge?", errors.Is(err, user.ErrInvalidAge))

	var custom *app.Error
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
