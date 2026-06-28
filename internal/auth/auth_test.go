package auth

import (
	"strings"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "super-secret-123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatalf("unexpected hash format: %s", hash)
	}
	ok, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword: %v", err)
	}
	if !ok {
		t.Fatal("expected VerifyPassword to return true for correct password")
	}
	ok, err = VerifyPassword("wrong-password", hash)
	if err != nil {
		t.Fatalf("VerifyPassword wrong: %v", err)
	}
	if ok {
		t.Fatal("expected VerifyPassword to return false for wrong password")
	}
}

func TestHashPasswordUnique(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	if h1 == h2 {
		t.Fatal("two hashes of the same password must differ (different salts)")
	}
}

func TestGenerateToken(t *testing.T) {
	raw, hash, err := generateToken()
	if err != nil {
		t.Fatalf("generateToken: %v", err)
	}
	if raw == "" || hash == "" {
		t.Fatal("raw and hash must not be empty")
	}
	if raw == hash {
		t.Fatal("raw token and its hash must differ")
	}
	// Hash is deterministic.
	if hashToken(raw) != hash {
		t.Fatal("hashToken(raw) must equal the returned hash")
	}
}

func TestGenerateAPIKey(t *testing.T) {
	raw, hash, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey: %v", err)
	}
	if !strings.HasPrefix(raw, apiKeyPrefix) {
		t.Fatalf("api key must start with %q, got %q", apiKeyPrefix, raw)
	}
	if hashToken(raw) != hash {
		t.Fatal("hashToken(raw) must match returned hash")
	}
}
