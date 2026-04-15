package scotty

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDotenv(t *testing.T) {
	t.Parallel()

	type tcase struct {
		input   string
		want    map[string]string
		wantErr bool
	}

	tests := map[string]tcase{
		"simple": {
			input: "FOO=bar\nBAZ=qux\n",
			want:  map[string]string{"FOO": "bar", "BAZ": "qux"},
		},
		"comments and blanks": {
			input: "# comment\n\nFOO=bar\n",
			want:  map[string]string{"FOO": "bar"},
		},
		"export prefix": {
			input: "export FOO=bar\n",
			want:  map[string]string{"FOO": "bar"},
		},
		"single quoted preserves": {
			input: "FOO='a b $NOPE'\n",
			want:  map[string]string{"FOO": "a b $NOPE"},
		},
		"double quoted escapes and expands": {
			input: "A=1\nB=\"x=${A}\\n\"\n",
			want:  map[string]string{"A": "1", "B": "x=1\n"},
		},
		"unquoted strips inline comment": {
			input: "FOO=bar # trailing\n",
			want:  map[string]string{"FOO": "bar"},
		},
		"missing equals": {
			input:   "FOO\n",
			wantErr: true,
		},
		"invalid key": {
			input:   "1FOO=bar\n",
			wantErr: true,
		},
		"unterminated single quote": {
			input:   "FOO='bar\n",
			wantErr: true,
		},
		"unterminated double quote": {
			input:   "FOO=\"bar\n",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseDotenv(strings.NewReader(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("len mismatch: got %v want %v", got, tc.want)
			}

			for k, v := range tc.want {
				if got[k] != v {
					t.Errorf("key %q: got %q want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestLoadDotenv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("SCOTTY_DOTENV_NEW=one\nSCOTTY_DOTENV_KEEP=file\n"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	t.Setenv("SCOTTY_DOTENV_KEEP", "env")

	if err := LoadDotenv(path); err != nil {
		t.Fatalf("LoadDotenv: %v", err)
	}

	if got := os.Getenv("SCOTTY_DOTENV_NEW"); got != "one" {
		t.Errorf("NEW: got %q want %q", got, "one")
	}

	if got := os.Getenv("SCOTTY_DOTENV_KEEP"); got != "env" {
		t.Errorf("KEEP should not be overridden: got %q", got)
	}

	if err := LoadDotenvOverride(path); err != nil {
		t.Fatalf("LoadDotenvOverride: %v", err)
	}

	if got := os.Getenv("SCOTTY_DOTENV_KEEP"); got != "file" {
		t.Errorf("KEEP should be overridden: got %q want %q", got, "file")
	}
}

func TestLoadDotenvMissingFile(t *testing.T) {
	t.Parallel()

	if err := LoadDotenv(filepath.Join(t.TempDir(), "does-not-exist.env")); err == nil {
		t.Fatal("expected error for missing file")
	}
}
