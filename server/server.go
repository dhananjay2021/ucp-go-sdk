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

// Package server provides HTTP handler utilities for implementing UCP endpoints.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/dhananjay2021/ucp-go-sdk/extensions"
	"github.com/dhananjay2021/ucp-go-sdk/models"
)

// Config contains server configuration.
type Config struct {
	// Version is the UCP protocol version (YYYY-MM-DD format).
	Version models.Version

	// Capabilities are the supported capabilities.
	Capabilities []models.CapabilityDiscovery

	// Services are the service definitions.
	Services models.Services

	// SigningKeys are the public keys for signature verification.
	SigningKeys []models.JWK

	// PaymentHandlers are the supported payment handlers.
	PaymentHandlers []models.PaymentHandlerResponse
}

// Server is a UCP server that handles HTTP requests.
type Server struct {
	config Config
	mux    *http.ServeMux

	// Handlers
	createCheckoutHandler   func(http.ResponseWriter, *http.Request)
	getCheckoutHandler      func(http.ResponseWriter, *http.Request)
	updateCheckoutHandler   func(http.ResponseWriter, *http.Request)
	completeCheckoutHandler func(http.ResponseWriter, *http.Request)
	cancelCheckoutHandler   func(http.ResponseWriter, *http.Request)
	getOrderHandler         func(http.ResponseWriter, *http.Request)
}

// NewServer creates a new UCP server.
func NewServer(config Config) *Server {
	s := &Server{
		config: config,
		mux:    http.NewServeMux(),
	}

	// Register routes
	s.mux.HandleFunc("GET /.well-known/ucp", s.handleDiscovery)
	s.mux.HandleFunc("POST /checkout-sessions", s.handleCreateCheckout)
	s.mux.HandleFunc("GET /checkout-sessions/{id}", s.handleGetCheckout)
	s.mux.HandleFunc("PATCH /checkout-sessions/{id}", s.handleUpdateCheckout)
	s.mux.HandleFunc("POST /checkout-sessions/{id}/complete", s.handleCompleteCheckout)
	s.mux.HandleFunc("POST /checkout-sessions/{id}/cancel", s.handleCancelCheckout)
	s.mux.HandleFunc("GET /orders/{id}", s.handleGetOrder)

	return s
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// CreateCheckoutHandler is a function that handles checkout creation.
type CreateCheckoutHandler func(r *http.Request, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error)

// GetCheckoutHandler is a function that handles checkout retrieval.
type GetCheckoutHandler func(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error)

// UpdateCheckoutHandler is a function that handles checkout updates.
type UpdateCheckoutHandler func(r *http.Request, id string, req *extensions.ExtendedCheckoutUpdateRequest) (*extensions.ExtendedCheckoutResponse, error)

// CompleteCheckoutHandler is a function that handles checkout completion.
type CompleteCheckoutHandler func(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error)

// CancelCheckoutHandler is a function that handles checkout cancellation.
type CancelCheckoutHandler func(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error)

// GetOrderHandler is a function that handles order retrieval.
type GetOrderHandler func(r *http.Request, id string) (*models.Order, error)

// HandleCreateCheckout registers a handler for creating checkout sessions.
func (s *Server) HandleCreateCheckout(handler CreateCheckoutHandler) {
	s.createCheckoutHandler = func(w http.ResponseWriter, r *http.Request) {
		var req extensions.ExtendedCheckoutCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid_request", "Failed to parse request body")
			return
		}

		resp, err := handler(r, &req)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusCreated, resp)
	}
}

// HandleGetCheckout registers a handler for retrieving checkout sessions.
func (s *Server) HandleGetCheckout(handler GetCheckoutHandler) {
	s.getCheckoutHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// HandleUpdateCheckout registers a handler for updating checkout sessions.
func (s *Server) HandleUpdateCheckout(handler UpdateCheckoutHandler) {
	s.updateCheckoutHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req extensions.ExtendedCheckoutUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "invalid_request", "Failed to parse request body")
			return
		}

		resp, err := handler(r, id, &req)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// HandleCompleteCheckout registers a handler for completing checkout sessions.
func (s *Server) HandleCompleteCheckout(handler CompleteCheckoutHandler) {
	s.completeCheckoutHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// HandleCancelCheckout registers a handler for canceling checkout sessions.
func (s *Server) HandleCancelCheckout(handler CancelCheckoutHandler) {
	s.cancelCheckoutHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// HandleGetOrder registers a handler for retrieving orders.
func (s *Server) HandleGetOrder(handler GetOrderHandler) {
	s.getOrderHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// Internal route handlers

func (s *Server) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	profile := models.UCPProfile{
		UCP: models.DiscoveryProfile{
			Version:      s.config.Version,
			Services:     s.config.Services,
			Capabilities: s.config.Capabilities,
		},
		SigningKeys: s.config.SigningKeys,
	}

	if len(s.config.PaymentHandlers) > 0 {
		profile.Payment = &models.PaymentConfig{
			Handlers: s.config.PaymentHandlers,
		}
	}

	WriteJSON(w, http.StatusOK, profile)
}

func (s *Server) handleCreateCheckout(w http.ResponseWriter, r *http.Request) {
	if s.createCheckoutHandler != nil {
		s.createCheckoutHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Checkout creation not implemented")
	}
}

func (s *Server) handleGetCheckout(w http.ResponseWriter, r *http.Request) {
	if s.getCheckoutHandler != nil {
		s.getCheckoutHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Checkout retrieval not implemented")
	}
}

func (s *Server) handleUpdateCheckout(w http.ResponseWriter, r *http.Request) {
	if s.updateCheckoutHandler != nil {
		s.updateCheckoutHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Checkout update not implemented")
	}
}

func (s *Server) handleCompleteCheckout(w http.ResponseWriter, r *http.Request) {
	if s.completeCheckoutHandler != nil {
		s.completeCheckoutHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Checkout completion not implemented")
	}
}

func (s *Server) handleCancelCheckout(w http.ResponseWriter, r *http.Request) {
	if s.cancelCheckoutHandler != nil {
		s.cancelCheckoutHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Checkout cancellation not implemented")
	}
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	if s.getOrderHandler != nil {
		s.getOrderHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Order retrieval not implemented")
	}
}
