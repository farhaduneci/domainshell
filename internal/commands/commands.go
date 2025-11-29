package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"domainshell/internal/api"
)

type Commands struct {
	apiClient api.ClientInterface
}

func NewCommands(apiClient api.ClientInterface) *Commands {
	return &Commands{
		apiClient: apiClient,
	}
}

func formatPrice(price int) string {
	if price >= 1000000 {
		return fmt.Sprintf("%.2fM", float64(price)/1000000)
	} else if price >= 1000 {
		return fmt.Sprintf("%.1fK", float64(price)/1000)
	}
	return fmt.Sprintf("%d", price)
}

func (c *Commands) Search(domainName string) error {
	white := color.New(color.FgWhite)
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)

	result, err := c.apiClient.CheckAvailability(domainName)
	if err != nil {
		red.Printf("Request error: %v\n", err)
		return err
	}

	if len(result.Data) == 0 {
		yellow.Println("No data returned")
		return nil
	}

	item := result.Data[0]
	if item.Available {
		green.Printf("%s is available", item.Domain)
		if item.Prices.Register.OneYear > 0 {
			white.Printf(" (%s Toman/year)", formatPrice(item.Prices.Register.OneYear))
		}
		if item.OnSale {
			yellow.Print(" [ON SALE]")
		}
		if item.Premium {
			yellow.Print(" [PREMIUM]")
		}
		fmt.Println()
	} else {
		red.Printf("%s is NOT available", item.Domain)
		if item.Reason != "" {
			yellow.Printf(" (%s)", item.Reason)
		}
		fmt.Println()
	}

	return nil
}

func (c *Commands) Suggest(domainName string) error {
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	white := color.New(color.FgWhite)

	result, err := c.apiClient.SuggestDomains(domainName)
	if err != nil {
		red.Printf("Request error: %v\n", err)
		return err
	}

	if len(result.Data) == 0 {
		yellow.Println("No suggestions found")
		return nil
	}

	white.Printf("Suggestions for %s:\n", domainName)
	for _, item := range result.Data {
		if item.Available {
			green.Printf("  %s", item.Domain)
			if item.Prices.Register.OneYear > 0 {
				white.Printf(" (%s Toman/year)", formatPrice(item.Prices.Register.OneYear))
			}
			if item.OnSale {
				yellow.Print(" [ON SALE]")
			}
			if item.Premium {
				yellow.Print(" [PREMIUM]")
			}
			fmt.Println()
		}
	}

	return nil
}

func ParseInput(input string) (command string, args string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", ""
	}

	first := strings.ToLower(parts[0])
	knownCommands := map[string]bool{
		"search":  true,
		"suggest": true,
		"exit":    true,
		"quit":    true,
		"help":    true,
		"history": true,
	}

	if knownCommands[first] {
		if len(parts) > 1 {
			return first, strings.Join(parts[1:], " ")
		}
		return first, ""
	}

	return "search", input
}

func (c *Commands) Help() {
	white := color.New(color.FgWhite)
	cyan := color.New(color.FgCyan)

	cyan.Println("Available commands:")
	white.Println("  <domain>           - Check domain availability (default)")
	white.Println("  search <domain>    - Check domain availability")
	white.Println("  suggest <domain>   - Get domain suggestions")
	white.Println("  history            - Show command history")
	white.Println("  help               - Show this help message")
	white.Println("  exit, quit         - Exit the program")
	fmt.Println()
}
