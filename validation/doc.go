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

// Package validation provides JSON Schema validation and capability negotiation.
//
// This package includes utilities for:
//
//   - Validating request/response payloads against UCP JSON schemas
//   - Capability negotiation between platforms and businesses
//   - Version compatibility checking
//   - Schema composition for extensions
//
// The validation logic ensures that all UCP messages conform to the
// official specification.
package validation
