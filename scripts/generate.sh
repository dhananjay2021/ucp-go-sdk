#!/bin/bash
# Copyright 2026 UCP Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Generate Go models from UCP JSON Schemas
# Usage: ./scripts/generate.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# UCP repository root. Defaults to a sibling checkout of the ucp repo, but can be
# overridden via the first argument or the UCP_REPO environment variable to point
# at a specific release checkout/worktree, e.g.:
#   ./scripts/generate.sh /path/to/ucp-release-2026-04-08
UCP_REPO="${1:-${UCP_REPO:-$ROOT_DIR/../ucp}}"

# Canonical JSON Schemas live under source/schemas in the ucp repo.
SCHEMA_DIR="$UCP_REPO/source/schemas"
TYPES_DIR="$SCHEMA_DIR/shopping/types"

# Output directory for generated code
GEN_DIR="$ROOT_DIR/models/generated"

# Check if schema directory exists
if [ ! -d "$SCHEMA_DIR" ]; then
    echo "Error: Schema directory not found at $SCHEMA_DIR"
    echo "Please ensure the UCP specification repository is available."
    echo ""
    echo "You can clone it with:"
    echo "  git clone https://github.com/Universal-Commerce-Protocol/ucp.git ../ucp"
    exit 1
fi

# Check if go-jsonschema is installed
GO_JSONSCHEMA=$(go env GOPATH)/bin/go-jsonschema
if [ ! -x "$GO_JSONSCHEMA" ]; then
    echo "go-jsonschema not found. Installing..."
    go install github.com/atombender/go-jsonschema@latest
fi

echo "=== UCP Go SDK Model Generation ==="
echo ""
echo "Schema directory: $SCHEMA_DIR"
echo "Output directory: $GEN_DIR"
echo ""

# Ensure output directory exists. Only the generated models.go is replaced;
# hand-maintained files in this package (e.g. doc.go) are preserved.
mkdir -p "$GEN_DIR"
rm -f "$GEN_DIR/models.go"

# The canonical schemas reference each other with relative $refs that
# go-jsonschema cannot resolve on its own (see scripts/bundle_schemas.py).
# Bundle them into a single self-contained schema first.
BUNDLE_FILE="$(mktemp -t ucp_bundle.XXXXXX.json)"
trap 'rm -f "$BUNDLE_FILE"' EXIT

echo "Bundling schemas from $SCHEMA_DIR..."
python3 "$SCRIPT_DIR/bundle_schemas.py" "$SCHEMA_DIR" "$BUNDLE_FILE"

echo "Generating models..."

# Generate all types into a single file. Type names come from the bundle's
# $defs keys, so --struct-name-from-title is intentionally omitted.
$GO_JSONSCHEMA \
    --package generated \
    --only-models \
    --tags json \
    --capitalization ID \
    --capitalization URL \
    --capitalization URI \
    --capitalization API \
    --output "$GEN_DIR/models.go" \
    "$BUNDLE_FILE" 2>&1 || {
        echo "Warning: Some schemas may have issues. Continuing..."
    }

# Check if file was generated
if [ ! -f "$GEN_DIR/models.go" ]; then
    echo "Error: Generation failed - no output file created"
    exit 1
fi

# Add file header
HEADER="// Code generated from UCP JSON Schemas. DO NOT EDIT.
// Source: https://github.com/Universal-Commerce-Protocol/ucp
// Generator: go-jsonschema (https://github.com/atombender/go-jsonschema)
//
// This file contains auto-generated types that match the UCP specification.
// For custom extensions and helper methods, see the parent models/ package.
"

# Keep the generated package dependency-free. go-jsonschema maps JSON Schema
# "format": "date" to its own types.SerializableDate, which would pull in the
# generator module as a runtime dependency. Represent bare dates as strings
# (date-time already maps to the standard library time.Time) and drop the import.
sed -i '' \
    -e 's#types\.SerializableDate#string#g' \
    -e '\#"github.com/atombender/go-jsonschema/pkg/types"#d' \
    "$GEN_DIR/models.go"

# Create temp file with header
echo "$HEADER" > "$GEN_DIR/models.go.tmp"
cat "$GEN_DIR/models.go" >> "$GEN_DIR/models.go.tmp"
mv "$GEN_DIR/models.go.tmp" "$GEN_DIR/models.go"

# Run gofmt
gofmt -w "$GEN_DIR/models.go"

echo ""
echo "=== Generation Complete ==="
echo ""
echo "Generated: $GEN_DIR/models.go"
echo ""
echo "Next steps:"
echo "  1. Review generated types"
echo "  2. Run 'go build ./...' to verify"
echo "  3. Update models/*.go to use generated types where appropriate"
echo ""
