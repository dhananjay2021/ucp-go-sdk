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

// Package internal contains internal utilities shared across the SDK.
package internal

import (
	"net/http"
	"time"
)

// DefaultHTTPClient returns a configured HTTP client with sensible defaults.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// RetryableClient wraps an HTTP client with retry logic.
type RetryableClient struct {
	client     *http.Client
	maxRetries int
	backoff    time.Duration
}

// NewRetryableClient creates a new retryable HTTP client.
func NewRetryableClient(client *http.Client, maxRetries int, backoff time.Duration) *RetryableClient {
	if client == nil {
		client = DefaultHTTPClient()
	}
	return &RetryableClient{
		client:     client,
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// Do executes an HTTP request with retry logic.
func (c *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	for i := 0; i <= c.maxRetries; i++ {
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(c.backoff * time.Duration(i+1))
			continue
		}

		// Retry on server errors
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = &ServerError{StatusCode: resp.StatusCode}
			time.Sleep(c.backoff * time.Duration(i+1))
			continue
		}

		return resp, nil
	}
	return nil, lastErr
}

// ServerError represents a server-side error.
type ServerError struct {
	StatusCode int
}

func (e *ServerError) Error() string {
	return http.StatusText(e.StatusCode)
}
