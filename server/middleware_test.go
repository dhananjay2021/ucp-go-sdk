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

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddlewareAddsVaryOriginForAllowedOrigin(t *testing.T) {
	handler := CORSMiddleware([]string{"https://platform.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}))

	req := httptest.NewRequest(http.MethodGet, "/checkout", nil)
	req.Header.Set("Origin", "https://platform.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://platform.example" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://platform.example")
	}
	if got := rec.Header().Values("Vary"); len(got) == 0 || got[0] != "Origin" {
		t.Fatalf("Vary = %q, want Origin", got)
	}
}

func TestCORSMiddlewarePreflightAddsVaryOriginForAllowedOrigin(t *testing.T) {
	handler := CORSMiddleware([]string{"https://platform.example"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run for preflight")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/checkout", nil)
	req.Header.Set("Origin", "https://platform.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://platform.example" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://platform.example")
	}
	if got := rec.Header().Values("Vary"); len(got) == 0 || got[0] != "Origin" {
		t.Fatalf("Vary = %q, want Origin", got)
	}
}
