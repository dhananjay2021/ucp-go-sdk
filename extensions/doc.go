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

// Package extensions provides extended types that compose UCP base schemas.
//
// UCP uses a capability-based extension mechanism where base schemas can be
// extended with additional fields. This package provides Go types that
// combine base schemas with common extensions:
//
//   - ExtendedCheckoutResponse: Base checkout + fulfillment + discounts + buyer consent
//   - ExtendedCheckoutCreateRequest: Base request + extension fields
//   - ExtendedPaymentCredential: Base credential + token field
//
// These types simplify working with the full UCP feature set.
package extensions
