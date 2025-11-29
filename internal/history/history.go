package history

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type History struct {
	filePath string
	items    []string
}

func NewHistory() (*History, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	historyDir := filepath.Join(homeDir, ".config", "domainshell")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	filePath := filepath.Join(historyDir, "history.txt")

	h := &History{
		filePath: filePath,
		items:    make([]string, 0),
	}

	if err := h.Load(); err != nil {
		return h, nil
	}

	return h, nil
}

func NewEmptyHistory() *History {
	return &History{
		filePath: "",
		items:    make([]string, 0),
	}
}

func (h *History) Load() error {
	file, err := os.Open(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			h.items = append(h.items, line)
		}
	}

	return scanner.Err()
}

func (h *History) Save() error {
	if h.filePath == "" {
		return nil
	}

	file, err := os.Create(h.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range h.items {
		if _, err := writer.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (h *History) Add(item string) {
	item = strings.TrimSpace(item)
	if item == "" {
		return
	}

	for i, existing := range h.items {
		if existing == item {
			h.items = append(h.items[:i], h.items[i+1:]...)
			break
		}
	}

	h.items = append(h.items, item)

	if len(h.items) > 1000 {
		h.items = h.items[len(h.items)-1000:]
	}

	_ = h.Save()
}

func (h *History) GetItems() []string {
	return h.items
}

func (h *History) GetDomains() []string {
	domains := make([]string, 0)
	seen := make(map[string]bool)

	for _, item := range h.items {
		parts := strings.Fields(item)
		if len(parts) > 0 {
			domain := parts[len(parts)-1]
			if !seen[domain] && isValidDomain(domain) {
				domains = append(domains, domain)
				seen[domain] = true
			}
		}
	}

	return domains
}

func (h *History) GetHistoryFilePath() string {
	return h.filePath
}

func isValidDomain(s string) bool {
	return strings.Contains(s, ".") && !strings.HasPrefix(s, "suggest") && !strings.HasPrefix(s, "search") && s != "exit" && s != "quit"
}
