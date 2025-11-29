package main

import (
	"fmt"
	"os"

	"domainshell/internal/api"
	"domainshell/internal/commands"
	"domainshell/internal/history"
	"domainshell/internal/repl"
)

func main() {
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
