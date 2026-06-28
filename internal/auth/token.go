package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	apiKeyPrefix = "wm_"
	tokenBytes   = 32
)

// generateToken returns a raw base64url token and its SHA-256 hash for storage.
func generateToken() (raw, hash string, err error) {
	b, err := randomBytes(tokenBytes)
	if err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, hashToken(raw), nil
}

// generateAPIKey returns a "wm_"-prefixed key and its SHA-256 hash.
func generateAPIKey() (raw, hash string, err error) {
	b, err := randomBytes(tokenBytes)
	if err != nil {
		return "", "", err
	}
	raw = apiKeyPrefix + base64.RawURLEncoding.EncodeToString(b)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}
	return b, nil
}
