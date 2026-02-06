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

	// Checkout Handlers
	createCheckoutHandler   func(http.ResponseWriter, *http.Request)
	getCheckoutHandler      func(http.ResponseWriter, *http.Request)
	updateCheckoutHandler   func(http.ResponseWriter, *http.Request)
	completeCheckoutHandler func(http.ResponseWriter, *http.Request)
	cancelCheckoutHandler   func(http.ResponseWriter, *http.Request)
	getOrderHandler         func(http.ResponseWriter, *http.Request)

	// Cart Handlers
	createCartHandler func(http.ResponseWriter, *http.Request)
	getCartHandler    func(http.ResponseWriter, *http.Request)
	updateCartHandler func(http.ResponseWriter, *http.Request)
	deleteCartHandler func(http.ResponseWriter, *http.Request)
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

	// Cart routes
	s.mux.HandleFunc("POST /carts", s.handleCreateCart)
	s.mux.HandleFunc("GET /carts/{id}", s.handleGetCart)
	s.mux.HandleFunc("PATCH /carts/{id}", s.handleUpdateCart)
	s.mux.HandleFunc("DELETE /carts/{id}", s.handleDeleteCart)

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

// CreateCartHandler is a function that handles cart creation.
type CreateCartHandler func(r *http.Request, req *models.CartCreateRequest) (*models.CartResponse, error)

// GetCartHandler is a function that handles cart retrieval.
type GetCartHandler func(r *http.Request, id string) (*models.CartResponse, error)

// UpdateCartHandler is a function that handles cart updates.
type UpdateCartHandler func(r *http.Request, id string, req *models.CartUpdateRequest) (*models.CartResponse, error)

// DeleteCartHandler is a function that handles cart deletion.
type DeleteCartHandler func(r *http.Request, id string) error

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

// HandleCreateCart registers a handler for creating carts.
func (s *Server) HandleCreateCart(handler CreateCartHandler) {
	s.createCartHandler = func(w http.ResponseWriter, r *http.Request) {
		var req models.CartCreateRequest
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

// HandleGetCart registers a handler for retrieving carts.
func (s *Server) HandleGetCart(handler GetCartHandler) {
	s.getCartHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		resp, err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, http.StatusOK, resp)
	}
}

// HandleUpdateCart registers a handler for updating carts.
func (s *Server) HandleUpdateCart(handler UpdateCartHandler) {
	s.updateCartHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req models.CartUpdateRequest
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

// HandleDeleteCart registers a handler for deleting carts.
func (s *Server) HandleDeleteCart(handler DeleteCartHandler) {
	s.deleteCartHandler = func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		err := handler(r, id)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
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

func (s *Server) handleCreateCart(w http.ResponseWriter, r *http.Request) {
	if s.createCartHandler != nil {
		s.createCartHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Cart creation not implemented")
	}
}

func (s *Server) handleGetCart(w http.ResponseWriter, r *http.Request) {
	if s.getCartHandler != nil {
		s.getCartHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Cart retrieval not implemented")
	}
}

func (s *Server) handleUpdateCart(w http.ResponseWriter, r *http.Request) {
	if s.updateCartHandler != nil {
		s.updateCartHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Cart update not implemented")
	}
}

func (s *Server) handleDeleteCart(w http.ResponseWriter, r *http.Request) {
	if s.deleteCartHandler != nil {
		s.deleteCartHandler(w, r)
	} else {
		WriteError(w, http.StatusNotImplemented, "not_implemented", "Cart deletion not implemented")
	}
}
