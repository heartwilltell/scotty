# scotty

🖖👨‍💻`Scotty` - Zero dependencies library to build simple commandline apps.

Basically it is a thin wrapper around standard `flag.FlagSet` type.

## Documentation

[![Go Report Card](https://goreportcard.com/badge/github.com/heartwilltell/scotty)](https://goreportcard.com/report/github.com/heartwilltell/scotty)
[![GoDoc](https://pkg.go.dev/badge/github.com/heartwilltell/scotty?utm_source=godoc)](https://pkg.go.dev/github.com/heartwilltell/scotty)
[![Build](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml/badge.svg?branch=main&event=push)](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml)
[![codecov](https://codecov.io/gh/heartwilltell/scotty/branch/main/graph/badge.svg?token=JFY9EQ4F2A)](https://codecov.io/gh/heartwilltell/scotty)

## Features

- 🤓 Simple API.
- 👌 Zero dependencies.
- 😘 Plays nice with standard `flag` package.
- 😌 Nice default `-help` information.
- 🏷️ Struct tag-based config binding with flag and env var support.

## Installation

```bash
go get github.com/heartwilltell/scotty
```

## Usage

The usage is pretty simple:

- 👉 Declare the root command.
- 👉 Attach subcommands and flags to it.
- 👉 Write your stuff inside the `Run` function.
- 👉 Call `Exec` function of the root command somewhere in `main` function.

```go
package main

import (
 "fmt"
 "log"
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
  Short: "Subcommands that does something",
  Run: func(cmd *scotty.Command, args []string) error {
   // Do some your stuff here.
   return nil
  },
 }

 // And here how you bind some flags to your command.
 var logLVL string

 subCmd.Flags().StringVar(&logLVL, "loglevel", "info",
  "set logging level: 'debug', 'info', 'warning', 'error'",
 )

 // Or use the SetFlags function.

 subCmd2 := scotty.Command{
  Name:  "subcommand2",
  Short: "Subcommands that does something",
  SetFlags: func(flags *FlagSet) {
   flags.StringVar(&logLVL, "loglevel", "info",
    "set logging level: 'debug', 'info', 'warning', 'error'",
   )
        },
  Run: func(cmd *scotty.Command, args []string) error {
   // Do some your stuff here.
   return nil
  },
 }

 // Attach subcommand to the root command. 
 rootCmd.AddSubcommands(&subCmd)

 // Execute the root command.
 if err := rootCmd.Exec(); err != nil {
  fmt.Println(err)
  os.Exit(2)
 }
}

```

## Config Binding

Scotty supports binding configuration structs to flags and environment variables using struct tags:

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
            fmt.Printf("Starting on %s:%d\n", cfg.Host, cfg.Port)
            return nil
        },
    }

    // Bind config before Exec
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
|-----|-------------|---------|
| `flag` | Flag name (required for binding) | `flag:"host"` |
| `env` | Environment variable name | `env:"APP_HOST"` |
| `default` | Default value | `default:"localhost"` |
| `usage` | Help text | `usage:"Server host"` |
| `required` | Must have non-zero value | `required:"true"` |

### Supported Types

`string`, `bool`, `int`, `int64`, `uint`, `uint64`, `float64`, `time.Duration`

### Generic Helpers

```go
// Type-safe config access
cfg := scotty.MustConfig[ServerConfig](cmd)

// Or with error handling
cfg, ok := scotty.GetConfig[ServerConfig](cmd)
```

## License

[MIT License](LICENSE).
