package csvparser

import "reflect"

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
