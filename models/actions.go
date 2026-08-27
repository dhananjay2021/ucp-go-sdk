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

// Actions represents outstanding extension-defined Action instances, keyed by
// reverse-domain Action type (not extension name). Each key maps to a non-empty,
// ordered slice of instances of that Action type; JSON preserves array order, and
// the declaring extension defines whether order carries processing semantics.
type Actions map[string][]ActionInstance

// ActionInstance represents one outstanding Action instance. The common fields are
// ID and an optional Config. The extension that declares the Action type defines
// the type-specific processing data carried under Config.
type ActionInstance struct {
	// ID is the identifier for this Action instance.
	ID string `json:"id"`

	// Config is configuration defined by the extension that declares this Action
	// type. Its structure is opaque here and interpreted by the declaring extension.
	Config map[string]interface{} `json:"config,omitempty"`
}
