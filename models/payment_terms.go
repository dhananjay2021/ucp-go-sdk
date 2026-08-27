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

package models

import "time"

// PaymentTerm is a way of paying for the checkout: one or more payment schedules
// that together cover its total. The Payment Terms extension (UCP spec change
// #602) lets a business offer alternative schedules for when payment for the
// checkout is due, and projects the accepted term onto the resulting order via
// payment.selected_term_id.
type PaymentTerm struct {
	// ID is the unique identifier for this payment term within the checkout.
	// Referenced by payment.selected_term_id.
	ID string `json:"id"`

	// Title is a short label distinguishing this term from its siblings (e.g.
	// "Pay now", "Pay in 4", "Deposit + balance at check-in").
	Title string `json:"title"`

	// Description is supplementary context for the title (e.g. "Save 5% by paying
	// today"). Directly renderable; MUST NOT repeat the title.
	Description *Description `json:"description,omitempty"`

	// Schedules are the payment schedules that settle this checkout under this
	// term, in the order they come due.
	Schedules []PaymentSchedule `json:"schedules"`
}

// PaymentScheduleType is the timing class for a payment schedule, drawn from an
// open vocabulary.
type PaymentScheduleType string

const (
	// PaymentScheduleTypeImmediate is the only value with defined meaning: the
	// payment is due when the checkout is completed. Any other value means the
	// payment is not due at completion and Description states when it is due.
	PaymentScheduleTypeImmediate PaymentScheduleType = "immediate"
)

// PaymentSchedule is a single payment that settles part or all of the checkout
// under a payment term. Timing is stated in buyer-facing text; Type and DueAt are
// supplementary machine-readable signals derived from it.
type PaymentSchedule struct {
	// ID is the identifier for this payment schedule, unique within its payment
	// term. Businesses SHOULD keep it stable across responses while the schedule
	// remains the same payment.
	ID string `json:"id"`

	// Amount is charged when this payment is taken, inclusive of tax and every
	// other charge, in the checkout currency's minor units (ISO 4217).
	Amount Amount `json:"amount"`

	// Description is the complete buyer-facing statement of when and how this
	// payment is due. Businesses MUST make this field sufficient on its own.
	Description Description `json:"description"`

	// DueAt is the absolute RFC 3339 date-time when this payment is due, when the
	// business can determine one at checkout. Supplementary to Description, never a
	// replacement for it. Omitted when the due date depends on a future event.
	DueAt *time.Time `json:"due_at,omitempty"`

	// Type is the timing class. "immediate" is the only value with defined
	// meaning; platforms MUST treat unrecognized values as not due at completion.
	Type string `json:"type"`
}
