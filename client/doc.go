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
//
// This package is intended for platforms and agents that need to interact
// with UCP-compliant merchants. It provides typed methods for all UCP
// operations including:
//
//   - Profile discovery
//   - Checkout session management
//   - Order retrieval
//   - Identity linking
//
// Example usage:
//
//	client := client.NewClient("https://merchant.example.com")
//	profile, err := client.FetchProfile(ctx)
//	checkout, err := client.CreateCheckout(ctx, req)
package client
