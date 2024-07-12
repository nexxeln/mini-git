package main

import (
	"fmt"
	"os"

	"github.com/nexxeln/mini-git/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mini-git <command> [<args>]")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		os.Exit(1)
	}

	switch command {
	case "init":
		if err := commands.Init(cwd); err != nil {
			fmt.Println("Error initializing repository:", err)
			os.Exit(1)
		}

	case "add":
		if len(args) < 1 {
			fmt.Println("Usage: mini-git add <file> [<file> ...]")
			fmt.Println("       mini-git add <directory> <file> [<file> ...] (any mix of directories and files in the current directory)")
			fmt.Println("       mini-git add . (to add all files in the current directory)")

			os.Exit(1)
		}

		if err := commands.Add(args); err != nil {
			fmt.Println("Error adding files:", err)
			os.Exit(1)
		}

	case "commit":
		if len(args) < 1 {
			fmt.Println("Usage: mini-git commit <message>")
			os.Exit(1)
		}
		author := "John Doe <john@example.com>" // TODO: Make this configurable
		if err := commands.Commit(cwd, args[0], author); err != nil {
			fmt.Println("Error committing changes:", err)
			os.Exit(1)
		}

	case "log":
		if err := commands.Log(cwd); err != nil {
			fmt.Println("Error displaying log:", err)
			os.Exit(1)
		}

	case "status":
		if err := commands.Status(cwd); err != nil {
			fmt.Println("Error displaying status:", err)
			os.Exit(1)
		}

	case "branch":
		if err := commands.Branch(cwd, args); err != nil {
			fmt.Println("Error handling branch command:", err)
			os.Exit(1)
		}

	case "checkout":
		if err := commands.Checkout(cwd, args); err != nil {
			fmt.Println("Error handling checkout command:", err)
			os.Exit(1)
		}

	case "merge":
		if err := commands.Merge(cwd, args); err != nil {
			fmt.Println("Error handling merge command:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
