package main

import (
	_ "embed"
	"encoding/json"
	stderrors "errors"
	"fmt"

	"github.com/alextanhongpin/errors"
	"github.com/alextanhongpin/errors/examples/app"
	"github.com/alextanhongpin/errors/examples/domain/user"
)

func main() {
	debugErrValidationErrors(user.ErrValidationErrors.SetParams([]user.ValidationFieldError{
		{Field: "name", Reason: "required"},
		{Field: "age", Reason: "required"},
	}))
	debugErrInvalidAge(user.ErrInvalidAge)
	debugErrInvalidName(user.ErrInvalidName.SetParams(user.InvalidNameParams{
		Name: "J@hn",
	}))
}

func debugErrValidationErrors(err error) {
	fmt.Println("ErrValidationErrors:", err)
	fmt.Println(stderrors.Is(err, user.ErrValidationErrors.Self()))

	{
		b, merr := json.Marshal(err)
		if merr != nil {
			panic(merr)
		}
		fmt.Println("MarshalErrValidationErrors: err", string(b))
	}
}

func debugErrInvalidAge(err error) {
	fmt.Println("ErrInvalidAge:", err)
	fmt.Println(stderrors.Is(err, user.ErrInvalidAge))

	var userError *errors.Error
	if stderrors.As(err, &userError) {
		fmt.Println("Cast err back to UserInvalidAge success")
		fmt.Println(userError.SetLanguage(app.MS))
	}

	{
		b, merr := json.Marshal(err)
		if merr != nil {
			panic(merr)
		}
		fmt.Println("MarshalErrInvalidAge: err", string(b))
	}

	{
		b, merr := json.Marshal(userError)
		if merr != nil {
			panic(merr)
		}
		fmt.Println("MarshalErrInvalidAge: userError EN", string(b))
	}

	{
		b, merr := json.Marshal(userError.SetLanguage(app.MS))
		if merr != nil {
			panic(merr)
		}
		fmt.Println("MarshalErrInvalidAge: userError MS", string(b))
	}
}

func debugErrInvalidName(err error) {
	fmt.Println("ErrInvalidName:", err)
	fmt.Println("IsErrInvalidName?", stderrors.Is(err, user.ErrInvalidName.Self()))

	{
		b, err := json.Marshal(err)
		if err != nil {
			panic(err)
		}
		fmt.Println("MarshalErrInvalidName: error", string(b))
	}

	var userError *errors.Error
	if stderrors.As(err, &userError) {
		fmt.Println("MarshalErrInvalidNameSuccess")
		{
			b, err := json.Marshal(userError)
			if err != nil {
				panic(err)
			}
			fmt.Println("MarshalErrInvalidName: userError EN", string(b))
		}
		{
			b, err := json.Marshal(userError.SetLanguage(app.MS))
			if err != nil {
				panic(err)
			}
			fmt.Println("MarshalErrInvalidName: userError MS", string(b))
		}
	} else {
		panic("MarshalErrInvalidNameFailed")
	}
}
