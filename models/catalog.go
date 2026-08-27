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

// Catalog capability models.
//
// The Catalog capability lets platforms search and browse a business's
// product catalog before checkout. It is split into two independently
// negotiated capabilities:
//
//   - dev.ucp.shopping.catalog.search — free-text/filter product search
//   - dev.ucp.shopping.catalog.lookup — retrieve products/variants by ID
//
// Catalog operations return Product/Variant IDs that flow directly into
// checkout via LineItemCreateRequest.Item.ID — Variant.ID is the value a
// platform passes as item.id when creating a checkout. Catalog pricing and
// availability reflect current terms but are NOT transactional commitments;
// checkout is authoritative.

// Price is a monetary value with an explicit currency, enabling
// multi-currency catalogs.
type Price struct {
	// Amount is the value in ISO 4217 minor units (e.g. cents for USD).
	// Use 0 for free items.
	Amount int `json:"amount"`

	// Currency is the ISO 4217 currency code (e.g. "USD", "EUR", "GBP").
	Currency string `json:"currency"`
}

// PriceRange is a min/max price span, e.g. across a product's variants.
type PriceRange struct {
	// Min is the minimum price in the range.
	Min Price `json:"min"`

	// Max is the maximum price in the range.
	Max Price `json:"max"`
}

// Description holds product/variant copy in one or more formats. At least
// one format SHOULD be populated.
type Description struct {
	// Plain is the plain-text description.
	Plain string `json:"plain,omitempty"`

	// Markdown is the markdown-formatted description.
	Markdown string `json:"markdown,omitempty"`

	// Html is HTML-formatted content. Security: platforms MUST sanitize before
	// rendering (strip scripts, event handlers, and untrusted elements). Treat
	// all rich text as untrusted input. Added by the 2026-08-25 UCP release.
	Html string `json:"html,omitempty"`
}

// Media represents a product or variant media asset (image, video, 3D model).
type Media struct {
	// Type is the media type. Well-known values: "image", "video", "model_3d".
	Type string `json:"type"`

	// URL is the location of the media resource.
	URL string `json:"url"`

	// AltText is accessibility text describing the media.
	AltText string `json:"alt_text,omitempty"`

	// Width is the width in pixels (for images/video).
	Width int `json:"width,omitempty"`

	// Height is the height in pixels (for images/video).
	Height int `json:"height,omitempty"`
}

// OptionValue is a possible value for a product option (e.g. "Small", "Blue").
type OptionValue struct {
	// ID is an optional server-assigned identifier for this option value.
	// When present in a SelectedOption, the server SHOULD use it for
	// matching instead of label.
	ID string `json:"id,omitempty"`

	// Label is the display text for this option value.
	Label string `json:"label"`
}

// ProductOption is a configurable dimension of a product (e.g. "Size", "Color").
type ProductOption struct {
	// Name is the option name (e.g. "Size", "Color").
	Name string `json:"name"`

	// Values are the available values for this option.
	Values []OptionValue `json:"values"`
}

// SelectedOption is an option value that defines a specific variant
// (e.g. Color: Blue).
type SelectedOption struct {
	// Name is the option name (e.g. "Size").
	Name string `json:"name"`

	// ID is an optional option value identifier from OptionValue.ID. When
	// present, the server SHOULD use it for matching; name and label remain
	// required for display.
	ID string `json:"id,omitempty"`

	// Label is the selected option label (e.g. "Large").
	Label string `json:"label"`
}

// ProductCategory is a catalog taxonomy entry a product or variant belongs to.
// Named ProductCategory (rather than Category) to avoid ambiguity with
// vertical/category concepts in consuming applications.
type ProductCategory struct {
	// ID is an optional taxonomy identifier for this category.
	ID string `json:"id,omitempty"`

	// Name is the human-readable category name.
	Name string `json:"name,omitempty"`
}

// Rating is an aggregate product or variant rating.
type Rating struct {
	// Value is the average rating value.
	Value float64 `json:"value"`

	// ScaleMin is the minimum value on the rating scale (e.g. 1 for 1-5 stars).
	ScaleMin float64 `json:"scale_min,omitempty"`

	// ScaleMax is the maximum value on the rating scale (e.g. 5 for 5-star).
	ScaleMax float64 `json:"scale_max"`

	// Count is the number of reviews contributing to the rating.
	Count int `json:"count,omitempty"`
}

