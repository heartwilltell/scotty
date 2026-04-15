package scotty

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseDotenv parses a dotenv-formatted stream and returns the parsed
// key/value pairs. Supports comments (#), blank lines, optional "export"
// prefix, single- and double-quoted values, escape sequences in
// double-quoted values, and ${VAR}/$VAR expansion in unquoted and
// double-quoted values.
func ParseDotenv(r io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseDotenvLine(line, result)
		if err != nil {
			return nil, fmt.Errorf("dotenv: line %d: %w", lineNum, err)
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv: reading: %w", err)
	}

	return result, nil
}

// LoadDotenv parses dotenv files at the given paths and sets the resulting
// variables into the process environment without overriding keys that are
// already set.
func LoadDotenv(paths ...string) error {
	return loadDotenv(false, paths)
}

// LoadDotenvOverride parses dotenv files at the given paths and sets the
// resulting variables into the process environment, overriding any existing
// values.
func LoadDotenvOverride(paths ...string) error {
	return loadDotenv(true, paths)
}

// loadDotenv is the shared implementation behind LoadDotenv variants.
func loadDotenv(override bool, paths []string) error {
	for _, path := range paths {
		vars, err := parseDotenvFile(path)
		if err != nil {
			return err
		}

		for k, v := range vars {
			if !override {
				if _, ok := os.LookupEnv(k); ok {
					continue
				}
			}

			if err := os.Setenv(k, v); err != nil {
				return fmt.Errorf("dotenv: setenv %s: %w", k, err)
			}
		}
	}

	return nil
}

// parseDotenvFile opens the file at path and parses it as dotenv.
func parseDotenvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("dotenv: open %s: %w", path, err)
	}
	defer file.Close()

	return ParseDotenv(file)
}

// parseDotenvLine parses a single non-empty, non-comment line into a
// key/value pair, using existing for variable expansion.
func parseDotenvLine(line string, existing map[string]string) (string, string, error) {
	trimmed := strings.TrimSpace(strings.TrimPrefix(line, "export "))

	eq := strings.IndexByte(trimmed, '=')
	if eq <= 0 {
		return "", "", ErrInvalidLineEqualSign
	}

	key := strings.TrimSpace(trimmed[:eq])
	if !isValidDotenvKey(key) {
		return "", "", fmt.Errorf("invalid key: %q", key)
	}

	value, err := parseDotenvValue(strings.TrimLeft(trimmed[eq+1:], " \t"), existing)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

// isValidDotenvKey reports whether k is a valid environment variable name.
func isValidDotenvKey(k string) bool {
	if k == "" {
		return false
	}

	for i, r := range k {
		switch {
		case r == '_':
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}

	return true
}

// parseDotenvValue parses the value portion of a dotenv line.
func parseDotenvValue(raw string, existing map[string]string) (string, error) {
	if raw == "" {
		return "", nil
	}

	switch raw[0] {
	case '\'':
		end := strings.IndexByte(raw[1:], '\'')
		if end < 0 {
			return "", ErrUnterminatedSingleQuote
		}

		return raw[1 : 1+end], nil

	case '"':
		return parseDoubleQuoted(raw, existing)

	default:
		return parseUnquoted(raw, existing), nil
	}
}

// parseDoubleQuoted parses a double-quoted value, honoring escape sequences
// and performing variable expansion on the decoded result.
func parseDoubleQuoted(raw string, existing map[string]string) (string, error) {
	var b strings.Builder

	for i := 1; i < len(raw); i++ {
		c := raw[i]
		if c == '"' {
			return expandDotenvVars(b.String(), existing), nil
		}

		if c == '\\' && i+1 < len(raw) {
			if err := b.WriteByte(unescapeDotenvByte(raw[i+1])); err != nil {
				return "", fmt.Errorf("failed to write byte: %w", err)
			}

			i++

			continue
		}

		if err := b.WriteByte(c); err != nil {
			return "", fmt.Errorf("write byte: %w", err)
		}
	}

	return "", ErrUnterminatedDoubleQuote
}

// unescapeDotenvByte returns the decoded byte for a backslash escape.
func unescapeDotenvByte(c byte) byte {
	switch c {
	case 'n':
		return '\n'

	case 'r':
		return '\r'

	case 't':
		return '\t'

	default:
		return c
	}
}

// parseUnquoted parses an unquoted value, stripping inline "# comment"
// suffixes and trailing whitespace, then expanding variables.
func parseUnquoted(raw string, existing map[string]string) string {
	toTrim := raw

	if hash := strings.IndexByte(raw, '#'); hash >= 0 {
		toTrim = raw[:hash]
	}

	return expandDotenvVars(strings.TrimRight(toTrim, " \t"), existing)
}

// expandDotenvVars expands $VAR and ${VAR} references using previously
// parsed values first, then falling back to the process environment.
func expandDotenvVars(s string, existing map[string]string) string {
	return os.Expand(s, func(name string) string {
		if v, ok := existing[name]; ok {
			return v
		}

		return os.Getenv(name)
	})
}
