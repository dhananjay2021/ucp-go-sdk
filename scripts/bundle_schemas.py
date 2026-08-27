#!/usr/bin/env python3
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

"""Bundle the UCP JSON Schemas into a single self-contained document.

The canonical UCP schemas under ``source/schemas`` are split across many files
that reference each other with relative ``$ref`` paths (e.g. ``category.json``,
``types/message.json``, ``../ucp.json#/$defs/foo``). The upstream ucp build tool
resolves these with a bespoke multi-directory search that ``go-jsonschema`` does
not replicate, so feeding the raw files to ``go-jsonschema`` fails on transitive
sibling refs.

This script resolves every ``$ref`` using standard JSON Schema base-URI rules
(each file declares an absolute ``$id``) and inlines the referenced schemas into
a single document under ``$defs`` with all refs rewritten to local
``#/$defs/<Name>`` pointers. The resulting bundle is safe to hand to
``go-jsonschema``.

Usage:
    bundle_schemas.py <schemas_root> <output_file>
"""

import json
import os
import re
import sys
from urllib.parse import urldefrag, urljoin

BASE_URL = "https://ucp.dev/schemas/"

# Top-level schema directories (relative to the schemas root) to seed bundling
# from. As of the 2026-08-25 release many capabilities and primitive types live
# under common/ (payment, location, loyalty, shared types) in addition to the
# shopping vertical.
ENTRY_GLOBS = ("shopping", "shopping/types", "common", "common/types")


def load_all(root):
    """Index every schema file under ``root`` by its declared ``$id``."""
    by_id = {}
    for dirpath, _dirs, files in os.walk(root):
        for name in files:
            if not name.endswith(".json"):
                continue
            path = os.path.join(dirpath, name)
            with open(path, encoding="utf-8") as handle:
                data = json.load(handle)
            schema_id = data.get("$id")
            if schema_id:
                by_id[schema_id] = data
    return by_id


def resolve_pointer(doc, fragment):
    """Return the sub-schema referenced by a JSON pointer ``fragment``."""
    if not fragment:
        return doc
    node = doc
    for raw in fragment.lstrip("/").split("/"):
        part = raw.replace("~1", "/").replace("~0", "~")
        node = node[part]
    return node


class Bundler:
    def __init__(self, docs_by_id):
        self.docs_by_id = docs_by_id
        self.defs = {}          # local $defs key -> inlined schema body
        self.key_by_uri = {}    # absolute uri (with fragment) -> local key
        self.used_keys = set()

    def _unique_key(self, schema, base_url, fragment):
        title = schema.get("title") if isinstance(schema, dict) else None
        if title:
            candidate = re.sub(r"[^0-9a-zA-Z]+", " ", title).title().replace(" ", "")
        else:
            stem = base_url[len(BASE_URL):] if base_url.startswith(BASE_URL) else base_url
            stem = stem.rsplit("/", 1)[-1].replace(".json", "")
            frag = fragment.replace("$defs", "").strip("/").replace("/", "_")
            raw = f"{stem}_{frag}" if frag else stem
            candidate = re.sub(r"[^0-9a-zA-Z]+", " ", raw).title().replace(" ", "")
        candidate = candidate or "Schema"
        if candidate[0].isdigit():
            candidate = "N" + candidate
        key = candidate
        suffix = 2
        while key in self.used_keys:
            key = f"{candidate}{suffix}"
            suffix += 1
        self.used_keys.add(key)
        return key

    def register(self, abs_uri, current_base):
        """Ensure ``abs_uri`` is inlined; return its local ``#/$defs`` ref."""
        base_url, fragment = urldefrag(abs_uri)
        if base_url not in self.docs_by_id:
            raise KeyError(f"unknown schema $id: {base_url} (from {current_base})")
        canonical = f"{base_url}#{fragment}"
        if canonical in self.key_by_uri:
            return f"#/$defs/{self.key_by_uri[canonical]}"

        doc = self.docs_by_id[base_url]
        target = resolve_pointer(doc, fragment)
        key = self._unique_key(target, base_url, fragment)
        self.key_by_uri[canonical] = key
        # Reserve the slot before recursing so cyclic refs terminate.
        self.defs[key] = None
        self.defs[key] = self._inline(target, base_url)
        return f"#/$defs/{key}"

    def _inline(self, node, base_url):
        """Deep-copy ``node`` rewriting external refs to local ``$defs`` refs."""
        if isinstance(node, dict):
            result = {}
            for prop, value in node.items():
                if prop in ("$id", "$schema", "$defs", "definitions"):
                    # Drop embedded identity/metadata and local $defs; referenced
                    # $defs entries are hoisted individually via register().
                    continue
                if prop == "$ref" and isinstance(value, str):
                    abs_uri = urljoin(base_url, value)
                    result[prop] = self.register(abs_uri, base_url)
                else:
                    result[prop] = self._inline(value, base_url)
            return result
        if isinstance(node, list):
            return [self._inline(item, base_url) for item in node]
        return node


def main():
    if len(sys.argv) != 3:
        sys.stderr.write("usage: bundle_schemas.py <schemas_root> <output_file>\n")
        return 1
    schemas_root, output_file = sys.argv[1], sys.argv[2]

    docs_by_id = load_all(schemas_root)
    bundler = Bundler(docs_by_id)

    entries = []
    for rel in ENTRY_GLOBS:
        directory = os.path.join(schemas_root, rel)
        if not os.path.isdir(directory):
            continue
        for name in sorted(os.listdir(directory)):
            if name.endswith(".json"):
                entries.append(os.path.join(directory, name))

    for path in entries:
        with open(path, encoding="utf-8") as handle:
            schema_id = json.load(handle).get("$id")
        if schema_id:
            bundler.register(schema_id, schema_id)

    bundle = {
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://ucp.dev/schemas/_bundle.json",
        "$defs": bundler.defs,
    }
    with open(output_file, "w", encoding="utf-8") as handle:
        json.dump(bundle, handle, indent=2)
        handle.write("\n")

    sys.stderr.write(
        f"Bundled {len(entries)} entry schemas into {len(bundler.defs)} "
        f"definitions -> {output_file}\n"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
