package scotty

import (
	"testing"
)

func TestCommand_IsSubcommand(t *testing.T) {
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
