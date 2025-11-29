package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"domainshell/pkg/domain"
)

func TestClient_CheckAvailability(t *testing.T) {
	tests := []struct {
		name           string
		domain         string
		response       domain.Response
		statusCode     int
		expectError    bool
		expectedDomain string
	}{
		{
			name:   "available domain",
			domain: "example.com",
			response: domain.Response{
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
						Reason: "",
					},
				},
			},
			statusCode:     http.StatusOK,
			expectError:    false,
			expectedDomain: "example.com",
		},
		{
			name:   "unavailable domain",
			domain: "example.com",
			response: domain.Response{
				Data: []domain.DomainData{
					{
						Available: false,
						Domain:    "example.com",
						Reason:    "Domain already registered",
					},
				},
			},
			statusCode:     http.StatusOK,
			expectError:    false,
			expectedDomain: "example.com",
		},
		{
			name:        "server error",
			domain:      "example.com",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/domain/check-availability" {
					t.Errorf("Expected path /v1/domain/check-availability, got %s", r.URL.Path)
				}

				domainParam := r.URL.Query().Get("domain[]")
				if domainParam != tt.domain {
					t.Errorf("Expected domain %s, got %s", tt.domain, domainParam)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.URL + "/v1/domain")
			client.httpClient = server.Client()

			result, err := client.CheckAvailability(tt.domain)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result.Data) == 0 {
				t.Error("Expected data in response")
				return
			}

			if result.Data[0].Domain != tt.expectedDomain {
				t.Errorf("Expected domain %s, got %s", tt.expectedDomain, result.Data[0].Domain)
			}
		})
	}
}

func TestClient_SuggestDomains(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		response    domain.Response
		statusCode  int
		expectError bool
		expectedLen int
	}{
		{
			name:   "successful suggestions",
			domain: "example",
			response: domain.Response{
				Data: []domain.DomainData{
					{Available: true, Domain: "example.com"},
					{Available: true, Domain: "example.org"},
					{Available: false, Domain: "example.net"},
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
			expectedLen: 3,
		},
		{
			name:        "server error",
			domain:      "example",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/domain/suggest" {
					t.Errorf("Expected path /v1/domain/suggest, got %s", r.URL.Path)
				}

				domainParam := r.URL.Query().Get("domain")
				if domainParam != tt.domain {
					t.Errorf("Expected domain %s, got %s", tt.domain, domainParam)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := NewClientWithBaseURL(server.URL + "/v1/domain")
			client.httpClient = server.Client()

			result, err := client.SuggestDomains(tt.domain)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result.Data) != tt.expectedLen {
				t.Errorf("Expected %d suggestions, got %d", tt.expectedLen, len(result.Data))
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Error("NewClient() httpClient is nil")
	}
}
