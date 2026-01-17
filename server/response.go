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
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// APIError represents an error that can be returned from handlers.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    any
}

func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error.
func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// NotFoundError creates a 404 not found error.
func NotFoundError(message string) *APIError {
	return NewAPIError(http.StatusNotFound, "not_found", message)
}

// BadRequestError creates a 400 bad request error.
func BadRequestError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, "bad_request", message)
}

// UnauthorizedError creates a 401 unauthorized error.
func UnauthorizedError(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "unauthorized", message)
}

// ForbiddenError creates a 403 forbidden error.
func ForbiddenError(message string) *APIError {
	return NewAPIError(http.StatusForbidden, "forbidden", message)
}

// ConflictError creates a 409 conflict error.
func ConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, "conflict", message)
}

// InternalError creates a 500 internal server error.
func InternalError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "internal_error", message)
}

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// WriteError writes an error response.
func WriteError(w http.ResponseWriter, statusCode int, code, message string) {
	WriteJSON(w, statusCode, ErrorResponse{
		Error:   code,
		Message: message,
	})
}

// handleError handles errors from handlers.
func handleError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		resp := ErrorResponse{
			Error:   apiErr.Code,
			Message: apiErr.Message,
			Details: apiErr.Details,
		}
		WriteJSON(w, apiErr.StatusCode, resp)
		return
	}

	// Default to internal server error
	WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
}
