package errors_test

import (
	"testing"

	"github.com/alextanhongpin/errors"
)

func TestTags(t *testing.T) {
	tags := make(errors.Tags)
	tags["op"] = "findUser"

	err := errors.Error{
		Code:    "CodeNotFound",
		Kind:    "NotFound",
		Message: "Not found",
		Params:  nil,
		Tags:    tags,
	}

	clone := err.WithTag(errors.Tag{"op", "edited"})
	tags["op"] = "FindUser"

	if exp, got := "FindUser", err.Tags["op"]; exp != got {
		t.Fatalf("expected parent to be %s, got %s", exp, got)
	}

	if exp, got := "edited", clone.Tags["op"]; exp != got {
		t.Fatalf("expected parent to be %s, got %s", exp, got)
	}
}
