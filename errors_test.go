package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/alextanhongpin/errors"
)

func TestTagsImmutable(t *testing.T) {
	parent := errors.Error{
		Code:    "user.NotFound",
		Kind:    "not_found",
		Message: "The user you requested does not exists or may have been deleted.",
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

func TestWrapError(t *testing.T) {
	err := errors.Error{
		Code:    "user.NotFound",
		Kind:    "not_found",
		Message: "The user you requested does not exists or may have been deleted.",
		Params:  nil,
		Tags:    nil,
	}

	var cause error
	cause = err.WithTag("repo.not_found")

	result := errors.Wrap(err.WithTag("login_usecase.not_found"), cause)

	var target *errors.Error
	if !stderrors.As(result, &target) {
		t.Fatal("type assertion to *errors.Error failed")
	}

	if exp, got := int(2), len(target.Tags); exp != got {
		t.Fatalf("expected length to be %d, got %d", exp, got)
	}

	t.Logf("error: %v", result)
	t.Logf("target: %#v", target)
}
