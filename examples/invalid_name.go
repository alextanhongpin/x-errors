package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alextanhongpin/errors/examples/app"
	"github.com/alextanhongpin/errors/examples/domain/user"
)

type ErrorWithoutParams struct {
	*app.Error
}

func (e *ErrorWithoutParams) MarshalJSON() ([]byte, error) {
	if e == nil || e.Error == nil {
		return []byte("null"), nil
	}

	type response struct {
		*app.Error
		Params *bool `json:"params,omitempty"`
	}

	return json.Marshal(response{Error: e.Error})
}

func main() {
	debugError(user.ErrInvalidName.SetParams(user.InvalidNameParams{
		Name: "john appleseed",
	}))
}

func debugError(err error) {
	fmt.Println("is ErrInvalidName?", errors.Is(err, user.ErrInvalidName.Self()))

	var custom *app.Error
	if errors.As(err, &custom) {
		fmt.Println("errors.As?", true)
	}

	fmt.Println("is partial?", custom.IsPartial())
	fmt.Println("message?", custom)

	localized := custom.Localize(app.MS)
	fmt.Println("localized?", localized)
	fmt.Println("is original modified?", err)
	fmt.Println("is parent modified?", user.ErrInvalidName.Self())
	fmt.Println("is parent partial?", user.ErrInvalidName.Self().IsPartial())

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

	b, err = json.MarshalIndent(&ErrorWithoutParams{Error: custom}, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("stripped params?", string(b))
}
