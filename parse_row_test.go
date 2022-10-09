package csvparser

import (
	"reflect"
	"testing"
)

func Test_setVal(t *testing.T) {
	type args struct {
		rv  *reflect.Value
		val reflect.Value
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setVal(tt.args.rv, tt.args.val)
		})
	}
}
