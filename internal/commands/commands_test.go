package commands

import (
	"errors"
	"testing"

	"domainshell/pkg/domain"
)

type mockAPIClient struct {
	checkAvailabilityFunc func(string) (*domain.Response, error)
	suggestDomainsFunc    func(string) (*domain.Response, error)
}

func (m *mockAPIClient) CheckAvailability(domainName string) (*domain.Response, error) {
	if m.checkAvailabilityFunc != nil {
		return m.checkAvailabilityFunc(domainName)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAPIClient) SuggestDomains(domainName string) (*domain.Response, error) {
	if m.suggestDomainsFunc != nil {
		return m.suggestDomainsFunc(domainName)
	}
	return nil, errors.New("not implemented")
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedCmd  string
		expectedArgs string
	}{
		{
			name:         "empty input",
			input:        "",
			expectedCmd:  "",
			expectedArgs: "",
		},
		{
			name:         "whitespace only",
			input:        "   ",
			expectedCmd:  "",
			expectedArgs: "",
		},
		{
			name:         "domain only (default search)",
			input:        "example.com",
			expectedCmd:  "search",
			expectedArgs: "example.com",
		},
		{
			name:         "search command",
			input:        "search example.com",
			expectedCmd:  "search",
			expectedArgs: "example.com",
		},
		{
			name:         "suggest command",
			input:        "suggest example",
			expectedCmd:  "suggest",
			expectedArgs: "example",
		},
		{
			name:         "suggest with multiple words",
			input:        "suggest example domain",
			expectedCmd:  "suggest",
			expectedArgs: "example domain",
		},
		{
			name:         "help command",
			input:        "help",
			expectedCmd:  "help",
			expectedArgs: "",
		},
		{
			name:         "exit command",
			input:        "exit",
			expectedCmd:  "exit",
			expectedArgs: "",
		},
		{
			name:         "quit command",
			input:        "quit",
			expectedCmd:  "quit",
			expectedArgs: "",
		},
		{
			name:         "history command",
			input:        "history",
			expectedCmd:  "history",
			expectedArgs: "",
		},
		{
			name:         "case insensitive command",
			input:        "SEARCH example.com",
			expectedCmd:  "search",
			expectedArgs: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, args := ParseInput(tt.input)
			if cmd != tt.expectedCmd {
				t.Errorf("Expected command %q, got %q", tt.expectedCmd, cmd)
			}
			if args != tt.expectedArgs {
				t.Errorf("Expected args %q, got %q", tt.expectedArgs, args)
			}
		})
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		name     string
		price    int
		expected string
	}{
		{
			name:     "price less than 1000",
			price:    500,
			expected: "500",
		},
		{
			name:     "price in thousands",
			price:    5000,
			expected: "5.0K",
		},
		{
			name:     "price in millions",
			price:    1500000,
			expected: "1.50M",
		},
		{
			name:     "exact thousand",
			price:    1000,
			expected: "1.0K",
		},
		{
			name:     "exact million",
			price:    1000000,
			expected: "1.00M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPrice(tt.price)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCommands_Search(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		response    *domain.Response
		apiError    error
		expectError bool
	}{
		{
			name:   "available domain",
			domain: "example.com",
			response: &domain.Response{
				Data: []domain.DomainData{
					{
						Available: true,
						Domain:    "example.com",
						OnSale:    false,
						Premium:   false,
						Prices: struct {
							Register struct {
								OneYear int `json:"1y"`
							} `json:"register"`
						}{
							Register: struct {
								OneYear int `json:"1y"`
							}{
								OneYear: 100000,
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:   "unavailable domain",
			domain: "example.com",
			response: &domain.Response{
				Data: []domain.DomainData{
					{
						Available: false,
						Domain:    "example.com",
						Reason:    "Already registered",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "api error",
			domain:      "example.com",
			apiError:    errors.New("network error"),
			expectError: true,
		},
		{
			name:   "empty response",
			domain: "example.com",
			response: &domain.Response{
				Data: []domain.DomainData{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAPIClient{
				checkAvailabilityFunc: func(domainName string) (*domain.Response, error) {
					if tt.apiError != nil {
						return nil, tt.apiError
					}
					return tt.response, nil
				},
			}

			cmds := &Commands{
				apiClient: mockClient,
			}

			err := cmds.Search(tt.domain)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCommands_Suggest(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		response    *domain.Response
		apiError    error
		expectError bool
	}{
		{
			name:   "successful suggestions",
			domain: "example",
			response: &domain.Response{
				Data: []domain.DomainData{
					{Available: true, Domain: "example.com"},
					{Available: true, Domain: "example.org"},
					{Available: false, Domain: "example.net"},
				},
			},
			expectError: false,
		},
		{
			name:        "api error",
			domain:      "example",
			apiError:    errors.New("network error"),
			expectError: true,
		},
		{
			name:   "empty response",
			domain: "example",
			response: &domain.Response{
				Data: []domain.DomainData{},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockAPIClient{
				suggestDomainsFunc: func(domainName string) (*domain.Response, error) {
					if tt.apiError != nil {
						return nil, tt.apiError
					}
					return tt.response, nil
				},
			}

			cmds := &Commands{
				apiClient: mockClient,
			}

			err := cmds.Suggest(tt.domain)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCommands_Help(t *testing.T) {
	mockClient := &mockAPIClient{}
	cmds := NewCommands(mockClient)

	cmds.Help()
}
