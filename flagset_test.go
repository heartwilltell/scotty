package scotty

import (
	"reflect"
	"testing"
)

func Test_tern(t *testing.T) {
	type tcase[T any] struct {
		cond       bool
		t, f, want T
	}

	tests := map[string]tcase[string]{
		"true": {
			cond: true,
			t:    "true",
			f:    "false",
			want: "true",
		},

		"false": {
			cond: false,
			t:    "true",
			f:    "false",
			want: "false",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tern(tc.cond, tc.t, tc.f)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("tern got := %v want := %v", got, tc.want)
			}
		})
	}
}
