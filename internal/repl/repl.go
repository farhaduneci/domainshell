package repl

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"

	"domainshell/internal/commands"
	"domainshell/internal/history"
)

type REPL struct {
	cmds  *commands.Commands
	hist  *history.History
	rl    *readline.Instance
	white *color.Color
}

func NewREPL(cmds *commands.Commands, hist *history.History) (*REPL, error) {
	white := color.New(color.FgWhite)

	historyFile := ""
	if hist != nil {
		historyFile = hist.GetHistoryFilePath()
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "domain → ",
		HistoryFile:       historyFile,
		AutoComplete:      nil,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize readline: %w", err)
	}

	completer := &Completer{hist: hist}
	rl.Config.AutoComplete = completer

	r := &REPL{
		cmds:  cmds,
		hist:  hist,
		rl:    rl,
		white: white,
	}

	return r, nil
}

func (r *REPL) Run() error {
	defer r.rl.Close()

	for {
		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		r.hist.Add(line)

		command, args := commands.ParseInput(line)

		switch command {
		case "exit", "quit":
			return nil
		case "help":
			r.cmds.Help()
		case "history":
			r.showHistory()
		case "search":
			if args == "" {
				r.white.Println("Usage: search <domain>")
				continue
			}
			_ = r.cmds.Search(args)
		case "suggest":
			if args == "" {
				r.white.Println("Usage: suggest <domain>")
				continue
			}
			_ = r.cmds.Suggest(args)
		default:
			if args == "" && command != "" {
				_ = r.cmds.Search(command)
			}
		}
	}

	return nil
}

func (r *REPL) showHistory() {
	items := r.hist.GetItems()
	if len(items) == 0 {
		r.white.Println("No history")
		return
	}

	start := 0
	if len(items) > 20 {
		start = len(items) - 20
	}

	for _, item := range items[start:] {
		r.white.Println("  ", item)
	}
}

type Completer struct {
	hist *history.History
}

func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if c.hist == nil {
		return nil, 0
	}

	if pos == 0 {
		return nil, 0
	}

	start := pos
	for start > 0 && line[start-1] != ' ' && line[start-1] != '\t' {
		start--
	}

	prefix := string(line[start:pos])
	if prefix == "" {
		return nil, 0
	}

	prefixLower := strings.ToLower(prefix)
	text := string(line[:pos])
	parts := strings.Fields(text)

	var candidates []string
	commands := []string{"search", "suggest", "help", "history", "exit", "quit"}

	if len(parts) == 0 || (len(parts) == 1 && strings.HasPrefix(text, prefix)) {
		for _, cmd := range commands {
			if strings.HasPrefix(cmd, prefixLower) {
				candidates = append(candidates, cmd)
			}
		}
		if c.hist != nil {
			domains := c.hist.GetDomains()
			for _, domain := range domains {
				if strings.HasPrefix(strings.ToLower(domain), prefixLower) {
					candidates = append(candidates, domain)
				}
			}
		}
	} else if len(parts) > 1 {
		firstCmd := strings.ToLower(parts[0])
		if firstCmd == "search" || firstCmd == "suggest" {
			if c.hist != nil {
				domains := c.hist.GetDomains()
				for _, domain := range domains {
					if strings.HasPrefix(strings.ToLower(domain), prefixLower) {
						candidates = append(candidates, domain)
					}
				}
			}
		}
	}

	if len(candidates) == 0 {
		return nil, 0
	}

	prefixLen := pos - start
	result := make([][]rune, len(candidates))
	for i, cand := range candidates {
		if len(cand) > prefixLen {
			result[i] = []rune(cand[prefixLen:])
		} else {
			result[i] = []rune("")
		}
	}

	return result, prefixLen
}
