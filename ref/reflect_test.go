package ref

import "testing"

type structWrapper struct {
	nonExported int
}

func TestReflect(t *testing.T) {
	s := &structWrapper{}
	field, err := ReflectField[int](s, "nonExported")
	if err != nil {
		t.Fatal("failed to ReflectField[int]", err)
	}
	*field = 1
	if s.nonExported != 1 {
		t.Fatal("failed to ReflectField[int]")
	}
}