// Availability describes whether a variant can be purchased.
type Availability struct {
	// Available indicates whether the variant is purchasable.
	Available bool `json:"available"`

	// Status qualifies Available with fulfillment state. Well-known values:
	// "in_stock", "backorder", "preorder", "out_of_stock", "discontinued".
	// Added by the 2026-08-25 UCP release.
	Status string `json:"status,omitempty"`
}

// VariantMatchType indicates how a requested identifier resolved to a variant
// in a lookup response.
type VariantMatchType string

const (
	// VariantMatchExact indicates the identifier matched this variant directly.
	VariantMatchExact VariantMatchType = "exact"

	// VariantMatchFeatured indicates the server selected this variant as the
	// featured representative (e.g. when a product ID was supplied).
	VariantMatchFeatured VariantMatchType = "featured"
)

// VariantInput correlates a request identifier to the variant it resolved to
// in a lookup response. See the catalog spec's Client Correlation section.
type VariantInput struct {
	// Value is the request identifier that resolved to this variant.
	Value string `json:"value"`

	// Match indicates whether the resolution was exact or featured.
	Match VariantMatchType `json:"match,omitempty"`
}

// Variant is a purchasable item with specific option selections, price, and
// availability. Variant.ID is used as item.id in checkout.
type Variant struct {
	// ID is the global ID uniquely identifying this variant. Used as item.id
	// in checkout.
	ID string `json:"id"`

	// SKU is the business-assigned identifier for inventory and fulfillment.
	SKU string `json:"sku,omitempty"`

	// Handle is a URL-safe variant slug.
	Handle string `json:"handle,omitempty"`

	// Title is the variant display title (e.g. "Blue / Large").
	Title string `json:"title"`

	// Description is the variant description in one or more formats.
	Description *Description `json:"description,omitempty"`

	// URL is the canonical variant page URL.
	URL string `json:"url,omitempty"`

	// Categories are the taxonomy entries for this variant.
	Categories []ProductCategory `json:"categories,omitempty"`

	// Price is the current selling price.
	Price Price `json:"price"`

	// ListPrice is the list price before discounts (for strikethrough display).
	ListPrice *Price `json:"list_price,omitempty"`

	// Availability indicates whether the variant is purchasable.
	Availability *Availability `json:"availability,omitempty"`

	// Options are the option values that define this variant.
	Options []SelectedOption `json:"options,omitempty"`

	// Media are the variant media assets; the first element is featured.
	Media []Media `json:"media,omitempty"`

	// Rating is the variant rating.
	Rating *Rating `json:"rating,omitempty"`

	// Tags categorize the variant for search.
	Tags []string `json:"tags,omitempty"`

	// Inputs correlates request identifiers to this variant in lookup
	// responses (which request IDs resolved here, and how).
	Inputs []VariantInput `json:"inputs,omitempty"`

	// Metadata is business-defined custom data extending the standard model.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Product is a catalog item with one or more purchasable variants.
type Product struct {
	// ID is the global ID uniquely identifying this product.
	ID string `json:"id"`

	// Handle is a URL-safe slug for SEO-friendly URLs (e.g. "blue-runner-pro").
	// Use ID for stable API references.
	Handle string `json:"handle,omitempty"`

	// Title is the product title.
	Title string `json:"title"`

	// Description is the product description in one or more formats.
	Description *Description `json:"description,omitempty"`

	// URL is the canonical product page URL.
	URL string `json:"url,omitempty"`

	// Categories are the taxonomy entries for this product.
	Categories []ProductCategory `json:"categories,omitempty"`

	// PriceRange is the price range across all variants.
	PriceRange PriceRange `json:"price_range"`

	// ListPriceRange is the list price range before discounts.
	ListPriceRange *PriceRange `json:"list_price_range,omitempty"`

	// Media are the product media assets; the first element is featured.
	Media []Media `json:"media,omitempty"`

	// Options are the product options (Size, Color, etc.).
	Options []ProductOption `json:"options,omitempty"`

	// Variants are the purchasable variants of this product. The first item
	// is the featured variant for listings.
	Variants []Variant `json:"variants"`

	// Rating is the aggregate product rating.
	Rating *Rating `json:"rating,omitempty"`

	// Tags categorize the product for search.
	Tags []string `json:"tags,omitempty"`

	// Metadata is business-defined custom data extending the standard model.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ResponseCatalog represents UCP metadata for catalog responses, mirroring the
// ucp envelope on other capability responses.
type ResponseCatalog struct {
	// Version is the UCP protocol version.
	Version Version `json:"version,omitempty"`

	// Capabilities lists the active capabilities for this response.
	Capabilities []CapabilityResponse `json:"capabilities,omitempty"`
}

// CatalogSearchRequest is a request for the dev.ucp.shopping.catalog.search
// capability — free-text and/or filter-based product search.
type CatalogSearchRequest struct {
	// Query is the free-text search query. Optional when filters are supplied.
	Query string `json:"query,omitempty"`

	// Filters are structured constraints (e.g. price, category) applied to the
	// search. Structure is business-defined; common keys are documented in the
	// capability spec.
	Filters map[string]interface{} `json:"filters,omitempty"`

	// Context provides buyer/location hints for relevance and localization.
	Context *Context `json:"context,omitempty"`

	// Signals contains platform-provided environment data for authorization
	// and abuse prevention.
	Signals *Signals `json:"signals,omitempty"`

	// Attribution provides platform-emitted referral/conversion context.
	Attribution Attribution `json:"attribution,omitempty"`

	// Limit caps the number of products returned.
	Limit int `json:"limit,omitempty"`

	// Cursor is an opaque pagination cursor from a prior response.
	Cursor string `json:"cursor,omitempty"`
}

// CatalogSearchResponse is the response to a catalog search.
type CatalogSearchResponse struct {
	// UCP contains protocol metadata (active capabilities).
	UCP *ResponseCatalog `json:"ucp,omitempty"`

	// Products are the matching products. Empty (without messages) means a
	// valid query with no results.
	Products []Product `json:"products"`

	// NextCursor is an opaque pagination cursor for the next page, if any.
	NextCursor string `json:"next_cursor,omitempty"`

	// Messages contains errors, warnings, or informational notices.
	Messages []Message `json:"messages,omitempty"`
}

// CatalogLookupRequest is a request for the dev.ucp.shopping.catalog.lookup
// capability — resolve one or more identifiers to products/variants. Use when
// you already have IDs (saved lists, deep links, cart validation).
type CatalogLookupRequest struct {
	// IDs are the product or variant identifiers to resolve. Businesses MAY
	// also support SKU, handle, or URL.
	IDs []string `json:"ids"`

	// Context provides buyer/location hints for localization.
	Context *Context `json:"context,omitempty"`

	// Signals contains platform-provided environment data.
	Signals *Signals `json:"signals,omitempty"`
}

// CatalogLookupResponse is the response to a catalog lookup. Found products are
// returned; identifiers not found are reported via Messages (info, not_found).
type CatalogLookupResponse struct {
	// UCP contains protocol metadata (active capabilities).
	UCP *ResponseCatalog `json:"ucp,omitempty"`

	// Products are the resolved products. Each variant's Inputs array
	// correlates which requested identifiers resolved to it.
	Products []Product `json:"products"`

	// Messages contains errors, warnings, or informational notices.
	Messages []Message `json:"messages,omitempty"`
}

// GetProductRequest retrieves full detail for a single product or variant,
// the authoritative source for purchase decisions (dev.ucp.shopping.catalog.lookup,
// get_product / POST /catalog/product).
type GetProductRequest struct {
	// ID is the product or variant identifier.
	ID string `json:"id"`

	// SelectedOptions narrows variant resolution for interactive selection.
	SelectedOptions []SelectedOption `json:"selected_options,omitempty"`

	// Context provides buyer/location hints for localization.
	Context *Context `json:"context,omitempty"`

	// Signals contains platform-provided environment data.
	Signals *Signals `json:"signals,omitempty"`
}

// GetProductResponse is the response to a single-product lookup.
type GetProductResponse struct {
	// UCP contains protocol metadata (active capabilities).
	UCP *ResponseCatalog `json:"ucp,omitempty"`

	// Product is the resolved product with full detail, or nil if not found.
	Product *Product `json:"product,omitempty"`

	// Messages contains errors, warnings, or informational notices.
	Messages []Message `json:"messages,omitempty"`
}
