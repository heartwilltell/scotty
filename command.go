package scotty

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	// flags holds set of commandline flags which are bind to this Command.
	// To avoid nil pointer exception it is better to work with flags via
	flags *flag.FlagSet
	// flagsState holds state of flags initialization.
	flagsState *sync.Once
	// subcommands holds set of Command who are a subcommand to this Command.
	subcommands map[string]*Command
	// parent holds a pointer to a parent Command.
	parent *Command
}

func (c *Command) Exec() error {
	// If the binary has been named differently that root command.
	if !c.IsSubcommand() {
		c.Name = filepath.Base(os.Args[0])
		flag.CommandLine = c.Flags()
	}

	// Parse all the program arguments.
	flag.Parse()

	return c.exec(os.Args[1:])
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

// TraverseToRoot traverse all command chain until it reaches root.
func (c *Command) TraverseToRoot() *Command {
	if c.IsSubcommand() {
		return c.parent.TraverseToRoot()
	}

	return c
}

func (c *Command) Flags() *flag.FlagSet {
	if c.flagsState == nil {
		c.flagsState = &sync.Once{}
	}

	c.flagsState.Do(func() {
		c.flags = flag.NewFlagSet(c.Name, flag.ExitOnError)
		c.flags.Usage = c.usage
	})

	return c.flags
}

func (c *Command) Args() []string { return c.Flags().Args() }

func (c *Command) exec(args []string) error {
	if err := c.Flags().Parse(args); err != nil {
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
		c.flags.Usage()

		return fmt.Errorf("%s: unknown command: %s", c.Name, args[0])
	}

	if c.Run == nil {
		c.Flags().Usage()
		return nil
	}

	if err := c.Run(c, args); err != nil {
		return fmt.Errorf("%s command failed: %w", c.Name, err)
	}

	return nil
}
