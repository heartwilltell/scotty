package scotty

import (
	"flag"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

func (c *Command) usage() {
	// Define the single strings.Builder
	// for the output of the command usage.
	var b strings.Builder

	root := c.TraverseToRoot()
	if root.Short != "" {
		b.WriteString(fmt.Sprintf("%s - %s\n\n", root.Name, root.Short))
	} else {
		b.WriteString(fmt.Sprintf("%s\n\n", root.Name))
	}

	b.WriteString("Usage:\n")
	printCommandCallUsage(&b, c)

	printSubcommands(&b, c.subcommands)
	printFlags(&b, c.Flags())
	printHelpSuggestion(&b, c)

	if _, err := fmt.Fprintln(c.Flags().Output(), b.String()); err != nil {
		fmt.Println(b.String())
	}
}

func printCommandCallUsage(b *strings.Builder, c *Command) {
	b.WriteString(fmt.Sprintf("  %s ", commandsChain(c)))

	if hasFlags(c.Flags()) {
		b.WriteString("<flags> ")
	}

	if len(c.subcommands) > 0 {
		b.WriteString("[command]\n")
	} else {
		b.WriteString("[arguments...]\n")
	}
}

func printSubcommands(b *strings.Builder, subcommands map[string]*Command) {
	if b == nil || len(subcommands) == 0 {
		return
	}

	b.WriteString("\nAvailable Commands:\n")

	sorted := sortedSubcommands(subcommands)
	longest := 0

	for _, c := range sorted {
		nameLen := utf8.RuneCountInString(c.Name)

		if longest < nameLen {
			longest = nameLen
		}
	}

	for _, c := range sorted {
		b.WriteString(fmt.Sprintf("  %s %s\n", c.Name+indent(c.Name, longest, 1), c.Short))
	}
}

func printFlags(b *strings.Builder, flags *flag.FlagSet) {
	if b == nil || flags == nil {
		return
	}

	flagCount := 0
	longest := 0

	flags.VisitAll(func(f *flag.Flag) {
		flagCount++

		fType := reflect.TypeOf(f.Value).Elem().Kind().String()
		nameLen := utf8.RuneCountInString(f.Name + " " + fType)

		if longest < nameLen {
			longest = nameLen
		}
	})

	if flagCount > 0 {
		b.WriteString("\nFlags:\n")
		flags.VisitAll(func(f *flag.Flag) {
			fType := reflect.TypeOf(f.Value).Elem().Kind().String()

			b.WriteString(fmt.Sprintf("  -%s %s\n",
				f.Name+" "+fType+indent(f.Name+" "+fType, longest, 1),
				f.Usage,
			))
		})
	}
}

func printHelpSuggestion(b *strings.Builder, c *Command) {
	b.WriteString(fmt.Sprintf("\nUse '%s -help' for more information about a command.\n", commandsChain(c)))
}

func sortedSubcommands(subcommands map[string]*Command) []*Command {
	commands := make([]*Command, 0, len(subcommands))

	for _, cmd := range subcommands {
		commands = append(commands, cmd)
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

func commandsChain(c *Command) string {
	commands := make([]string, 0, 1)

	walkCommandsChain(c, func(name string) {
		commands = append(commands, name)
	})

	// Reversing tha slice to make the correct order.
	for i, j := 0, len(commands)-1; i < j; i, j = i+1, j-1 {
		commands[i], commands[j] = commands[j], commands[i]
	}

	var chain string

	for i, cmd := range commands {
		chain += cmd
		if i != len(commands)-1 {
			chain += " "
		}
	}

	return chain
}

func walkCommandsChain(c *Command, f func(name string)) {
	f(c.Name)

	if c.parent == nil {
		return
	}

	walkCommandsChain(c.parent, f)
}

func indent(name string, longest, offset int) string {
	whitespace := ""
	for i := utf8.RuneCountInString(name); i < longest+offset; i++ {
		whitespace += " "
	}

	return whitespace
}

func hasFlags(flags *flag.FlagSet) bool {
	has := false

	flags.VisitAll(func(_ *flag.Flag) { has = true })

	return has
}
