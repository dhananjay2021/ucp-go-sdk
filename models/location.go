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

// Geo represents WGS 84 geographic coordinates in decimal degrees.
type Geo struct {
	// Latitude is the WGS 84 latitude in decimal degrees (range -90 to 90).
	Latitude float64 `json:"latitude"`

	// Longitude is the WGS 84 longitude in decimal degrees (range -180 to 180).
	Longitude float64 `json:"longitude"`
}

// Locality is a coarse geographic location — country, region, and postal code.
// A lightweight alternative to a full postal address.
type Locality struct {
	// AddressCountry is the country as a 2-letter ISO 3166-1 alpha-2 code
	// (e.g. "US"). A 3-letter alpha-3 code or full country name MAY also be used.
	AddressCountry *string `json:"address_country,omitempty"`

	// AddressRegion is the first-level administrative region within the country
	// (e.g. a state or province such as California).
	AddressRegion *string `json:"address_region,omitempty"`

	// PostalCode is the postal code (e.g. "94043").
	PostalCode *string `json:"postal_code,omitempty"`
}

// Amenity is buyer-facing presentation metadata for one amenity identifier. The
// containing map key, not this metadata, defines amenity identity and filter
// matching.
type Amenity struct {
	// Description is a short, plain-text, buyer-facing label or phrase for the
	// amenity, suitable for direct use in a compact list (e.g., 'Curbside
	// pickup'). This content does not participate in amenity identity or filter
	// matching.
	Description string `json:"description"`
}

// LocationAmenities are static features, services, or capabilities of a
// Location, keyed by reverse-domain amenity identifier (e.g.,
// "dev.ucp.amenity.wi_fi"). Each value provides a buyer-facing description; the
// key alone defines amenity identity and filter matching.
type LocationAmenities map[string]Amenity

// DailyHourDay is a stable UCP day-of-week identifier for a recurring operating
// interval. It is not localized display text.
type DailyHourDay string

const (
	// DailyHourDayMonday identifies Monday.
	DailyHourDayMonday DailyHourDay = "monday"

	// DailyHourDayTuesday identifies Tuesday.
	DailyHourDayTuesday DailyHourDay = "tuesday"

	// DailyHourDayWednesday identifies Wednesday.
	DailyHourDayWednesday DailyHourDay = "wednesday"

	// DailyHourDayThursday identifies Thursday.
	DailyHourDayThursday DailyHourDay = "thursday"

	// DailyHourDayFriday identifies Friday.
	DailyHourDayFriday DailyHourDay = "friday"

	// DailyHourDaySaturday identifies Saturday.
	DailyHourDaySaturday DailyHourDay = "saturday"

	// DailyHourDaySunday identifies Sunday.
	DailyHourDaySunday DailyHourDay = "sunday"
)

// DailyHour is a regular weekly operating interval. Its Day, Opens, and Closes
// are recurring local civil values interpreted in the containing Location's
// timezone. Multiple entries for the same day support split shifts.
type DailyHour struct {
	// Day is the day-of-week on which this recurring local civil-time interval
	// begins in the containing Location's timezone.
	Day DailyHourDay `json:"day"`

	// Opens is the opening time in 24-hour HH:MM format.
	Opens string `json:"opens"`

	// Closes is the closing time in 24-hour HH:MM format.
	Closes string `json:"closes"`
}

// ExceptionHour is a date-specific operating interval or full closure. Its
// ValidFrom, ValidThrough, Opens, and Closes are local civil values interpreted
// in the containing Location's timezone. Date bounds are inclusive.
type ExceptionHour struct {
	// Title is a short human-readable heading naming the exception (for example,
	// 'Thanksgiving'). Presentation metadata that does not affect schedule
	// evaluation.
	Title *string `json:"title,omitempty"`

	// ValidFrom is the first local civil date (YYYY-MM-DD) to which this
	// exception applies, interpreted in the containing Location's timezone.
	ValidFrom string `json:"valid_from"`

	// ValidThrough is the last local civil date (YYYY-MM-DD) to which this
	// exception applies, interpreted in the containing Location's timezone.
	ValidThrough string `json:"valid_through"`

	// Opens is the opening time in 24-hour HH:MM format. When omitted (together
	// with Closes) the exception represents a full closure.
	Opens *string `json:"opens,omitempty"`

	// Closes is the closing time in 24-hour HH:MM format.
	Closes *string `json:"closes,omitempty"`
}

// LocationSummary is a summary of a physical business location.
type LocationSummary struct {
	// ID is a stable, opaque, Business-scoped Location identifier.
	ID string `json:"id"`

	// Name is the buyer-facing, Business-owned display name.
	Name string `json:"name"`

	// Address is the physical address of the location.
	Address *PostalAddress `json:"address,omitempty"`
}

