package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alextanhongpin/errors/examples/domain/user"
	"github.com/alextanhongpin/errors/examples/errors"
)

var (
	ErrCreateUserAlreadyExists = user.ErrAlreadyExists.WithTag(errors.T("op", "CreateUser"))
	ErrUpdateUserAlreadyExists = user.ErrAlreadyExists.WithTag(errors.T("op", "UpdateUser"))
)

func main() {
	debug(user.ErrInvalidAge)
}

func debug(err error) {
	{
		err := createUser()
		var custom *errors.Error
		if errors.As(err, &custom) {
			fmt.Println("errors.As?", true)
		}
		fmt.Printf("%#v\n", custom)
	}

	{
		err := updateUser()
		var custom *errors.Error
		if errors.As(err, &custom) {
			fmt.Println("errors.As?", true)
		}
		fmt.Printf("%#v\n", custom)
	}

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

	for i, valid := range []bool{
		errors.Is(createUser(), user.ErrAlreadyExists),
		errors.Is(createUser(), ErrCreateUserAlreadyExists),
		errors.Is(updateUser(), user.ErrAlreadyExists),
		errors.Is(updateUser(), ErrUpdateUserAlreadyExists),
	} {
		if !valid {
			log.Fatalf("error at index: %d", i)
		} else {
			fmt.Printf("%d is valid\n", i)
		}
	}

}

func createUser() error {
	return ErrCreateUserAlreadyExists
}

func updateUser() error {
	return ErrUpdateUserAlreadyExists
}
