package history

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewEmptyHistory(t *testing.T) {
	h := NewEmptyHistory()
	if h == nil {
		t.Fatal("NewEmptyHistory() returned nil")
	}
	if h.filePath != "" {
		t.Errorf("Expected empty filePath, got %q", h.filePath)
	}
	if len(h.items) != 0 {
		t.Errorf("Expected empty items, got %d items", len(h.items))
	}
}

func TestHistory_Add(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "history.txt")

	h := &History{
		filePath: filePath,
		items:    make([]string, 0),
	}

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "add single item",
			input:    "example.com",
			expected: 1,
		},
		{
			name:     "add multiple items",
			input:    "example.org",
			expected: 2,
		},
		{
			name:     "add duplicate item",
			input:    "example.com",
			expected: 2,
		},
		{
			name:     "add empty string",
			input:    "",
			expected: 2,
		},
		{
			name:     "add whitespace only",
			input:    "   ",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Add(tt.input)
			if len(h.items) != tt.expected {
				t.Errorf("Expected %d items, got %d", tt.expected, len(h.items))
			}
		})
	}

	if h.items[len(h.items)-1] != "example.com" {
		t.Errorf("Expected last item to be 'example.com' (duplicate moved to end), got %q", h.items[len(h.items)-1])
	}
}

func TestHistory_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "history.txt")

	h1 := &History{
		filePath: filePath,
		items:    []string{"example.com", "example.org", "test.com"},
	}

	if err := h1.Save(); err != nil {
		t.Fatalf("Failed to save history: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("History file was not created")
	}

	h2 := &History{
		filePath: filePath,
		items:    make([]string, 0),
	}

	if err := h2.Load(); err != nil {
		t.Fatalf("Failed to load history: %v", err)
	}

	if len(h2.items) != len(h1.items) {
		t.Errorf("Expected %d items after load, got %d", len(h1.items), len(h2.items))
	}

	for i, item := range h1.items {
		if h2.items[i] != item {
			t.Errorf("Item %d mismatch: expected %q, got %q", i, item, h2.items[i])
		}
	}
}

func TestHistory_GetItems(t *testing.T) {
	h := &History{
		items: []string{"item1", "item2", "item3"},
	}

	items := h.GetItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	if items[0] != "item1" {
		t.Errorf("Expected first item to be 'item1', got %q", items[0])
	}
}

func TestHistory_GetDomains(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		expected []string
	}{
		{
			name:     "extract domains from commands",
			items:    []string{"search example.com", "suggest test.org", "example.net"},
			expected: []string{"example.com", "test.org", "example.net"},
		},
		{
			name:     "filter out commands",
			items:    []string{"search example.com", "help", "exit", "example.org"},
			expected: []string{"example.com", "example.org"},
		},
		{
			name:     "deduplicate domains",
			items:    []string{"example.com", "search example.com", "example.com"},
			expected: []string{"example.com"},
		},
		{
			name:     "no valid domains",
			items:    []string{"help", "exit", "quit"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &History{
				items: tt.items,
			}

			domains := h.GetDomains()

			if len(domains) != len(tt.expected) {
				t.Errorf("Expected %d domains, got %d", len(tt.expected), len(domains))
			}

			domainMap := make(map[string]bool)
			for _, d := range domains {
				domainMap[d] = true
			}

			for _, expected := range tt.expected {
				if !domainMap[expected] {
					t.Errorf("Expected domain %q not found in results", expected)
				}
			}
		})
	}
}

func TestHistory_GetHistoryFilePath(t *testing.T) {
	filePath := "/tmp/test/history.txt"
	h := &History{
		filePath: filePath,
	}

	if h.GetHistoryFilePath() != filePath {
		t.Errorf("Expected filePath %q, got %q", filePath, h.GetHistoryFilePath())
	}
}

func TestHistory_SaveWithEmptyPath(t *testing.T) {
	h := &History{
		filePath: "",
		items:    []string{"test"},
	}

	if err := h.Save(); err != nil {
		t.Errorf("Save() with empty path should not return error, got %v", err)
	}
}

func TestHistory_LoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.txt")

	h := &History{
		filePath: filePath,
		items:    make([]string, 0),
	}

	if err := h.Load(); err != nil {
		t.Errorf("Load() on non-existent file should not return error, got %v", err)
	}

	if len(h.items) != 0 {
		t.Errorf("Expected empty items after loading non-existent file, got %d items", len(h.items))
	}
}

func TestHistory_AddMaxItems(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "history.txt")

	h := &History{
		filePath: filePath,
		items:    make([]string, 0),
	}

	for i := 0; i < 1001; i++ {
		h.Add(fmt.Sprintf("example%d.com", i))
	}

	if len(h.items) > 1000 {
		t.Errorf("Expected max 1000 items, got %d", len(h.items))
	}

	if len(h.items) != 1000 {
		t.Errorf("Expected exactly 1000 items, got %d", len(h.items))
	}
}

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid domain",
			input:    "example.com",
			expected: true,
		},
		{
			name:     "command prefix",
			input:    "suggest",
			expected: false,
		},
		{
			name:     "search prefix",
			input:    "search",
			expected: false,
		},
		{
			name:     "exit command",
			input:    "exit",
			expected: false,
		},
		{
			name:     "quit command",
			input:    "quit",
			expected: false,
		},
		{
			name:     "no dot",
			input:    "example",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidDomain(tt.input)
			if result != tt.expected {
				t.Errorf("isValidDomain(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
