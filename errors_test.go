package errors_test

import (
	"testing"

	"github.com/alextanhongpin/errors"
)

func TestTagsImmutable(t *testing.T) {

	parent := errors.Error{
		Code:    "CodeNotFound",
		Kind:    "NotFound",
		Message: "Not found",
		Params:  nil,
		Tags:    nil,
	}

	child := parent.WithTag("repo.FindUser")

	if exp, got := int(0), len(parent.Tags); exp != got {
		t.Fatalf("expected parent tags length to be %d, got %d", exp, got)
	}
	t.Logf("parent: %v", parent.Tags)

	if exp, got := int(1), len(child.Tags); exp != got {
		t.Fatalf("expected parent to be %d, got %d", exp, got)
	}
	t.Logf("child: %v", child.Tags)
}
