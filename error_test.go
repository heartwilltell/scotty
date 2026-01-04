package scotty

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	type tcase struct {
		err  Error
		want string
	}

	tests := map[string]tcase{
		"ErrRequiredField": {
			err:  ErrRequiredField,
			want: "required field not set",
		},
		"CustomError": {
			err:  Error("custom error message"),
			want: "custom error message",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("Error.Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRequiredFieldError_Error(t *testing.T) {
	err := &RequiredFieldError{
		FieldName: "Host",
		FlagName:  "host",
		EnvName:   "APP_HOST",
	}

	want := "required field not set: field=Host, flag=host, env=APP_HOST"
	if got := err.Error(); got != want {
		t.Errorf("RequiredFieldError.Error() = %q, want %q", got, want)
	}
}

func TestRequiredFieldError_Unwrap(t *testing.T) {
	err := &RequiredFieldError{
		FieldName: "Port",
		FlagName:  "port",
		EnvName:   "APP_PORT",
	}

	unwrapped := err.Unwrap()
	if unwrapped != ErrRequiredField {
		t.Errorf("RequiredFieldError.Unwrap() = %v, want %v", unwrapped, ErrRequiredField)
	}
}

func TestRequiredFieldError_Is(t *testing.T) {
	err := &RequiredFieldError{
		FieldName: "Name",
		FlagName:  "name",
		EnvName:   "APP_NAME",
	}

	if !errors.Is(err, ErrRequiredField) {
		t.Error("errors.Is(RequiredFieldError, ErrRequiredField) = false, want true")
	}

	if errors.Is(err, Error("other error")) {
		t.Error("errors.Is(RequiredFieldError, other error) = true, want false")
	}
}

func TestRequiredFieldError_As(t *testing.T) {
	err := &RequiredFieldError{
		FieldName: "Debug",
		FlagName:  "debug",
		EnvName:   "APP_DEBUG",
	}

	var reqErr *RequiredFieldError
	if !errors.As(err, &reqErr) {
		t.Fatal("errors.As failed to match RequiredFieldError")
	}

	if reqErr.FieldName != "Debug" {
		t.Errorf("FieldName = %q, want %q", reqErr.FieldName, "Debug")
	}

	if reqErr.FlagName != "debug" {
		t.Errorf("FlagName = %q, want %q", reqErr.FlagName, "debug")
	}

	if reqErr.EnvName != "APP_DEBUG" {
		t.Errorf("EnvName = %q, want %q", reqErr.EnvName, "APP_DEBUG")
	}
}
