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

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuth2Config contains OAuth2 configuration for identity linking.
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURL  string
	Scopes       []string
}

// OAuth2Token represents an OAuth2 access token.
type OAuth2Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

// IsExpired checks if the token has expired.
func (t *OAuth2Token) IsExpired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(t.ExpiresAt.Add(-time.Minute)) // 1 minute buffer
}

// OAuth2Client handles OAuth2 token operations.
type OAuth2Client struct {
	config     OAuth2Config
	httpClient *http.Client
}

// NewOAuth2Client creates a new OAuth2 client.
func NewOAuth2Client(config OAuth2Config) *OAuth2Client {
	return &OAuth2Client{
		config:     config,
		httpClient: DefaultHTTPClient(),
	}
}

// AuthCodeURL returns the URL for the authorization code flow.
func (c *OAuth2Client) AuthCodeURL(state string) string {
	u, _ := url.Parse(c.config.AuthURL)
	q := u.Query()
	q.Set("client_id", c.config.ClientID)
	q.Set("redirect_uri", c.config.RedirectURL)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(c.config.Scopes, " "))
	q.Set("state", state)
	u.RawQuery = q.Encode()
	return u.String()
}

// ExchangeCode exchanges an authorization code for tokens.
func (c *OAuth2Client) ExchangeCode(ctx context.Context, code string) (*OAuth2Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.config.RedirectURL)
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)

	return c.tokenRequest(ctx, data)
}

// RefreshToken refreshes an access token using a refresh token.
func (c *OAuth2Client) RefreshToken(ctx context.Context, refreshToken string) (*OAuth2Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)

	return c.tokenRequest(ctx, data)
}

func (c *OAuth2Client) tokenRequest(ctx context.Context, data url.Values) (*OAuth2Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	token := &OAuth2Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
	}

	if tokenResp.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return token, nil
}
