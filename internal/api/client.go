package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"domainshell/pkg/domain"
)

const defaultBaseURL = "https://edge.limoo.host/v1/domain"

type ClientInterface interface {
	CheckAvailability(domainName string) (*domain.Response, error)
	SuggestDomains(domainName string) (*domain.Response, error)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    defaultBaseURL,
	}
}

func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

func (c *Client) CheckAvailability(domainName string) (*domain.Response, error) {
	q := url.Values{}
	q.Add("domain[]", domainName)

	reqURL := fmt.Sprintf("%s/check-availability?%s", c.baseURL, q.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	var result domain.Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &result, nil
}

func (c *Client) SuggestDomains(domainName string) (*domain.Response, error) {
	reqURL := fmt.Sprintf("%s/suggest?domain=%s", c.baseURL, url.QueryEscape(domainName))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	var result domain.Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return &result, nil
}
