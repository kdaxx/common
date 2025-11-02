package ref

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

func ReflectField[T any](structPtr any, fieldName string) (*T, error) {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("can not reflect non struct pointer")
	}

	elem := v.Elem()
	field := elem.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("can't find field %s", fieldName)
	}

	// ensure addressable
	if !field.CanAddr() {
		return nil, fmt.Errorf("field %s is not addressable", fieldName)
	}

	return (*T)(unsafe.Pointer(field.UnsafeAddr())), nil
}
