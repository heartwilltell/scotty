package scotty

import (
	"fmt"
	"os"
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

		if err := bindField(f, fieldVal, flagName, envName, defaultVal, usage); err != nil {
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
//nolint:cyclop // Switch on types is inherently complex.
func bindField(f *FlagSet, fieldVal reflect.Value, flagName, envName, defaultVal, usage string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		ptr := fieldVal.Addr().Interface().(*string)
		if envName != "" {
			f.StringVarE(ptr, flagName, envName, defaultVal, usage)
		} else {
			f.StringVar(ptr, flagName, defaultVal, usage)
		}

	case reflect.Bool:
		ptr := fieldVal.Addr().Interface().(*bool)
		def, _ := strconv.ParseBool(defaultVal)
		if envName != "" {
			f.BoolVarE(ptr, flagName, envName, def, usage)
		} else {
			f.BoolVar(ptr, flagName, def, usage)
		}

	case reflect.Int:
		ptr := fieldVal.Addr().Interface().(*int)
		def, _ := strconv.Atoi(defaultVal)
		if envName != "" {
			f.IntVarE(ptr, flagName, envName, def, usage)
		} else {
			f.IntVar(ptr, flagName, def, usage)
		}

	case reflect.Int64:
		// Check for time.Duration specifically.
		if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
			ptr := fieldVal.Addr().Interface().(*time.Duration)
			def, _ := time.ParseDuration(defaultVal)
			if envName != "" {
				f.DurationVarE(ptr, flagName, envName, def, usage)
			} else {
				f.DurationVar(ptr, flagName, def, usage)
			}
		} else {
			ptr := fieldVal.Addr().Interface().(*int64)
			def, _ := strconv.ParseInt(defaultVal, 10, 64)
			if envName != "" {
				f.Int64VarE(ptr, flagName, envName, def, usage)
			} else {
				f.Int64Var(ptr, flagName, def, usage)
			}
		}

	case reflect.Uint:
		ptr := fieldVal.Addr().Interface().(*uint)
		def, _ := strconv.ParseUint(defaultVal, 10, 64)
		if envName != "" {
			f.UintVarE(ptr, flagName, envName, uint(def), usage)
		} else {
			f.UintVar(ptr, flagName, uint(def), usage)
		}

	case reflect.Uint64:
		ptr := fieldVal.Addr().Interface().(*uint64)
		def, _ := strconv.ParseUint(defaultVal, 10, 64)
		if envName != "" {
			f.Uint64VarE(ptr, flagName, envName, def, usage)
		} else {
			f.Uint64Var(ptr, flagName, def, usage)
		}

	case reflect.Float64:
		ptr := fieldVal.Addr().Interface().(*float64)
		def, _ := strconv.ParseFloat(defaultVal, 64)
		if envName != "" {
			f.Float64VarE(ptr, flagName, envName, def, usage)
		} else {
			f.Float64Var(ptr, flagName, def, usage)
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

// getEnvOrDefault returns the environment variable value if set, otherwise the default.
func getEnvOrDefault(envName, defaultVal string) string {
	if envName == "" {
		return defaultVal
	}

	if val := os.Getenv(envName); val != "" {
		return val
	}

	return defaultVal
}
