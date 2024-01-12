package scotty

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// FlagSet wraps flag.FlagSet and adds a few methods like
// StringVarE, BoolVarE and similar methods for other types.
type FlagSet struct{ *flag.FlagSet }

// StringVarE defines a string flag and environment variable with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) StringVarE(p *string, flagName, envName, value, usage string) {
	f.FlagSet.StringVar(p, flagName, tern[string](os.Getenv(envName) != "", os.Getenv(envName), value), usage)
}

// BoolVarE defines a bool flag and environment variable with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) BoolVarE(p *bool, flagName, envName string, value bool, usage string) {
	parsed, err := strconv.ParseBool(os.Getenv(envName))
	f.FlagSet.BoolVar(p, flagName, tern[bool](err == nil, parsed, value), usage)
}

// IntVarE defines an int flag and environment variable with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) IntVarE(p *int, flagName, envName string, value int, usage string) {
	parsed, err := strconv.Atoi(os.Getenv(envName))
	f.FlagSet.IntVar(p, flagName, tern[int](err == nil, parsed, value), usage)
}

// Int64VarE defines an int64 flag and environment variable with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) Int64VarE(p *int64, flagName, envName string, value int64, usage string) {
	parsed, err := strconv.Atoi(os.Getenv(envName))
	f.FlagSet.Int64Var(p, flagName, tern[int64](err == nil, int64(parsed), value), usage)
}

// Float64VarE defines a float64 flag and environment variable with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) Float64VarE(p *float64, flagName, envName string, value float64, usage string) {
	parsed, err := strconv.ParseFloat(os.Getenv(envName), 64)
	f.FlagSet.Float64Var(p, flagName, tern[float64](err == nil, parsed, value), usage)
}

// UintVarE defines an uint flag and environment variable with specified name, default value, and usage string.
// The argument p points to an uint variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) UintVarE(p *uint, flagName, envName string, value uint, usage string) {
	parsed, err := strconv.Atoi(os.Getenv(envName))
	f.FlagSet.UintVar(p, flagName, tern[uint](err == nil, uint(parsed), value), usage)
}

// Uint64VarE defines an uint64 flag and environment variable with specified name, default value, and usage string.
// The argument p points to an uint64 variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) Uint64VarE(p *uint64, flagName, envName string, value uint64, usage string) {
	parsed, err := strconv.Atoi(os.Getenv(envName))
	f.FlagSet.Uint64Var(p, flagName, tern[uint64](err == nil, uint64(parsed), value), usage)
}

// DurationVarE defines a time.Duration flag and environment variable with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag or environment variable.
// Flag has priority over environment variable. If flag not set the environment variable value will be used.
// If the value of environment variable can't be parsed to destination type the default value will be used.
func (f *FlagSet) DurationVarE(p *time.Duration, flagName, envName string, value time.Duration, usage string) {
	parsed, err := time.ParseDuration(os.Getenv(envName))
	fmt.Println(parsed, err, tern[time.Duration](err == nil, parsed, value))
	f.FlagSet.DurationVar(p, flagName, tern[time.Duration](err == nil, parsed, value), usage)
}

//nolint:revive // flag-parameter is ok here.
func tern[T any](cond bool, t, f T) T {
	if cond {
		return t
	}

	return f
}
