package errorx

import (
	"fmt"
	"testing"
)

type mockErr struct{}

func (*mockErr) Error() string {
	return "mock error"
}

func TestWarp(t *testing.T) {
	var err error = &mockErr{}
	err2 := fmt.Errorf("wrap %w", err)
	if err != Unwrap(err2) {
		t.Errorf("got %v want: %v", err, Unwrap(err2))
	}
	if !Is(err2, err) {
		t.Errorf("Is(err2, err) got %v want: %v", Is(err2, err), true)
	}
	err3 := &mockErr{}
	if !As(err2, &err3) {
		t.Errorf("As(err2, &err3) got %v want: %v", As(err2, &err3), true)
	}
}

func TestAsLast(t *testing.T) {
	err := fmt.Errorf("first orignal error")
	err1 := New(1, "reason1", "message1").WithCause(err)
	err2 := fmt.Errorf("wrap2 %w", err1)
	err3 := New(3, "reason3", "message3").WithCause(err2)
	err4 := fmt.Errorf("wrap2 %w", err3)
	errLast := &Error{}

	if !AsLast(err4, &errLast) || errLast != err1 {
		t.Errorf("As(err4, &errLast) got %v want: %v", AsLast(err4, &errLast), true)
	}

}
