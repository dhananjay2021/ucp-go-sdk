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
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/dhananjay2021/ucp-go-sdk/models"
)

// WebhookVerifier verifies webhook signatures.
type WebhookVerifier struct {
	keys map[string]crypto.PublicKey
}

// NewWebhookVerifier creates a new webhook verifier from JWKs.
func NewWebhookVerifier(jwks []models.JWK) (*WebhookVerifier, error) {
	v := &WebhookVerifier{
		keys: make(map[string]crypto.PublicKey),
	}

	for _, jwk := range jwks {
		key, err := jwkToPublicKey(jwk)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JWK %s: %w", jwk.Kid, err)
		}
		v.keys[jwk.Kid] = key
	}

	return v, nil
}

// VerifyRequest verifies the signature of an HTTP request.
func (v *WebhookVerifier) VerifyRequest(r *http.Request, body []byte) error {
	// Get the signature header
	sig := r.Header.Get("X-Detached-JWT")
	if sig == "" {
		return errors.New("missing X-Detached-JWT header")
	}

	// Parse the detached JWS
	parts := strings.Split(sig, ".")
	if len(parts) != 3 {
		return errors.New("invalid JWS format")
	}

	// Decode header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("failed to decode JWS header: %w", err)
	}

	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return fmt.Errorf("failed to parse JWS header: %w", err)
	}

	// Get the signing key
	key, ok := v.keys[header.Kid]
	if !ok {
		return fmt.Errorf("unknown key ID: %s", header.Kid)
	}

	// For detached JWS, the payload is the request body
	payloadB64 := base64.RawURLEncoding.EncodeToString(body)

	// Reconstruct the signing input
	signingInput := parts[0] + "." + payloadB64

	// Decode signature
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Verify based on algorithm
	switch header.Alg {
	case "ES256":
		return verifyES256(key, signingInput, signature)
	case "RS256":
		return verifyRS256(key, signingInput, signature)
	default:
		return fmt.Errorf("unsupported algorithm: %s", header.Alg)
	}
}

// jwkToPublicKey converts a JWK to a crypto.PublicKey.
func jwkToPublicKey(jwk models.JWK) (crypto.PublicKey, error) {
	switch jwk.Kty {
	case "EC":
		return jwkToECDSAPublicKey(jwk)
	case "RSA":
		return jwkToRSAPublicKey(jwk)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}
}

func jwkToECDSAPublicKey(jwk models.JWK) (*ecdsa.PublicKey, error) {
	if jwk.X == "" || jwk.Y == "" {
		return nil, errors.New("missing EC key coordinates")
	}

	xBytes, err := base64.RawURLEncoding.DecodeString(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("failed to decode X coordinate: %w", err)
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(jwk.Y)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Y coordinate: %w", err)
	}

	curve, err := getCurve(jwk.Crv)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}, nil
}

func jwkToRSAPublicKey(jwk models.JWK) (*rsa.PublicKey, error) {
	if jwk.N == "" || jwk.E == "" {
		return nil, errors.New("missing RSA key components")
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	// Convert E to int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

func getCurve(name string) (elliptic.Curve, error) {
	switch name {
	case "P-256":
		return elliptic.P256(), nil
	case "P-384":
		return elliptic.P384(), nil
	case "P-521":
		return elliptic.P521(), nil
	default:
		return nil, fmt.Errorf("unsupported curve: %s", name)
	}
}

func verifyES256(key crypto.PublicKey, signingInput string, signature []byte) error {
	ecKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("invalid key type for ES256")
	}

	hash := sha256.Sum256([]byte(signingInput))

	// ES256 signature is R || S, each 32 bytes
	if len(signature) != 64 {
		return errors.New("invalid ES256 signature length")
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	if !ecdsa.Verify(ecKey, hash[:], r, s) {
		return errors.New("signature verification failed")
	}

	return nil
}

func verifyRS256(key crypto.PublicKey, signingInput string, signature []byte) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return errors.New("invalid key type for RS256")
	}

	hash := sha256.Sum256([]byte(signingInput))

	if err := rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, hash[:], signature); err != nil {
		return errors.New("signature verification failed")
	}

	return nil
}