// Location is the full, rich representation of a physical business location. It
// builds on LocationSummary with discovery-centric details such as geographic
// coordinates, operating hours, timezone, and amenities.
type Location struct {
	// ID is a stable, opaque, Business-scoped Location identifier.
	ID string `json:"id"`

	// Name is the buyer-facing, Business-owned display name.
	Name string `json:"name"`

	// Address is the physical address of the location.
	Address *PostalAddress `json:"address,omitempty"`

	// Geo holds the geographic coordinates for the location.
	Geo *Geo `json:"geo,omitempty"`

	// Amenities are static features, services, or capabilities of the Location,
	// keyed by reverse-domain amenity identifier.
	Amenities LocationAmenities `json:"amenities,omitempty"`

	// Hours are regular weekly operating hours whose day and time values use this
	// Location's canonical local civil-time frame. Multiple entries for the same
	// day support split shifts. Omission means the regular schedule is unknown.
	Hours []DailyHour `json:"hours,omitempty"`

	// ExceptionHours are date-specific operating-hour exceptions, including full
	// closures, whose date and time values use this Location's canonical local
	// civil-time frame.
	ExceptionHours []ExceptionHour `json:"exception_hours,omitempty"`

	// Timezone is the Business-owned IANA Time Zone Database identifier (e.g.,
	// 'America/New_York') defining this Location's canonical local civil-time
	// frame for all returned schedule day, time, and date fields. Required when
	// Hours or ExceptionHours is present.
	Timezone *string `json:"timezone,omitempty"`
}

// LocationDistance is an explicit-center inclusive-radius predicate. The
// Business compares the unrounded shortest WGS 84 ellipsoidal geodesic distance
// in meters from Center to the Location's authoritative geo; a value less than
// or equal to Max matches. A Business unable to honor the request MUST reject it
// rather than clamp the radius or substitute an operand.
type LocationDistance struct {
	// Center is the explicit center of the radius. The Platform MUST supply it;
	// the Business MUST NOT derive it from context, signals, an IP address, or
	// serves.
	Center Geo `json:"center"`

	// Max is the inclusive maximum distance in meters (RFC 7035 distance unit).
	Max float64 `json:"max"`
}

// LocationServesAddress is the coarse locality of a service target. At least one
// of AddressCountry, AddressRegion, or PostalCode MUST be supplied.
type LocationServesAddress struct {
	// AddressCountry is the country as a 2-letter ISO 3166-1 alpha-2 code
	// (e.g. "US"). A 3-letter alpha-3 code or full country name MAY also be used.
	AddressCountry *string `json:"address_country,omitempty"`

	// AddressRegion is the first-level administrative region within the country
	// (e.g. a state or province such as California).
	AddressRegion *string `json:"address_region,omitempty"`

	// PostalCode is the postal code (e.g. "94043").
	PostalCode *string `json:"postal_code,omitempty"`
}

// LocationServes is an authoritative service-target predicate. The Platform MUST
// supply exactly one target form (Point or Address). A Business that cannot
// evaluate a well-formed target MUST reject the request rather than ignore it,
// fall back, or broaden results.
type LocationServes struct {
	// Point holds the WGS 84 coordinates of the service target.
	Point *Geo `json:"point,omitempty"`

	// Address holds the coarse locality of the service target.
	Address *LocationServesAddress `json:"address,omitempty"`
}

// LocationFilterHours filters by operating hours, evaluated at the one supplied
// instant.
type LocationFilterHours struct {
	// OpenAt is the RFC 3339 instant at which matching Locations must be open.
	// The Business evaluates it exactly as supplied using each Location's
	// authoritative timezone; the supplied offset does not identify that
	// timezone.
	OpenAt time.Time `json:"open_at"`
}

// LocationFilter holds criteria to narrow Location Search and Lookup results.
// All supplied filters combine with AND.
type LocationFilter struct {
	// Hours filters by operating hours, evaluated at the one supplied instant.
	Hours *LocationFilterHours `json:"hours,omitempty"`

	// Amenities filters by amenity identifier. A Location matches only when its
	// amenities map contains every supplied reverse-domain identifier as an exact
	// key; descriptions and namespace prefixes do not participate in matching.
	Amenities []string `json:"amenities,omitempty"`

	// Items is a current item-availability filter. A candidate Location matches
	// only when the Business can currently provide every referenced opaque
	// Business-scoped item identifier at that Location; all references AND
	// together.
	Items []string `json:"items,omitempty"`
}

// LocationPaginationRequest holds cursor-based pagination parameters for
// requests.
type LocationPaginationRequest struct {
	// Cursor is an opaque cursor from a previous response.
	Cursor *string `json:"cursor,omitempty"`

	// Limit is the requested page size, not a guaranteed result count. When
	// omitted, the Business applies a default page size.
	Limit *int `json:"limit,omitempty"`
}

