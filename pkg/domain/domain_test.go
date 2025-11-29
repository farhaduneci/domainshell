package domain

import (
	"encoding/json"
	"testing"
)

func TestDomainData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected DomainData
	}{
		{
			name: "full domain data",
			json: `{
				"available": true,
				"domain": "example.com",
				"on_sale": true,
				"premium": false,
				"prices": {
					"register": {
						"1y": 100000
					}
				},
				"reason": ""
			}`,
			expected: DomainData{
				Available: true,
				Domain:    "example.com",
				OnSale:    true,
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
		{
			name: "unavailable domain",
			json: `{
				"available": false,
				"domain": "example.com",
				"on_sale": false,
				"premium": false,
				"prices": {
					"register": {
						"1y": 0
					}
				},
				"reason": "Already registered"
			}`,
			expected: DomainData{
				Available: false,
				Domain:    "example.com",
				OnSale:    false,
				Premium:   false,
				Reason:    "Already registered",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data DomainData
			if err := json.Unmarshal([]byte(tt.json), &data); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if data.Available != tt.expected.Available {
				t.Errorf("Available: expected %v, got %v", tt.expected.Available, data.Available)
			}
			if data.Domain != tt.expected.Domain {
				t.Errorf("Domain: expected %q, got %q", tt.expected.Domain, data.Domain)
			}
			if data.OnSale != tt.expected.OnSale {
				t.Errorf("OnSale: expected %v, got %v", tt.expected.OnSale, data.OnSale)
			}
			if data.Premium != tt.expected.Premium {
				t.Errorf("Premium: expected %v, got %v", tt.expected.Premium, data.Premium)
			}
			if data.Prices.Register.OneYear != tt.expected.Prices.Register.OneYear {
				t.Errorf("Price: expected %d, got %d", tt.expected.Prices.Register.OneYear, data.Prices.Register.OneYear)
			}
			if data.Reason != tt.expected.Reason {
				t.Errorf("Reason: expected %q, got %q", tt.expected.Reason, data.Reason)
			}
		})
	}
}

func TestResponse_JSON(t *testing.T) {
	jsonStr := `{
		"data": [
			{
				"available": true,
				"domain": "example.com",
				"on_sale": false,
				"premium": false,
				"prices": {
					"register": {
						"1y": 100000
					}
				},
				"reason": ""
			},
			{
				"available": false,
				"domain": "example.org",
				"on_sale": false,
				"premium": false,
				"prices": {
					"register": {
						"1y": 0
					}
				},
				"reason": "Already registered"
			}
		]
	}`

	var response Response
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Expected 2 items in data, got %d", len(response.Data))
	}

	if response.Data[0].Domain != "example.com" {
		t.Errorf("Expected first domain to be 'example.com', got %q", response.Data[0].Domain)
	}

	if response.Data[1].Domain != "example.org" {
		t.Errorf("Expected second domain to be 'example.org', got %q", response.Data[1].Domain)
	}
}

func TestResponse_EmptyData(t *testing.T) {
	jsonStr := `{"data": []}`

	var response Response
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(response.Data) != 0 {
		t.Errorf("Expected empty data, got %d items", len(response.Data))
	}
}
