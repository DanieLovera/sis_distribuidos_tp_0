package util

import (
	"reflect"
	"unsafe"
)

func SizeOfField(t any, fieldName string) int {
	value := reflect.ValueOf(t)
	field := value.FieldByName(fieldName)
	return int(field.Type().Size())
}

func SizeOfType(t any) int {
	return int(unsafe.Sizeof(t))
}
