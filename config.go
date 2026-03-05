package scotty

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Struct tag names for config binding.
const (
	tagFlag     = "flag"
	tagEnv      = "env"
	tagDefault  = "default"
	tagUsage    = "usage"
	tagRequired = "required"
)

// ConfigValidator holds logic of validation the config parameters.
type ConfigValidator interface {
	// Validate validates the config parameters.
	Validate() error
}

// MustConfig returns the bound config with type assertion.
// Panics if config is not bound or wrong type.
func MustConfig[T any](cmd *Command) *T {
	cfg, ok := GetConfig[T](cmd)
	if !ok {
		panic("config not bound or wrong type")
	}

	return cfg
}

// GetConfig returns the bound config with type assertion.
// Returns nil and false if config is not bound or wrong type.
func GetConfig[T any](cmd *Command) (*T, bool) {
	if cmd.Config() == nil {
		return nil, false
	}

	cfg, ok := cmd.Config().(*T)

	return cfg, ok
}

// requiredFieldInfo tracks a required field for validation after parsing.
type requiredFieldInfo struct {
	fieldName string
	flagName  string
	envName   string
	fieldPtr  reflect.Value
}

// fieldOpts holds options for binding a struct field to a flag.
type fieldOpts struct {
	flagName   string
	envName    string
	defaultVal string
	usage      string
}

// bindConfigToFlagSet uses reflection to bind struct fields to flags.
// cfg must be a pointer to a struct.
func bindConfigToFlagSet(f *FlagSet, cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("config must be a pointer, got %T", cfg)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct, got pointer to %s", v.Kind())
	}

	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Skip unexported fields.
		if !fieldVal.CanSet() {
			continue
		}

		flagName := field.Tag.Get(tagFlag)
		if flagName == "" {
			continue // Skip fields without flag tag.
		}

		envName := field.Tag.Get(tagEnv)
		defaultVal := field.Tag.Get(tagDefault)
		usage := field.Tag.Get(tagUsage)
		required := field.Tag.Get(tagRequired) == "true"

		opts := fieldOpts{
			flagName:   flagName,
			envName:    envName,
			defaultVal: defaultVal,
			usage:      usage,
		}

		if err := bindField(f, fieldVal, opts); err != nil {
			return fmt.Errorf("binding field %s: %w", field.Name, err)
		}

		if required {
			f.requiredFields = append(f.requiredFields, requiredFieldInfo{
				fieldName: field.Name,
				flagName:  flagName,
				envName:   envName,
				fieldPtr:  fieldVal,
			})
		}
	}

	return nil
}

// bindField binds a single struct field to a flag based on its type.
//
//nolint:cyclop,revive // Switch on types is inherently complex.
func bindField(f *FlagSet, fieldVal reflect.Value, opts fieldOpts) error {
	switch fieldVal.Kind() {
	case reflect.String:
		ptr, ok := fieldVal.Addr().Interface().(*string)
		if !ok {
			return fmt.Errorf("invalid type for string field")
		}
		if opts.envName != "" {
			f.StringVarE(ptr, opts.flagName, opts.envName, opts.defaultVal, opts.usage)
		} else {
			f.StringVar(ptr, opts.flagName, opts.defaultVal, opts.usage)
		}

	case reflect.Bool:
		ptr, ok := fieldVal.Addr().Interface().(*bool)
		if !ok {
			return fmt.Errorf("invalid type for bool field")
		}
		parsed, err := strconv.ParseBool(opts.defaultVal)
		def := tern(err == nil, parsed, false)
		if opts.envName != "" {
			f.BoolVarE(ptr, opts.flagName, opts.envName, def, opts.usage)
		} else {
			f.BoolVar(ptr, opts.flagName, def, opts.usage)
		}

	case reflect.Int:
		ptr, ok := fieldVal.Addr().Interface().(*int)
		if !ok {
			return fmt.Errorf("invalid type for int field")
		}
		parsed, err := strconv.Atoi(opts.defaultVal)
		def := tern(err == nil, parsed, 0)
		if opts.envName != "" {
			f.IntVarE(ptr, opts.flagName, opts.envName, def, opts.usage)
		} else {
			f.IntVar(ptr, opts.flagName, def, opts.usage)
		}

	case reflect.Int64:
		if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
			ptr, ok := fieldVal.Addr().Interface().(*time.Duration)
			if !ok {
				return fmt.Errorf("invalid type for duration field")
			}
			parsed, err := time.ParseDuration(opts.defaultVal)
			def := tern(err == nil, parsed, 0)
			if opts.envName != "" {
				f.DurationVarE(ptr, opts.flagName, opts.envName, def, opts.usage)
			} else {
				f.DurationVar(ptr, opts.flagName, def, opts.usage)
			}
		} else {
			ptr, ok := fieldVal.Addr().Interface().(*int64)
			if !ok {
				return fmt.Errorf("invalid type for int64 field")
			}
			parsed, err := strconv.ParseInt(opts.defaultVal, 10, 64)
			def := tern(err == nil, parsed, int64(0))
			if opts.envName != "" {
				f.Int64VarE(ptr, opts.flagName, opts.envName, def, opts.usage)
			} else {
				f.Int64Var(ptr, opts.flagName, def, opts.usage)
			}
		}

	case reflect.Uint:
		ptr, ok := fieldVal.Addr().Interface().(*uint)
		if !ok {
			return fmt.Errorf("invalid type for uint field")
		}
		parsed, err := strconv.ParseUint(opts.defaultVal, 10, 64)
		def := tern(err == nil, uint(parsed), uint(0))
		if opts.envName != "" {
			f.UintVarE(ptr, opts.flagName, opts.envName, def, opts.usage)
		} else {
			f.UintVar(ptr, opts.flagName, def, opts.usage)
		}

	case reflect.Uint64:
		ptr, ok := fieldVal.Addr().Interface().(*uint64)
		if !ok {
			return fmt.Errorf("invalid type for uint64 field")
		}
		parsed, err := strconv.ParseUint(opts.defaultVal, 10, 64)
		def := tern(err == nil, parsed, uint64(0))
		if opts.envName != "" {
			f.Uint64VarE(ptr, opts.flagName, opts.envName, def, opts.usage)
		} else {
			f.Uint64Var(ptr, opts.flagName, def, opts.usage)
		}

	case reflect.Float64:
		ptr, ok := fieldVal.Addr().Interface().(*float64)
		if !ok {
			return fmt.Errorf("invalid type for float64 field")
		}
		parsed, err := strconv.ParseFloat(opts.defaultVal, 64)
		def := tern(err == nil, parsed, 0.0)
		if opts.envName != "" {
			f.Float64VarE(ptr, opts.flagName, opts.envName, def, opts.usage)
		} else {
			f.Float64Var(ptr, opts.flagName, def, opts.usage)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
	}

	return nil
}

// validateRequiredFields checks that all required fields have non-zero values.
func validateRequiredFields(fields []requiredFieldInfo) error {
	for _, f := range fields {
		if isZeroValue(f.fieldPtr) {
			return &RequiredFieldError{
				FieldName: f.fieldName,
				FlagName:  f.flagName,
				EnvName:   f.envName,
			}
		}
	}

	return nil
}

// isZeroValue checks if a reflect.Value holds the zero value for its type.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return false
	}
}

// validateConfig calls Validate() if cfg implements ConfigValidator.
func validateConfig(cfg any) error {
	if validator, ok := cfg.(ConfigValidator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	return nil
}
