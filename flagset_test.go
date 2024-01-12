package scotty

import (
	"flag"
	"reflect"
	"testing"
	"time"
)

func TestFlagSet_DurationVarE(t *testing.T) {
	type tcase[T any] struct {
		got, want time.Duration
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[time.Duration]{
		"Flag": {
			want: 11 * time.Second,
			before: func(t *testing.T, f *FlagSet, got *time.Duration) {
				t.Helper()

				f.DurationVarE(got, "f1", "TEST_E1", 10*time.Second, "")
			},
			args: []string{"-f1=11s"},
		},

		"Env": {
			want: 11 * time.Second,
			before: func(t *testing.T, f *FlagSet, got *time.Duration) {
				t.Helper()
				t.Setenv("TEST_E1", "11s")
				f.DurationVarE(got, "f1", "TEST_E1", 10*time.Second, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 10 * time.Second,
			before: func(t *testing.T, f *FlagSet, got *time.Duration) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.DurationVarE(got, "f1", "TEST_E1", 10*time.Second, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 11 * time.Second,
			before: func(t *testing.T, f *FlagSet, got *time.Duration) {
				t.Helper()
				t.Setenv("TEST_E1", "12s")
				f.DurationVarE(got, "f1", "TEST_E1", 10*time.Second, "")
			},
			args: []string{"-f1=11s"},
		},

		"Default": {
			want: 10 * time.Second,
			before: func(t *testing.T, f *FlagSet, got *time.Duration) {
				t.Helper()
				f.DurationVarE(got, "f1", "TEST_E1", 10*time.Second, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_BoolVarE(t *testing.T) {
	type tcase[T any] struct {
		got, want bool
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[bool]{
		"Flag": {
			want: true,
			before: func(t *testing.T, f *FlagSet, got *bool) {
				t.Helper()

				f.BoolVarE(got, "f1", "TEST_E1", false, "")
			},
			args: []string{"-f1"},
		},

		"Env": {
			want: true,
			before: func(t *testing.T, f *FlagSet, got *bool) {
				t.Helper()
				t.Setenv("TEST_E1", "true")
				f.BoolVarE(got, "f1", "TEST_E1", false, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: false,
			before: func(t *testing.T, f *FlagSet, got *bool) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.BoolVarE(got, "f1", "TEST_E1", false, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: true,
			before: func(t *testing.T, f *FlagSet, got *bool) {
				t.Helper()
				t.Setenv("TEST_E1", "false")
				f.BoolVarE(got, "f1", "TEST_E1", false, "")
			},
			args: []string{"-f1"},
		},

		"Default": {
			want: false,
			before: func(t *testing.T, f *FlagSet, got *bool) {
				t.Helper()
				f.BoolVarE(got, "f1", "TEST_E1", false, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_Float64VarE(t *testing.T) {
	type tcase[T any] struct {
		got, want float64
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[float64]{
		"Flag": {
			want: 10.2,
			before: func(t *testing.T, f *FlagSet, got *float64) {
				t.Helper()

				f.Float64VarE(got, "f1", "TEST_E1", 10.1, "")
			},
			args: []string{"-f1=10.2"},
		},

		"Env": {
			want: 10.2,
			before: func(t *testing.T, f *FlagSet, got *float64) {
				t.Helper()
				t.Setenv("TEST_E1", "10.2")
				f.Float64VarE(got, "f1", "TEST_E1", 10.1, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 10.1,
			before: func(t *testing.T, f *FlagSet, got *float64) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.Float64VarE(got, "f1", "TEST_E1", 10.1, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 10.2,
			before: func(t *testing.T, f *FlagSet, got *float64) {
				t.Helper()
				t.Setenv("TEST_E1", "10.3")
				f.Float64VarE(got, "f1", "TEST_E1", 10.1, "")
			},
			args: []string{"-f1=10.2"},
		},

		"Default": {
			want: 10.1,
			before: func(t *testing.T, f *FlagSet, got *float64) {
				t.Helper()
				f.Float64VarE(got, "f1", "TEST_E1", 10.1, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_Int64VarE(t *testing.T) {
	type tcase[T any] struct {
		got, want int64
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[int64]{
		"Flag": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int64) {
				t.Helper()

				f.Int64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Env": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int64) {
				t.Helper()
				t.Setenv("TEST_E1", "10")
				f.Int64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *int64) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.Int64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int64) {
				t.Helper()
				t.Setenv("TEST_E1", "12")
				f.Int64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Default": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *int64) {
				t.Helper()
				f.Int64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_IntVarE(t *testing.T) {
	type tcase[T any] struct {
		got, want int
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[int]{
		"Flag": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int) {
				t.Helper()

				f.IntVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Env": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int) {
				t.Helper()
				t.Setenv("TEST_E1", "10")
				f.IntVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *int) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.IntVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *int) {
				t.Helper()
				t.Setenv("TEST_E1", "12")
				f.IntVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Default": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *int) {
				t.Helper()
				f.IntVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_StringVarE(t *testing.T) {
	type tcase[T any] struct {
		got, want string
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[string]{
		"Flag": {
			want: "plainq",
			before: func(t *testing.T, f *FlagSet, got *string) {
				t.Helper()

				f.StringVarE(got, "f1", "TEST_E1", "", "")
			},
			args: []string{"-f1=plainq"},
		},

		"Env": {
			want: "plainq",
			before: func(t *testing.T, f *FlagSet, got *string) {
				t.Helper()
				t.Setenv("TEST_E1", "plainq")
				f.StringVarE(got, "f1", "TEST_E1", "", "")
			},
			args: []string{},
		},

		"BothSet": {
			want: "plainq",
			before: func(t *testing.T, f *FlagSet, got *string) {
				t.Helper()
				t.Setenv("TEST_E1", "plainq2")
				f.StringVarE(got, "f1", "TEST_E1", "", "")
			},
			args: []string{"-f1=plainq"},
		},

		"Default": {
			want: "plainq",
			before: func(t *testing.T, f *FlagSet, got *string) {
				t.Helper()
				f.StringVarE(got, "f1", "TEST_E1", "plainq", "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_Uint64VarE(t *testing.T) {
	type tcase[T any] struct {
		got, want uint64
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[uint64]{
		"Flag": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint64) {
				t.Helper()

				f.Uint64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Env": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint64) {
				t.Helper()
				t.Setenv("TEST_E1", "10")
				f.Uint64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *uint64) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.Uint64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint64) {
				t.Helper()
				t.Setenv("TEST_E1", "12")
				f.Uint64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Default": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *uint64) {
				t.Helper()
				f.Uint64VarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

func TestFlagSet_UintVarE(t *testing.T) {
	type tcase[T any] struct {
		got, want uint
		before    func(t *testing.T, f *FlagSet, got *T)
		args      []string
	}

	tests := map[string]tcase[uint]{
		"Flag": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint) {
				t.Helper()

				f.UintVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Env": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint) {
				t.Helper()
				t.Setenv("TEST_E1", "10")
				f.UintVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"InvalidEnv": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *uint) {
				t.Helper()
				t.Setenv("TEST_E1", "lalala")
				f.UintVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},

		"BothSet": {
			want: 10,
			before: func(t *testing.T, f *FlagSet, got *uint) {
				t.Helper()
				t.Setenv("TEST_E1", "12")
				f.UintVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{"-f1=10"},
		},

		"Default": {
			want: 11,
			before: func(t *testing.T, f *FlagSet, got *uint) {
				t.Helper()
				f.UintVarE(got, "f1", "TEST_E1", 11, "")
			},
			args: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := &FlagSet{
				FlagSet: flag.NewFlagSet("test", flag.ExitOnError),
			}

			tc.before(t, f, &tc.got)

			if err := f.Parse(tc.args); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.want, tc.got) {
				t.Errorf("want := %v, got := %v", tc.want, tc.got)
			}
		})
	}
}

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
