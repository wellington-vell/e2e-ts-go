package internal

import "reflect"

func GetUnderlyingValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return v.Elem()
	}
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func GetReflectType[T any]() reflect.Type {
	var zero T
	return reflect.TypeOf(zero)
}
