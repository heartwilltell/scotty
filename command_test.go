package scotty

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestCommand_IsSubcommand(t *testing.T) {
	helperDisableStdout(t)

	t.Run("IsSubcommand false", func(t *testing.T) {
		c := Command{Name: "root"}
		got := c.IsSubcommand()

		if got != false {
			t.Errorf("IsSubcommand(): unexpected result: '%s' command does not have parent command", c.Name)
		}
	})

	t.Run("IsSubcommand true", func(t *testing.T) {
		c := Command{Name: "root"}
		c2 := Command{Name: "c2"}
		c.AddSubcommands(&c2)

		got := c2.IsSubcommand()

		if got != true {
			t.Errorf("IsSubcommand(): unexpected result: '%s' command should have parent command", c2.Name)
		}
	})
}

func TestCommand_TraverseToRoot(t *testing.T) {
	helperDisableStdout(t)

	t.Run("TraverseToRoot traverse", func(t *testing.T) {
		c3 := Command{Name: "c3"}
		c2 := Command{Name: "c2"}
		c1 := Command{Name: "c1"}
		root := Command{Name: "root"}

		root.AddSubcommands(&c1)
		c1.AddSubcommands(&c2)
		c2.AddSubcommands(&c3)

		if c3.TraverseToRoot() != &root {
			t.Errorf("TraverseToRoot(): unexpected pointer to root command: expected := %p, got := %p", &root, c3.TraverseToRoot())
		}

		if c2.TraverseToRoot() != &root {
			t.Errorf("TraverseToRoot(): unexpected pointer to root command: expected := %p, got := %p", &root, c2.TraverseToRoot())
		}

		if c1.TraverseToRoot() != &root {
			t.Errorf("TraverseToRoot(): unexpected pointer to root command: expected := %p, got := %p", &root, c1.TraverseToRoot())
		}
	})
}

func TestCommand_Exec(t *testing.T) {
	helperDisableStdout(t)

	t.Run("OK", func(t *testing.T) {
		cmd := &Command{
			Name: "test",
			Run: func(cmd *Command, args []string) error {
				return nil
			},
		}

		// Set the testing flags to the command so the
		// command will not fail with the error like:
		// "flag provided but not defined: -test.timeout".
		helperSetTestingFlags(t, cmd)

		got := cmd.Exec()
		if !reflect.DeepEqual(got, nil) {
			t.Errorf("Expected := nil, got := %#v", got)
		}
	})
}

func TestCommand_AddSubcommands(t *testing.T) {
	helperDisableStdout(t)

	type tcase struct {
		cmd      *Command
		subCmd   *Command
		panicVal error
	}

	// Additional values for some test cases.
	testCmd1 := &Command{Name: "panic-add-self"}
	testCmd2 := &Command{Name: "panic-already-attached"}
	testCmd3 := &Command{Name: "panic-already-attached_attached"}
	testCmd2.AddSubcommands(testCmd3)

	tests := map[string]tcase{
		"OK": {
			cmd:      &Command{Name: "ok"},
			subCmd:   &Command{Name: "ok-sub"},
			panicVal: nil,
		},

		"Panic add self": {
			cmd:      testCmd1,
			subCmd:   testCmd1,
			panicVal: fmt.Errorf("command '%s' can't be a subcommand to itself", testCmd1.Name),
		},

		"Already attached": {
			cmd:      testCmd2,
			subCmd:   testCmd3,
			panicVal: nil,
		},

		"Panic same name attached": {
			cmd:      testCmd2,
			subCmd:   &Command{Name: "panic-already-attached_attached"},
			panicVal: fmt.Errorf("different command with a name '%s' already attached to '%s' command", testCmd3.Name, testCmd2.Name),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer helperCatchPanic(t, tc.panicVal)

			tc.cmd.AddSubcommands(tc.subCmd)

			subCmd, ok := tc.cmd.subcommands[tc.subCmd.Name]
			if !ok {
				t.Errorf("Expected subcommand := %+v, got := nil", tc.subCmd)
			}

			if !reflect.DeepEqual(subCmd, tc.subCmd) {
				t.Errorf("Expected subcommand := %+v, got := %+v", tc.subCmd, subCmd)
			}
		})
	}
}

