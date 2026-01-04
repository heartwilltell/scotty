package scotty

import (
	"errors"
	"os"
	"testing"
	"time"
)

type testConfig struct {
	Host    string        `flag:"host" env:"TEST_HOST" default:"localhost" usage:"Server host"`
	Port    int           `flag:"port" env:"TEST_PORT" default:"8080" usage:"Server port"`
	Debug   bool          `flag:"debug" env:"TEST_DEBUG" default:"false" usage:"Enable debug"`
	Timeout time.Duration `flag:"timeout" env:"TEST_TIMEOUT" default:"30s" usage:"Timeout"`
	Rate    float64       `flag:"rate" env:"TEST_RATE" default:"1.5" usage:"Rate limit"`
	Count   int64         `flag:"count" env:"TEST_COUNT" default:"100" usage:"Count"`
	Size    uint          `flag:"size" env:"TEST_SIZE" default:"1024" usage:"Size"`
	MaxSize uint64        `flag:"max-size" env:"TEST_MAX_SIZE" default:"65536" usage:"Max size"`
}

type testRequiredConfig struct {
	Name string `flag:"name" env:"TEST_NAME" required:"true" usage:"Required name"`
	Port int    `flag:"port" env:"TEST_PORT" default:"8080" usage:"Port"`
}

type testValidatorConfig struct {
	Port int `flag:"port" env:"TEST_PORT" default:"8080" usage:"Port number"`
}

func (c *testValidatorConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	return nil
}

func TestBindConfig_AllTypes(t *testing.T) {
	cfg := &testConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse with default values.
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check defaults were applied.
	if cfg.Host != "localhost" {
		t.Errorf("Host = %q, want %q", cfg.Host, "localhost")
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want %d", cfg.Port, 8080)
	}

	if cfg.Debug != false {
		t.Errorf("Debug = %v, want %v", cfg.Debug, false)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}

	if cfg.Rate != 1.5 {
		t.Errorf("Rate = %f, want %f", cfg.Rate, 1.5)
	}

	if cfg.Count != 100 {
		t.Errorf("Count = %d, want %d", cfg.Count, 100)
	}

	if cfg.Size != 1024 {
		t.Errorf("Size = %d, want %d", cfg.Size, 1024)
	}

	if cfg.MaxSize != 65536 {
		t.Errorf("MaxSize = %d, want %d", cfg.MaxSize, 65536)
	}
}

func TestBindConfig_FlagOverride(t *testing.T) {
	cfg := &testConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse with flag values.
	if err := cmd.Flags().Parse([]string{"-host=0.0.0.0", "-port=3000", "-debug=true"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Host != "0.0.0.0" {
		t.Errorf("Host = %q, want %q", cfg.Host, "0.0.0.0")
	}

	if cfg.Port != 3000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 3000)
	}

	if cfg.Debug != true {
		t.Errorf("Debug = %v, want %v", cfg.Debug, true)
	}
}

func TestBindConfig_EnvFallback(t *testing.T) {
	// Set environment variables.
	os.Setenv("TEST_HOST", "envhost")
	os.Setenv("TEST_PORT", "9000")

	defer func() {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_PORT")
	}()

	cfg := &testConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse without flags - should use env vars.
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if cfg.Host != "envhost" {
		t.Errorf("Host = %q, want %q", cfg.Host, "envhost")
	}

	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want %d", cfg.Port, 9000)
	}
}

func TestBindConfig_RequiredField_Missing(t *testing.T) {
	cfg := &testRequiredConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse without required field - should fail during exec.
	if err := cmd.Flags().Parse([]string{}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// execCommand should fail due to missing required field.
	err := cmd.execCommand([]string{})
	if err == nil {
		t.Fatal("expected error for missing required field, got nil")
	}

	if !errors.Is(err, ErrRequiredField) {
		t.Errorf("error = %v, want error wrapping ErrRequiredField", err)
	}
}

func TestBindConfig_RequiredField_Provided(t *testing.T) {
	cfg := &testRequiredConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse with required field.
	if err := cmd.Flags().Parse([]string{"-name=myapp"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// execCommand should succeed.
	err := cmd.execCommand([]string{"-name=myapp"})
	if err != nil {
		t.Fatalf("execCommand failed: %v", err)
	}

	if cfg.Name != "myapp" {
		t.Errorf("Name = %q, want %q", cfg.Name, "myapp")
	}
}

func TestBindConfig_Validator_Valid(t *testing.T) {
	cfg := &testValidatorConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse with valid port.
	if err := cmd.Flags().Parse([]string{"-port=8080"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	err := cmd.execCommand([]string{"-port=8080"})
	if err != nil {
		t.Fatalf("execCommand failed: %v", err)
	}
}

func TestBindConfig_Validator_Invalid(t *testing.T) {
	cfg := &testValidatorConfig{}

	cmd := &Command{
		Name: "test",
		Run: func(cmd *Command, args []string) error {
			return nil
		},
	}

	if err := cmd.BindConfig(cfg); err != nil {
		t.Fatalf("BindConfig failed: %v", err)
	}

	// Parse with invalid port (0).
	if err := cmd.Flags().Parse([]string{"-port=0"}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	err := cmd.execCommand([]string{"-port=0"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if !contains(err.Error(), "validation failed") {
		t.Errorf("error = %v, want error containing 'validation failed'", err)
	}
}

func TestBindConfig_InvalidInput(t *testing.T) {
	type tcase struct {
		cfg any
	}

	tests := map[string]tcase{
		"Non-pointer": {
			cfg: testConfig{},
		},
		"Non-struct pointer": {
			cfg: new(int),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &Command{Name: "test"}
			if err := cmd.BindConfig(tc.cfg); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestMustConfig(t *testing.T) {
	cfg := &testConfig{}
	cmd := &Command{Name: "test"}
	_ = cmd.BindConfig(cfg)

	got := MustConfig[testConfig](cmd)
	if got != cfg {
		t.Errorf("MustConfig returned %p, want %p", got, cfg)
	}
}

func TestGetConfig(t *testing.T) {
	type tcase struct {
		before func(cmd *Command)
		want   any
		wantOk bool
	}

	cfg := &testConfig{}

	tests := map[string]tcase{
		"OK": {
			before: func(cmd *Command) {
				_ = cmd.BindConfig(cfg)
			},
			want:   cfg,
			wantOk: true,
		},
		"Wrong type": {
			before: func(cmd *Command) {
				_ = cmd.BindConfig(cfg)
			},
			want:   nil,
			wantOk: false,
		},
		"No bound config": {
			before: func(cmd *Command) {
				// No binding
			},
			want:   nil,
			wantOk: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &Command{Name: "test"}
			tc.before(cmd)

			if name == "Wrong type" {
				got, ok := GetConfig[testRequiredConfig](cmd)
				if ok != tc.wantOk {
					t.Errorf("GetConfig ok = %v, want %v", ok, tc.wantOk)
				}
				if got != nil {
					t.Errorf("GetConfig got = %v, want nil", got)
				}
			} else {
				got, ok := GetConfig[testConfig](cmd)
				if ok != tc.wantOk {
					t.Errorf("GetConfig ok = %v, want %v", ok, tc.wantOk)
				}
				if tc.want == nil {
					if got != nil {
						t.Errorf("GetConfig got = %v, want nil", got)
					}
				} else if got != tc.want {
					t.Errorf("GetConfig got = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}

	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
