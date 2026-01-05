# scotty

🖖👨‍💻`Scotty` - Zero dependencies library to build simple commandline apps.

Basically it is a thin wrapper around standard `flag.FlagSet` type.

## Documentation

[![Go Report Card](https://goreportcard.com/badge/github.com/heartwilltell/scotty)](https://goreportcard.com/report/github.com/heartwilltell/scotty)
[![GoDoc](https://pkg.go.dev/badge/github.com/heartwilltell/scotty?utm_source=godoc)](https://pkg.go.dev/github.com/heartwilltell/scotty)
[![Build](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml/badge.svg?branch=main&event=push)](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml)
[![codecov](https://codecov.io/gh/heartwilltell/scotty/branch/main/graph/badge.svg?token=JFY9EQ4F2A)](https://codecov.io/gh/heartwilltell/scotty)

- 🤓 Simple API.
- 👌 Zero dependencies.
- 😘 Plays nice with standard `flag` package.
- 😌 Nice default `-help` information.
- 🏷️ Struct tag-based config binding with flag and environment variable support.
- 🌍 Support for environment variables in flags (e.g., `StringVarE`, `BoolVarE`).

## Installation

```bash
go get github.com/heartwilltell/scotty
```

## Usage

The usage is pretty simple:

1. Declare the root command.
2. Attach subcommands and flags to it.
3. Write your logic inside the `Run` function.
4. Call `Exec` function of the root command.

```go
package main

import (
 "fmt"
 "os"

 "github.com/heartwilltell/scotty"
)

func main() {
 // Declare the root command. 
 rootCmd := scotty.Command{
  Name:  "app",
  Short: "Main command which holds all subcommands",
 }

 // Declare the subcommand.
 subCmd := scotty.Command{
  Name:  "subcommand",
  Short: "Subcommand that does something",
  Run: func(cmd *scotty.Command, args []string) error {
   fmt.Println("Running subcommand")
   return nil
  },
 }

 // Bind flags to your command.
 var logLVL string
 subCmd.Flags().StringVar(&logLVL, "loglevel", "info", "set logging level")

 // Or use the SetFlags function.
 subCmd2 := scotty.Command{
  Name:  "subcommand2",
  Short: "Another subcommand",
  SetFlags: func(flags *scotty.FlagSet) {
   flags.StringVar(&logLVL, "loglevel", "info", "set logging level")
  },
  Run: func(cmd *scotty.Command, args []string) error {
   fmt.Println("Running subcommand2")
   return nil
  },
 }

 // Attach subcommands to the root command. 
 rootCmd.AddSubcommands(&subCmd, &subCmd2)

 // Execute the root command.
 if err := rootCmd.Exec(); err != nil {
  fmt.Println(err)
  os.Exit(1)
 }
}
```

Scotty supports binding configuration structs to flags and environment variables using struct tags. You can use `BindConfig` on both `Command` and `FlagSet`.

```go
type ServerConfig struct {
    Host    string        `flag:"host" env:"HOST" default:"localhost" usage:"Server host"`
    Port    int           `flag:"port" env:"PORT" default:"8080" usage:"Server port" required:"true"`
    Debug   bool          `flag:"debug" env:"DEBUG" default:"false" usage:"Enable debug mode"`
    Timeout time.Duration `flag:"timeout" env:"TIMEOUT" default:"30s" usage:"Request timeout"`
}

// Optional: implement ConfigValidator for custom validation
func (c *ServerConfig) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    return nil
}

func main() {
    cfg := &ServerConfig{}

    cmd := &scotty.Command{
        Name:  "serve",
        Short: "Start the server",
        Run: func(cmd *scotty.Command, args []string) error {
            // Access config using the generic helper.
            cfg := scotty.MustConfig[ServerConfig](cmd)
            fmt.Printf("Starting on %s:%d\n", cfg.Host, cfg.Port)
            return nil
        },
    }

    // Bind config before Exec.
    if err := cmd.BindConfig(cfg); err != nil {
        log.Fatal(err)
    }

    if err := cmd.Exec(); err != nil {
        log.Fatal(err)
    }
}
```

### Supported Tags

| Tag | Description | Example |
| ----- | ------------- | --------- |
| `flag` | Flag name (required for binding) | `flag:"host"` |
| `env` | Environment variable name | `env:"APP_HOST"` |
| `default` | Default value | `default:"localhost"` |
| `usage` | Help text | `usage:"Server host"` |
| `required` | Must have non-zero value | `required:"true"` |

### Supported Types

`string`, `bool`, `int`, `int64`, `uint`, `uint64`, `float64`, `time.Duration`

### Generic Helpers

```go
// Type-safe config access (panics if not bound).
cfg := scotty.MustConfig[ServerConfig](cmd)

// Or with error handling.
cfg, ok := scotty.GetConfig[ServerConfig](cmd)
```

## Manual Environment Variable Binding

If you don't want to use struct tags, you can use the `*VarE` methods on `FlagSet` to bind flags with environment variable support. The flag value has priority over the environment variable.

Supported methods: `StringVarE`, `BoolVarE`, `IntVarE`, `Int64VarE`, `UintVarE`, `Uint64VarE`, `Float64VarE`, `DurationVarE`.

```go
var (
    apiPath string
    debug   bool
)

cmd := &scotty.Command{
    Name: "app",
    SetFlags: func(flags *scotty.FlagSet) {
        flags.StringVarE(&apiPath, "api-path", "API_PATH", "/v1", "API base path")
        flags.BoolVarE(&debug, "debug", "DEBUG", false, "Enable debug mode")
    },
}
```

## License

[MIT License](LICENSE).
