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

package client

import "net/http"

// RequestOption is a function that modifies an HTTP request.
type RequestOption func(*http.Request)

// WithIdempotencyKey adds an idempotency key header.
func WithIdempotencyKey(key string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Idempotency-Key", key)
	}
}

// WithRequestID adds a request ID header for tracing.
func WithRequestID(id string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("X-Request-ID", id)
	}
}

// WithAcceptLanguage sets the Accept-Language header.
func WithAcceptLanguage(lang string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Accept-Language", lang)
	}
}

// WithHeader sets a custom header.
func WithHeader(key, value string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
