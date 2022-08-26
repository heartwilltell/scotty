package main

import (
	"fmt"
	"os"

	"github.com/heartwilltell/scotty"
)

func main() {
	root := scotty.Command{
		Name:  "scotty",
		Short: "This is a description of api",
	}

	var (
		a, b, c string
	)

	root.Flags().StringVar(&a, "flag-1", "", "usage of the flag")
	root.Flags().StringVar(&b, "flag-2", "", "usage of the flag")
	root.Flags().StringVar(&c, "flag-3", "", "usage of the flag")

	sub := scotty.Command{
		Name:  "a",
		Short: "sub shor description",
	}

	sub2 := scotty.Command{
		Name:  "b",
		Short: "sub shor description",
	}

	sub3 := scotty.Command{
		Name:  "c",
		Short: "sub shor description",
	}

	root.AddSubcommands(&sub, &sub2, &sub3)

	if err := root.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
