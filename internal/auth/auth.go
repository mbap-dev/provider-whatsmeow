// Package auth implements password hashing, session management, and API key
// validation for the WhatsMeow dashboard.
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    uint32 = 2
	argon2Memory  uint32 = 64 * 1024 // 64 MB
	argon2Threads uint8  = 4
	argon2KeyLen  uint32 = 32
	saltLen              = 16
)

// HashPassword returns an argon2id-encoded hash of password.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory, argon2Time, argon2Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks password against an argon2id-encoded hash.
func VerifyPassword(password, encoded string) (bool, error) {
	salt, storedHash, t, m, p, err := decodeHash(encoded)
	if err != nil {
		return false, err
	}
	cmp := argon2.IDKey([]byte(password), salt, t, m, p, uint32(len(storedHash)))
	return subtle.ConstantTimeCompare(cmp, storedHash) == 1, nil
}

func decodeHash(encoded string) (salt, hash []byte, time, memory uint32, threads uint8, err error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil, nil, 0, 0, 0, errors.New("invalid hash format")
	}
	var m, t, p uint32
	if _, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("parse params: %w", err)
	}
	if salt, err = base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("decode salt: %w", err)
	}
	if hash, err = base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("decode hash: %w", err)
	}
	return salt, hash, t, m, uint8(p), nil
}