// DO NOT RUN MANUALLY from GoLand, or VSCode by 'play' button.
func TestCommand_Args(t *testing.T) {
	helperDisableStdout(t)

	t.Run("OK", func(t *testing.T) {
		cmd := &Command{Name: "test"}

		// Set the testing flags to the command so the
		// command will not fail with the error like:
		// "flag provided but not defined: -test.timeout".
		helperSetTestingFlags(t, cmd)

		want := make([]string, 0)

		if err := cmd.Exec(); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		got := cmd.Args()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected := %+v, got := %+v", want, got)
		}

	})
}

func TestCommand_exec(t *testing.T) {
	helperDisableStdout(t)

	type tcase struct {
		cmd     *Command
		args    []string
		wantErr error
	}

	tests := map[string]tcase{
		"OK": {
			cmd:     &Command{Name: "test"},
			args:    nil,
			wantErr: nil,
		},

		"Subcommand execCommand": {
			cmd: func() *Command {
				cmd := &Command{Name: "test"}
				cmd.AddSubcommands(&Command{
					Name: "sub",
					Run: func(cmd *Command, args []string) error {
						return errors.New("sub error")
					},
				})
				return cmd
			}(),
			args:    []string{"sub"},
			wantErr: fmt.Errorf("command failed: %w", errors.New("sub error")),
		},

		"Unknown subcommand": {
			cmd: func() *Command {
				cmd := &Command{Name: "test"}
				cmd.AddSubcommands(&Command{
					Name: "sub",
					Run:  func(cmd *Command, args []string) error { return nil },
				})
				return cmd
			}(),
			args:    []string{"base"},
			wantErr: fmt.Errorf("unknown command: %s", "base"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.cmd.execCommand(tc.args)

			if !reflect.DeepEqual(got, tc.wantErr) {
				t.Errorf("Expected := %#v, got := %#v", tc.wantErr, got)
			}
		})
	}
}