// LocationPaginationResponse holds cursor-based pagination information in
// responses.
type LocationPaginationResponse struct {
	// Cursor fetches the next page of results. MUST be present when HasNextPage
	// is true.
	Cursor *string `json:"cursor,omitempty"`

	// HasNextPage indicates whether more results are available.
	HasNextPage bool `json:"has_next_page"`

	// TotalCount is the total number of matching items, if available.
	TotalCount *int `json:"total_count,omitempty"`
}

// ResponseLocation represents UCP metadata for location responses, mirroring the
// ucp envelope on other capability responses.
type ResponseLocation struct {
	// Version is the UCP protocol version.
	Version Version `json:"version,omitempty"`

	// Capabilities lists the active capabilities for this response.
	Capabilities []CapabilityResponse `json:"capabilities,omitempty"`
}

// LocationSearchRequest is a request for the dev.ucp.common.location.search
// capability. The distance and serves relations and every supplied filters
// predicate combine with AND; query does not relax them.
type LocationSearchRequest struct {
	// Query is a free-text search query for natural language location search
	// (e.g., 'restaurants near me that deliver', 'hotels with pool').
	Query *string `json:"query,omitempty"`

	// Context provides buyer/location hints for relevance and localization.
	Context *Context `json:"context,omitempty"`

	// Signals contains platform-provided environment data for authorization and
	// abuse prevention.
	Signals *Signals `json:"signals,omitempty"`

	// Distance is an optional explicit-center radius predicate.
	Distance *LocationDistance `json:"distance,omitempty"`

	// Serves is an optional authoritative service-target predicate.
	Serves *LocationServes `json:"serves,omitempty"`

	// Filters holds structured constraints applied to the search.
	Filters *LocationFilter `json:"filters,omitempty"`

	// Pagination holds cursor-based pagination parameters.
	Pagination *LocationPaginationRequest `json:"pagination,omitempty"`
}

// LocationSearchResponse is the response to a location search.
type LocationSearchResponse struct {
	// UCP contains protocol metadata (active capabilities).
	UCP *ResponseLocation `json:"ucp,omitempty"`

	// Locations are the locations matching the search criteria.
	Locations []Location `json:"locations"`

	// Pagination holds cursor-based pagination information for the next page.
	Pagination *LocationPaginationResponse `json:"pagination,omitempty"`

	// Messages contains errors, warnings, or informational notices about the
	// search results.
	Messages []Message `json:"messages,omitempty"`
}

// LocationLookupRequest is a request for the dev.ucp.common.location.lookup
// capability. The Business resolves and deduplicates IDs before applying
// Distance, Serves, and every supplied Filters predicate; all structured
// predicates combine with AND.
type LocationLookupRequest struct {
	// IDs are identifiers of the Locations to look up. The Business MUST support
	// canonical Location.ID values and MAY support secondary or alias
	// identifiers.
	IDs []string `json:"ids"`

	// Distance is an optional explicit-center radius predicate applied after ID
	// resolution.
	Distance *LocationDistance `json:"distance,omitempty"`

	// Serves is an optional authoritative service-target predicate applied after
	// ID resolution.
	Serves *LocationServes `json:"serves,omitempty"`

	// Filters holds structured constraints applied after ID resolution.
	Filters *LocationFilter `json:"filters,omitempty"`

	// Context provides buyer/location hints for relevance and localization.
	Context *Context `json:"context,omitempty"`

	// Signals contains platform-provided environment data.
	Signals *Signals `json:"signals,omitempty"`
}

// LocationLookupInput preserves one request identifier exactly as supplied,
// correlating it to a resolved Location.
type LocationLookupInput struct {
	// ID is the identifier exactly as supplied in the lookup request.
	ID string `json:"id"`
}

// LocationLookupResult is a Location with required correlation metadata for
// lookup responses.
type LocationLookupResult struct {
	Location

	// Inputs records which request identifiers resolved to this Location. Each
	// entry preserves one identifier exactly as supplied in the request.
	Inputs []LocationLookupInput `json:"inputs"`
}

// LocationLookupResponse is the response to a batch location lookup. It may
// contain fewer Locations if some identifiers do not resolve or are filtered
// out, or more if one identifier resolves to multiple Locations.
type LocationLookupResponse struct {
	// UCP contains protocol metadata (active capabilities).
	UCP *ResponseLocation `json:"ucp,omitempty"`

	// Locations are the locations matching the requested identifiers and
	// refinements. When multiple identifiers resolve to the same Location, one
	// returned Location carries all corresponding Inputs entries.
	Locations []LocationLookupResult `json:"locations"`

	// Messages contains errors, warnings, or informational notices about the
	// requested Locations.
	Messages []Message `json:"messages,omitempty"`
}
