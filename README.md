# scotty

🖖👨‍💻`Scotty` - Zero dependencies library to build simple commandline apps.

Basically it is a thin wrapper around standard `flag.FlagSet` type.

## Documentation

[![](https://goreportcard.com/badge/github.com/heartwilltell/scotty)](https://goreportcard.com/report/github.com/heartwilltell/scotty)
[![](https://pkg.go.dev/badge/github.com/heartwilltell/scotty?utm_source=godoc)](https://pkg.go.dev/github.com/heartwilltell/scotty)
[![Build](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml/badge.svg?branch=master&event=push)](https://github.com/heartwilltell/scotty/actions/workflows/pr.yml)
[![codecov](https://codecov.io/gh/heartwilltell/scotty/branch/master/graph/badge.svg?token=JFY9EQ4F2A)](https://codecov.io/gh/heartwilltell/scotty)

## Features

- 🤓 Simple API.
- 👌 Zero dependencies.
- 😘 Plays nice with standard `flag` package.
- 😌 Nice default `-help` information.

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

	// Attach subcommand to the root command. 
	rootCmd.AddSubcommands(&subCmd)

	// Execute the root command.
	if err := rootCmd.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

```

## License

[MIT License](LICENSE).
