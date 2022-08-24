package scotty

import (
	"flag"
	"fmt"
	"time"
)

// flags represents a thin wrapper around standard flag.FlagSet,
// which hides verbose api.
type flags struct {
	name  string
	flags *flag.FlagSet
}

func (f *flags) parse(args []string) error {
	if len(args) == 0 {
		return nil
	}

	if err := f.flags.Parse(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	return nil
}

func (f *flags) Bool(p *bool, name string, value bool, usage string) {
	f.flags.BoolVar(p, name, value, usage)
}

func (f *flags) Uint(p *uint, name string, value uint, usage string) {
	f.flags.UintVar(p, name, value, usage)
}

func (f *flags) Uint64(p *uint64, name string, value uint64, usage string) {
	f.flags.Uint64Var(p, name, value, usage)
}

func (f *flags) Int(p *int, name string, value int, usage string) {
	f.flags.IntVar(p, name, value, usage)
}

func (f *flags) Int64(p *int64, name string, value int64, usage string) {
	f.flags.Int64Var(p, name, value, usage)
}

func (f *flags) Float64(p *float64, name string, value float64, usage string) {
	f.flags.Float64Var(p, name, value, usage)
}

func (f *flags) String(p *string, name, value, usage string) {
	f.flags.StringVar(p, name, value, usage)
}

func (f *flags) Duration(p *time.Duration, name string, value time.Duration, usage string) {
	f.flags.DurationVar(p, name, value, usage)
}
