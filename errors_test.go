package errors_test

import (
	"embed"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/alextanhongpin/errors"
)

// NOTE: Golang's embed does not support recursive globbing like **/*.json
// Every path must be defined manually.
//
//go:embed testdata/*.json testdata/*/*.json
var errorFiles embed.FS
var ErrTest = stderrors.New("test")

func TestTagsImmutable(t *testing.T) {
	parent := &errors.Error{
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
	parentErr := &errors.Error{
		Code:    "user.NotFound",
		Kind:    "not_found",
		Message: "The user you requested does not exists or may have been deleted.",
		Params:  nil,
		Tags:    nil,
	}

	childErr := fmt.Errorf("not found: %w", parentErr.WithTag("login_usecase"))

	if !stderrors.Is(childErr, parentErr) {
		t.Fatal("child not equal to parent error")
	}

	var target *errors.Error
	if !stderrors.As(childErr, &target) {
		t.Fatal("type assertion to *errors.Error failed")
	}

	if exp, got := int(1), len(target.Tags); exp != got {
		t.Fatalf("expected length to be %d, got %d", exp, got)
	}

	if exp, got := true, stderrors.Is(target.Wrap(ErrTest), ErrTest); exp != got {
		t.Fatal("expected error to be ErrTest")
	}

	if exp, got := true, stderrors.Is(target.Wrap(ErrTest), parentErr); exp != got {
		t.Fatal("expected error to be ErrTest")
	}

	t.Logf("error: %v", childErr)
	t.Logf("target: %#v", target)
}

func TestEmbed(t *testing.T) {
	bundle := errors.NewBundle()
	if !bundle.MustLoadFS(errorFiles, json.Unmarshal) {
		t.Fatal("failed to load error bundle")
	}

	if exp, got := 6, bundle.Len(); exp != got {
		t.Fatalf("expected %d errors to be loaded, got %d", exp, got)
	}

	if err := bundle.Get(errors.Code("user.notFound")); err == nil {
		t.Fatalf("expected error user.notFound to be loaded, got nil")
	}

	if err := bundle.Get(errors.Code("user.invalidName")); err == nil {
		t.Fatalf("expected error user.invalidName to be loaded, got nil")
	}
}

func TestDefaultBundle(t *testing.T) {
	errors.MustLoadFS(errorFiles, json.Unmarshal)
	if err := errors.AddKinds("not_found", "bad_input"); err != nil {
		t.Fatalf("expected set kinds to return nil, got %v", err)
	}

	if exp, got := 6, errors.Len(); exp != got {
		t.Fatalf("expected %d errors to be loaded, got %d", exp, got)
	}

	if err := errors.Get(errors.Code("user.notFound")); err == nil {
		t.Fatalf("expected error user.notFound to be loaded, got nil")
	}

	if err := errors.Get(errors.Code("user.invalidName")); err == nil {
		t.Fatalf("expected error user.invalidName to be loaded, got nil")
	}

	t.Run("marshal", func(t *testing.T) {
		type userInvalidNameError struct {
			Name string
		}

		errRes := errors.ToPartial[userInvalidNameError](errors.Get("user.invalidName")).WithParams(userInvalidNameError{
			Name: "john",
		})
		b, err := json.Marshal(errRes)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", b)
	})
}
