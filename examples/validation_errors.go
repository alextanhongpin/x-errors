package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/examples/app"
	"github.com/alextanhongpin/errors/examples/domain/user"
)

func main() {
	debugError(user.ErrValidationErrors.SetParams(
		[]error{
			user.ErrInvalidName.SetParams(user.InvalidNameParams{
				Name: "john appleseed",
			}),
			user.ErrInvalidAge,
		},
	))
}

func debugError(err error) {
	fmt.Println("is ErrValidationErrors?", errors.Is(err, user.ErrValidationErrors.Self()))

	var custom *app.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("is partial?", custom.IsPartial())
	fmt.Println("message?", custom)

	localized := custom.Localize(app.MS)
	fmt.Println("localized?", localized)
	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrValidationErrors.Self())
	fmt.Println("is parent partial?", user.ErrValidationErrors.Self().IsPartial())

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("original json?", string(b))

	b, err = json.MarshalIndent(localized, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("localized json?", string(b))
}
