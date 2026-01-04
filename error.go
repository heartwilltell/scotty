package scotty

import "fmt"

// Error is a custom error type for scotty errors.
type Error string

func (e Error) Error() string { return string(e) }

// ErrRequiredField is returned when a required field is not set.
const ErrRequiredField Error = "required field not set"

// RequiredFieldError provides details about which required field was not set.
type RequiredFieldError struct {
	FieldName string
	FlagName  string
	EnvName   string
}

func (e *RequiredFieldError) Error() string {
	return fmt.Sprintf("%s: field=%s, flag=%s, env=%s", ErrRequiredField, e.FieldName, e.FlagName, e.EnvName)
}

func (e *RequiredFieldError) Unwrap() error {
	return ErrRequiredField
}
