package main

import (
	"fmt"
	"os"

	"domainshell/internal/api"
	"domainshell/internal/commands"
	"domainshell/internal/history"
	"domainshell/internal/repl"
	"domainshell/internal/version"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("domainshell %s\n", version.Version)
		if version.BuildDate != "" {
			fmt.Printf("Build date: %s\n", version.BuildDate)
		}
		if version.GitCommit != "" {
			fmt.Printf("Git commit: %s\n", version.GitCommit)
		}
		os.Exit(0)
	}

	apiClient := api.NewClient()
	cmds := commands.NewCommands(apiClient)

	hist, err := history.NewHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize history: %v\n", err)
		hist = history.NewEmptyHistory()
	}

	r, err := repl.NewREPL(cmds, hist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
