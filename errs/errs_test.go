package errs

import (
	"errors"
	"fmt"
	"testing"
)

func TestCombine(t *testing.T) {
	test1 := errors.New("test1")
	test2 := errors.New("test2")
	err := Combine(test1, test2)

	exp := fmt.Sprintf("%s; %s", test1, test2)
	if !(fmt.Sprintf("%s", err) == exp) {
		t.Fatalf("Combine() error = %s, expected %s", err, exp)
	}

	if !(errors.Is(err, test1) && errors.Is(err, test2)) {
		t.Fatalf("Combine() error does not compatible with errors.Is()")
	}

}
