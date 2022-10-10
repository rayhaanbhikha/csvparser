package csvparser

import (
	"reflect"
)

func isKind(t reflect.Type, k2 reflect.Kind) bool {
	switch k1 := t.Kind(); k1 {
	case k2:
		return true
	case reflect.Pointer:
		return isKind(t.Elem(), k2)
	default:
		return false
	}
}

func extractStruct(val reflect.Value) (reflect.Value, error) {
	switch k1 := val.Kind(); k1 {
	case reflect.Struct:
		return val, nil
	case reflect.Pointer:
		v1 := val.Elem()
		return extractStruct(v1)
	default:
		return reflect.Value{}, ErrCSVRowMustBeAStruct
	}
}

func setVal[T any](rv *reflect.Value, actualVal T) {
	val := reflect.ValueOf(actualVal)
	if !rv.CanSet() {
		return
	}

	if rv.Kind() == reflect.Pointer {
		if val.IsZero() {
			return
		}
		rv.Set(reflect.ValueOf(&actualVal))
		return
	}
	rv.Set(val)
}