func TestCommand_SetPersistentFlags(t *testing.T) {
	helperDisableStdout(t)

	t.Run("Inherited by subcommand", func(t *testing.T) {
		var verbose bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		sub := &Command{
			Name: "sub",
			Run:  func(cmd *Command, args []string) error { return nil },
		}

		root.AddSubcommands(sub)

		got := sub.execCommand([]string{"-verbose"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}
	})

	t.Run("Inherited by grandchild", func(t *testing.T) {
		var verbose bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		mid := &Command{Name: "mid"}
		leaf := &Command{
			Name: "leaf",
			Run:  func(cmd *Command, args []string) error { return nil },
		}

		root.AddSubcommands(mid)
		mid.AddSubcommands(leaf)

		got := leaf.execCommand([]string{"-verbose"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}
	})

	t.Run("Works with local flags", func(t *testing.T) {
		var verbose bool
		var port int

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		sub := &Command{
			Name: "sub",
			SetFlags: func(flags *FlagSet) {
				flags.IntVar(&port, "port", 8080, "server port")
			},
			Run: func(cmd *Command, args []string) error { return nil },
		}

		root.AddSubcommands(sub)

		got := sub.execCommand([]string{"-verbose", "-port", "9090"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}

		if port != 9090 {
			t.Errorf("Expected port to be 9090, got %d", port)
		}
	})

	t.Run("Parsed at parent level", func(t *testing.T) {
		var verbose bool
		var called bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		sub := &Command{
			Name: "sub",
			Run: func(cmd *Command, args []string) error {
				called = true
				return nil
			},
		}

		root.AddSubcommands(sub)

		// Simulate: args after root flag parsing where -verbose was
		// already consumed by the parent. Subcommand receives remaining args.
		got := root.execCommand([]string{"sub"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !called {
			t.Error("Expected subcommand Run to be called")
		}
	})

	t.Run("Multi-level persistent flags", func(t *testing.T) {
		var verbose bool
		var format string

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		mid := &Command{
			Name: "mid",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.StringVar(&format, "format", "text", "output format")
			},
		}

		leaf := &Command{
			Name: "leaf",
			Run:  func(cmd *Command, args []string) error { return nil },
		}

		root.AddSubcommands(mid)
		mid.AddSubcommands(leaf)

		got := leaf.execCommand([]string{"-verbose", "-format", "json"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}

		if format != "json" {
			t.Errorf("Expected format to be 'json', got '%s'", format)
		}
	})

	t.Run("Persistent flag between subcommands", func(t *testing.T) {
		var verbose bool
		var called bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		sub := &Command{
			Name: "sub",
			Run: func(cmd *Command, args []string) error {
				called = true
				return nil
			},
		}

		root.AddSubcommands(sub)

		got := root.execCommand([]string{"-verbose", "sub"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}

		if !called {
			t.Error("Expected subcommand Run to be called")
		}
	})

	t.Run("Persistent flag before subcommand in chain", func(t *testing.T) {
		var verbose bool
		var called bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		mid := &Command{Name: "mid"}
		leaf := &Command{
			Name: "leaf",
			Run: func(cmd *Command, args []string) error {
				called = true
				return nil
			},
		}

		root.AddSubcommands(mid)
		mid.AddSubcommands(leaf)

		got := root.execCommand([]string{"mid", "-verbose", "leaf"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}

		if !called {
			t.Error("Expected leaf Run to be called")
		}
	})

	t.Run("Multi-level persistent flags between subcommands", func(t *testing.T) {
		var verbose bool
		var format string
		var called bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
		}

		mid := &Command{
			Name: "mid",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.StringVar(&format, "format", "text", "output format")
			},
		}

		leaf := &Command{
			Name: "leaf",
			Run: func(cmd *Command, args []string) error {
				called = true
				return nil
			},
		}

		root.AddSubcommands(mid)
		mid.AddSubcommands(leaf)

		got := root.execCommand([]string{"mid", "-verbose", "-format", "json", "leaf"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}

		if format != "json" {
			t.Errorf("Expected format to be 'json', got '%s'", format)
		}

		if !called {
			t.Error("Expected leaf Run to be called")
		}
	})

	t.Run("Available on defining command itself", func(t *testing.T) {
		var verbose bool

		root := &Command{
			Name: "root",
			SetPersistentFlags: func(flags *FlagSet) {
				flags.BoolVar(&verbose, "verbose", false, "verbose output")
			},
			Run: func(cmd *Command, args []string) error { return nil },
		}

		got := root.execCommand([]string{"-verbose"})
		if got != nil {
			t.Fatalf("Unexpected error: %v", got)
		}

		if !verbose {
			t.Error("Expected verbose to be true, got false")
		}
	})
}

func helperDisableStdout(t *testing.T) {
	tmpStdout := os.Stdout
	tmpStderr := os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr, _ = os.Open(os.DevNull)

	t.Cleanup(func() {
		os.Stdout = tmpStdout
		os.Stderr = tmpStderr
	})
}

func helperCatchPanic(t *testing.T, expected error) {
	t.Helper()
	r := recover()
	if !reflect.DeepEqual(r, expected) {
		t.Errorf("Expected recover value := %+v, got := %#v", expected, r)
	}
}

func helperSetTestingFlags(t *testing.T, cmd *Command) {
	t.Helper()

	cmd.Flags().Bool("test.v", true, "")
	cmd.Flags().Bool("test.run", true, "")
	cmd.Flags().Bool("test.paniconexit0", true, "")
	cmd.Flags().String("test.testlogfile", "", "")
	cmd.Flags().String("test.timeout", "", "")
	cmd.Flags().String("test.coverprofile", "", "")
	cmd.Flags().String("test.gocoverdir", "", "")
}
