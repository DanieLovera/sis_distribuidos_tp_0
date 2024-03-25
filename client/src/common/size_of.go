package common

import (
	"reflect"
	"unsafe"
)

func SizeOfField(t interface{}, fieldName string) int {
	value := reflect.ValueOf(t)
	field := value.FieldByName(fieldName)
	return int(field.Type().Size())
}

func SizeOfType(t interface{}) int {
	return int(unsafe.Sizeof(t))
}
