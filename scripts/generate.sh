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

# Schema directory (assumes ucp repo is sibling to go-sdk)
SCHEMA_DIR="$ROOT_DIR/../ucp/spec/schemas"
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

# Clean and create output directory
rm -rf "$GEN_DIR"
mkdir -p "$GEN_DIR"

# Collect all schema files
SCHEMA_FILES=""
for f in "$TYPES_DIR"/*.json; do
    if [ -f "$f" ]; then
        SCHEMA_FILES="$SCHEMA_FILES $f"
    fi
done

# Add top-level shopping schemas
for f in "$SCHEMA_DIR/shopping"/*.json; do
    if [ -f "$f" ]; then
        SCHEMA_FILES="$SCHEMA_FILES $f"
    fi
done

echo "Generating models from $(echo $SCHEMA_FILES | wc -w | tr -d ' ') schema files..."

# Generate all types into a single file
$GO_JSONSCHEMA \
    --package generated \
    --only-models \
    --struct-name-from-title \
    --tags json \
    --resolve-extension json \
    --capitalization ID \
    --capitalization URL \
    --capitalization URI \
    --capitalization API \
    --output "$GEN_DIR/models.go" \
    $SCHEMA_FILES 2>&1 || {
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

# Create temp file with header
echo "$HEADER" > "$GEN_DIR/models.go.tmp"
cat "$GEN_DIR/models.go" >> "$GEN_DIR/models.go.tmp"
mv "$GEN_DIR/models.go.tmp" "$GEN_DIR/models.go"

# Run gofmt
gofmt -w "$GEN_DIR/models.go"

# Post-processing is optional - the generator already handles most cases well
# Uncomment to run post-processing for additional cleanup:
# if [ -f "$SCRIPT_DIR/postprocess.go" ]; then
#     echo "Running post-processing..."
#     go run "$SCRIPT_DIR/postprocess.go" "$GEN_DIR/models.go"
# fi

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
