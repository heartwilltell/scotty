package scotty

import (
	"flag"
	"testing"
)

func Test_hasFlags(t *testing.T) {
	type tcase struct {
		flags *FlagSet
		want  bool
	}

	tests := map[string]tcase{
		"Has": {
			flags: func() *FlagSet {
				set := &FlagSet{FlagSet: flag.NewFlagSet("test", flag.ExitOnError)}
				set.String("test-flag", "", "")
				return set
			}(),

			want: true,
		},
		"Doesn't": {
			flags: func() *FlagSet {
				set := &FlagSet{FlagSet: flag.NewFlagSet("test", flag.ExitOnError)}
				return set
			}(),

			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := hasFlags(tc.flags); got != tc.want {
				t.Errorf("Expected := %v, got := %v", tc.want, got)
			}
		})
	}
}
