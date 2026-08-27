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

// Payment authentication Action types contributed by the Payment Authentication
// extension (UCP spec change #458). These reverse-domain identifiers key entries
// in a checkout's actions map.
const (
	// ActionTypeDeviceDataCollection is the 3DS device data collection Action
	// type.
	ActionTypeDeviceDataCollection = "dev.ucp.common.payment.device_data_collection"

	// ActionTypeThreeDSChallenge is the 3DS challenge Action type.
	ActionTypeThreeDSChallenge = "dev.ucp.common.payment.three_ds_challenge"
)

// PaymentDeviceDataCollectionConfig is the config for a 3DS device data
// collection Action. It runs on an invisible surface during payment
// authentication.
type PaymentDeviceDataCollectionConfig struct {
	// PaymentInstrumentID is the ID of the payment instrument in the containing
	// checkout associated with this device data collection Action.
	PaymentInstrumentID string `json:"payment_instrument_id"`

	// URL is the URL for the invisible device data collection surface.
	URL string `json:"url"`
}

// PaymentThreeDSChallengeConfig is the config for a 3DS challenge Action. It runs
// on a buyer-facing surface during payment authentication.
type PaymentThreeDSChallengeConfig struct {
	// PaymentInstrumentID is the ID of the payment instrument in the containing
	// checkout associated with this 3DS challenge Action.
	PaymentInstrumentID string `json:"payment_instrument_id"`

	// URL is the URL for the buyer-facing 3DS challenge surface.
	URL string `json:"url"`
}
