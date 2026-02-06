// Copyright 2026 UCP Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client provides a REST client for consuming UCP APIs.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dhananjay2021/ucp-go-sdk/extensions"
	"github.com/dhananjay2021/ucp-go-sdk/models"
)

const (
	// DefaultTimeout is the default HTTP request timeout.
	DefaultTimeout = 30 * time.Second

	// WellKnownPath is the discovery profile path.
	WellKnownPath = "/.well-known/ucp"

	// CheckoutSessionsPath is the checkout sessions endpoint.
	CheckoutSessionsPath = "/checkout-sessions"

	// OrdersPath is the orders endpoint.
	OrdersPath = "/orders"

	// CartsPath is the shopping carts endpoint.
	CartsPath = "/carts"
)

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithAccessToken sets the OAuth access token for authentication.
func WithAccessToken(token string) ClientOption {
	return func(c *Client) {
		c.accessToken = token
	}
}

// WithUserAgent sets the User-Agent header.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithUCPAgent sets the UCP-Agent header with the platform's profile URL.
// This header is required on all UCP requests and identifies the calling platform.
// Format: profile="https://platform.example/.well-known/ucp"
func WithUCPAgent(profileURL string) ClientOption {
	return func(c *Client) {
		c.ucpAgentProfile = profileURL
	}
}

// Client is a UCP REST API client.
type Client struct {
	baseURL         string
	httpClient      *http.Client
	timeout         time.Duration
	apiKey          string
	accessToken     string
	userAgent       string
	ucpAgentProfile string

	// Cached discovery profile
	profile *models.UCPProfile
}

// NewClient creates a new UCP client.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:   baseURL,
		timeout:   DefaultTimeout,
		userAgent: "ucp-go-sdk/1.0",
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.timeout,
		}
	}

	return c
}

// Error represents an API error response.
type Error struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("UCP API error (status %d): %s", e.StatusCode, e.Message)
}

// doRequest performs an HTTP request and decodes the response.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path

	// Encode body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	if c.ucpAgentProfile != "" {
		req.Header.Set("UCP-Agent", fmt.Sprintf(`profile="%s"`, c.ucpAgentProfile))
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		apiErr := &Error{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
		if len(respBody) > 0 {
			var errDetails map[string]interface{}
			if json.Unmarshal(respBody, &errDetails) == nil {
				apiErr.Details = errDetails
				if msg, ok := errDetails["message"].(string); ok {
					apiErr.Message = msg
				}
			}
		}
		return apiErr
	}

	// Decode response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// FetchProfile fetches the discovery profile from /.well-known/ucp.
func (c *Client) FetchProfile(ctx context.Context) (*models.UCPProfile, error) {
	var profile models.UCPProfile
	if err := c.doRequest(ctx, http.MethodGet, WellKnownPath, nil, &profile); err != nil {
		return nil, err
	}
	c.profile = &profile
	return &profile, nil
}

// GetCachedProfile returns the cached discovery profile, fetching it if necessary.
func (c *Client) GetCachedProfile(ctx context.Context) (*models.UCPProfile, error) {
	if c.profile != nil {
		return c.profile, nil
	}
	return c.FetchProfile(ctx)
}

// CreateCheckout creates a new checkout session.
func (c *Client) CreateCheckout(ctx context.Context, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error) {
	var resp extensions.ExtendedCheckoutResponse
	if err := c.doRequest(ctx, http.MethodPost, CheckoutSessionsPath, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCheckout retrieves a checkout session by ID.
func (c *Client) GetCheckout(ctx context.Context, id string) (*extensions.ExtendedCheckoutResponse, error) {
	var resp extensions.ExtendedCheckoutResponse
	path := fmt.Sprintf("%s/%s", CheckoutSessionsPath, id)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCheckout updates a checkout session.
func (c *Client) UpdateCheckout(ctx context.Context, id string, req *extensions.ExtendedCheckoutUpdateRequest) (*extensions.ExtendedCheckoutResponse, error) {
	var resp extensions.ExtendedCheckoutResponse
	path := fmt.Sprintf("%s/%s", CheckoutSessionsPath, id)
	if err := c.doRequest(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CompleteCheckout completes a checkout session.
func (c *Client) CompleteCheckout(ctx context.Context, id string) (*extensions.ExtendedCheckoutResponse, error) {
	var resp extensions.ExtendedCheckoutResponse
	path := fmt.Sprintf("%s/%s/complete", CheckoutSessionsPath, id)
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelCheckout cancels a checkout session.
func (c *Client) CancelCheckout(ctx context.Context, id string) (*extensions.ExtendedCheckoutResponse, error) {
	var resp extensions.ExtendedCheckoutResponse
	path := fmt.Sprintf("%s/%s/cancel", CheckoutSessionsPath, id)
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrder retrieves an order by ID.
func (c *Client) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	var resp models.Order
	path := fmt.Sprintf("%s/%s", OrdersPath, id)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCart creates a new shopping cart.
// Carts provide lightweight pre-purchase exploration with estimated pricing
// before committing to a checkout session.
func (c *Client) CreateCart(ctx context.Context, req *models.CartCreateRequest) (*models.CartResponse, error) {
	var resp models.CartResponse
	if err := c.doRequest(ctx, http.MethodPost, CartsPath, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCart retrieves a cart by ID.
func (c *Client) GetCart(ctx context.Context, id string) (*models.CartResponse, error) {
	var resp models.CartResponse
	path := fmt.Sprintf("%s/%s", CartsPath, id)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCart updates a cart with new items.
// Line items are fully replaced on update.
func (c *Client) UpdateCart(ctx context.Context, id string, req *models.CartUpdateRequest) (*models.CartResponse, error) {
	var resp models.CartResponse
	path := fmt.Sprintf("%s/%s", CartsPath, id)
	if err := c.doRequest(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCart deletes a cart.
func (c *Client) DeleteCart(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s", CartsPath, id)
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return err
	}
	return nil
}

// CreateCheckoutFromCart creates a checkout session from an existing cart.
// This converts the cart to a checkout, using the cart's line_items, context, and buyer.
func (c *Client) CreateCheckoutFromCart(ctx context.Context, cartID string, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error) {
	// Create a wrapper that includes cart_id
	type checkoutWithCart struct {
		*extensions.ExtendedCheckoutCreateRequest
		CartID string `json:"cart_id"`
	}
	wrapped := &checkoutWithCart{
		ExtendedCheckoutCreateRequest: req,
		CartID:                        cartID,
	}

	var resp extensions.ExtendedCheckoutResponse
	if err := c.doRequest(ctx, http.MethodPost, CheckoutSessionsPath, wrapped, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
