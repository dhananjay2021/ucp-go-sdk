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

// Package server provides HTTP handler utilities for implementing UCP endpoints.
//
// This package is intended for merchants and businesses that need to implement
// UCP-compliant APIs. It provides:
//
//   - HTTP handler helpers
//   - Request parsing and validation
//   - Response serialization with UCP metadata
//   - Middleware for authentication and capability negotiation
//   - Webhook signature generation and verification
//
// Example usage:
//
//	srv := server.NewServer(server.Config{Version: "2026-01-11"})
//	srv.HandleCreateCheckout(myHandler)
//	http.ListenAndServe(":8080", srv)
package server
