package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"your.org/provider-whatsmeow/internal/store/pg/pgdb"
)

const defaultSessionTTL = 7 * 24 * time.Hour

// Store provides high-level auth operations backed by Postgres.
type Store struct {
	q *pgdb.Queries
}

// NewStore creates a Store using the given connection pool.
func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{q: pgdb.New(pool)}
}

// BootstrapOwner creates the first owner account when no users exist.
// Returns (true, nil) if the owner was created, (false, nil) if users already exist.
func (s *Store) BootstrapOwner(ctx context.Context, email, password string) (bool, error) {
	count, err := s.q.CountUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return false, nil
	}
	hash, err := HashPassword(password)
	if err != nil {
		return false, fmt.Errorf("hash password: %w", err)
	}
	if _, err = s.q.CreateUser(ctx, pgdb.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		Role:         "owner",
	}); err != nil {
		return false, fmt.Errorf("create owner: %w", err)
	}
	return true, nil
}

// Authenticate verifies email+password and returns the user on success.
func (s *Store) Authenticate(ctx context.Context, email, password string) (*pgdb.User, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		// Do not leak whether the email exists.
		return nil, fmt.Errorf("invalid credentials")
	}
	ok, err := VerifyPassword(password, u.PasswordHash)
	if err != nil || !ok {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &u, nil
}

// CreateSession generates a session token, persists its hash, and returns the raw token.
func (s *Store) CreateSession(ctx context.Context, userID pgtype.UUID, ttl time.Duration, ip, userAgent string) (string, error) {
	if ttl <= 0 {
		ttl = defaultSessionTTL
	}
	raw, hash, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	if _, err = s.q.CreateSession(ctx, pgdb.CreateSessionParams{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
		Ip:        pgtype.Text{String: ip, Valid: ip != ""},
		UserAgent: pgtype.Text{String: userAgent, Valid: userAgent != ""},
	}); err != nil {
		return "", fmt.Errorf("store session: %w", err)
	}
	return raw, nil
}

// ValidateSession looks up a session by raw token and returns the owner.
func (s *Store) ValidateSession(ctx context.Context, rawToken string) (*pgdb.User, error) {
	sess, err := s.q.GetSessionByTokenHash(ctx, hashToken(rawToken))
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}
	u, err := s.q.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &u, nil
}

// DeleteSession invalidates a session by its raw token.
func (s *Store) DeleteSession(ctx context.Context, rawToken string) error {
	return s.q.DeleteSession(ctx, hashToken(rawToken))
}

// PurgeExpiredSessions deletes all sessions whose expires_at is in the past.
func (s *Store) PurgeExpiredSessions(ctx context.Context) error {
	return s.q.DeleteExpiredSessions(ctx)
}

// CreateAPIKey generates a new API key, stores its hash, and returns the raw key + key ID.
func (s *Store) CreateAPIKey(ctx context.Context, userID pgtype.UUID, name string, expiresAt *time.Time) (rawKey string, id pgtype.UUID, err error) {
	raw, hash, err := generateAPIKey()
	if err != nil {
		return "", pgtype.UUID{}, fmt.Errorf("generate api key: %w", err)
	}
	var expTS pgtype.Timestamptz
	if expiresAt != nil {
		expTS = pgtype.Timestamptz{Time: *expiresAt, Valid: true}
	}
	k, err := s.q.CreateAPIKey(ctx, pgdb.CreateAPIKeyParams{
		UserID:    userID,
		KeyHash:   hash,
		Name:      name,
		ExpiresAt: expTS,
	})
	if err != nil {
		return "", pgtype.UUID{}, fmt.Errorf("store api key: %w", err)
	}
	return raw, k.ID, nil
}

// ValidateAPIKey resolves a raw key to its owner.
func (s *Store) ValidateAPIKey(ctx context.Context, rawKey string) (*pgdb.User, pgtype.UUID, error) {
	k, err := s.q.GetAPIKeyByHash(ctx, hashToken(rawKey))
	if err != nil {
		return nil, pgtype.UUID{}, fmt.Errorf("api key not found")
	}
	// Best-effort last_used_at update.
	_ = s.q.UpdateAPIKeyLastUsed(ctx, k.ID)
	u, err := s.q.GetUserByID(ctx, k.UserID)
	if err != nil {
		return nil, pgtype.UUID{}, fmt.Errorf("user not found")
	}
	return &u, k.ID, nil
}

// ListAPIKeys returns all API keys owned by userID.
func (s *Store) ListAPIKeys(ctx context.Context, userID pgtype.UUID) ([]pgdb.ApiKey, error) {
	return s.q.ListAPIKeysByUser(ctx, userID)
}

// DeleteAPIKey removes a key only if it belongs to userID.
func (s *Store) DeleteAPIKey(ctx context.Context, id, userID pgtype.UUID) error {
	return s.q.DeleteAPIKey(ctx, pgdb.DeleteAPIKeyParams{ID: id, UserID: userID})
}
