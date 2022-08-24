package scotty

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Command represents a program command.
type Command struct {
	// Name represents command name and argument by which command will be called.
	Name string
	// Short represents short description of the command.
	Short string
	// Long represents short description of the command.
	Long string
	// Run represents a function which wraps and executes the logic of the command.
	Run func(cmd *Command, args []string) error
	// args holds a slice of commandline arguments.
	args []string
	// flags holds set of commandline flags which are bind to this Command.
	flags *flags
	// subcommands holds set of Command who are a subcommand to this Command.
	subcommands map[string]*Command
	// parent holds a pointer to a parent Command.
	parent *Command
}

func (c *Command) Exec() error {
	if !c.IsSubcommand() {
		c.Name = filepath.Base(os.Args[0])
		flag.CommandLine = c.Flags().flags
	}

	flag.Parse()

	return c.exec(os.Args[1:])
}

func (c *Command) Flags() *flags {
	if c.flags != nil {
		return c.flags
	}

	flagSet := flag.NewFlagSet(c.Name, flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Printf("usage: %s [flags] COMMAND\n", c.Name)
	}

	c.flags = &flags{
		name:  c.Name,
		flags: flagSet,
	}

	return c.flags
}

func (c *Command) AddSubcommands(commands ...*Command) {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command, len(commands))
	}

	for _, command := range commands {
		if command == c {
			panic(fmt.Errorf("command '%s' can't be a subcommand to itself", command.Name))
		}

		// Check if command has already been attached.
		if cmd, ok := c.subcommands[command.Name]; ok {
			if cmd != command {
				panic(fmt.Errorf("different command with a name '%s' already attached to '%s' command", command.Name, c.Name))
			}

			continue
		}

		// Attach the pointer to a parent to the subcommand.
		command.parent = c

		// Add command to the subcommand list of the parent command.
		c.subcommands[command.Name] = command
	}
}

// IsSubcommand return whether the command is subcommand for another command.
func (c *Command) IsSubcommand() bool {
	if c.parent == nil || c.parent == c {
		return false
	}

	return true
}

func (c *Command) exec(args []string) error {
	c.args = args

	if err := c.Flags().parse(args); err != nil {
		return fmt.Errorf("%s command failed: %w", c.Name, err)
	}

	// Check if the positional arguments is a subcommand.
	if len(args) > 0 && len(c.subcommands) > 0 {
		if subcommand, ok := c.subcommands[args[0]]; ok {
			// Subcommand has been found and should be executed.
			if err := subcommand.exec(args[1:]); err != nil {
				return err
			}

			return nil
		}

		// Looks line the argument is not it the list of known subcommands.
		// Let's print the usage and return an error.
		c.flags.flags.Usage()

		return fmt.Errorf("%s: unknown command: %s", c.Name, args[0])
	}

	if c.Run == nil {
		return nil
	}

	if err := c.Run(c, args); err != nil {
		return fmt.Errorf("%s command failed: %w", c.Name, err)
	}

	return nil
}
