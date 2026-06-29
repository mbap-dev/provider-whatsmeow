package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"your.org/provider-whatsmeow/internal/store/pg/pgdb"
)

const defaultSessionTTL = 7 * 24 * time.Hour

var (
	ErrInvalidUserInput = errors.New("invalid user input")
	ErrLastOwner        = errors.New("cannot remove the last owner")
)

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
		Email:        normalizeEmail(email),
		PasswordHash: hash,
		Role:         "owner",
	}); err != nil {
		return false, fmt.Errorf("create owner: %w", err)
	}
	return true, nil
}

// Authenticate verifies email+password and returns the user on success.
func (s *Store) Authenticate(ctx context.Context, email, password string) (*pgdb.User, error) {
	trimmedEmail := strings.TrimSpace(email)
	u, err := s.q.GetUserByEmail(ctx, normalizeEmail(trimmedEmail))
	if err != nil && trimmedEmail != normalizeEmail(trimmedEmail) {
		u, err = s.q.GetUserByEmail(ctx, trimmedEmail)
	}
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

// ListUsers returns all dashboard users.
func (s *Store) ListUsers(ctx context.Context) ([]pgdb.User, error) {
	return s.q.ListUsers(ctx)
}

// CreateUser creates a dashboard user with a hashed password.
func (s *Store) CreateUser(ctx context.Context, email, password, role string, assignedSessions []string) (*pgdb.User, error) {
	email = normalizeEmail(email)
	role = normalizeRole(role)
	if email == "" || password == "" || !validRole(role) {
		return nil, ErrInvalidUserInput
	}
	hash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u, err := s.q.CreateUser(ctx, pgdb.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.ReplaceUserSessionAccess(ctx, u.ID, assignedSessions); err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUser updates email/role and optionally replaces the password.
func (s *Store) UpdateUser(ctx context.Context, id pgtype.UUID, email, password, role string, assignedSessions []string) (*pgdb.User, error) {
	email = normalizeEmail(email)
	role = normalizeRole(role)
	if !id.Valid || email == "" || !validRole(role) {
		return nil, ErrInvalidUserInput
	}
	current, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if current.Role == "owner" && role != "owner" {
		if err := s.ensureAnotherOwner(ctx); err != nil {
			return nil, err
		}
	}
	var updated pgdb.User
	if strings.TrimSpace(password) == "" {
		updated, err = s.q.UpdateUser(ctx, pgdb.UpdateUserParams{
			ID:    id,
			Email: email,
			Role:  role,
		})
	} else {
		hash, hashErr := HashPassword(password)
		if hashErr != nil {
			return nil, fmt.Errorf("hash password: %w", hashErr)
		}
		updated, err = s.q.UpdateUserWithPassword(ctx, pgdb.UpdateUserWithPasswordParams{
			ID:           id,
			Email:        email,
			Role:         role,
			PasswordHash: hash,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if err := s.ReplaceUserSessionAccess(ctx, id, assignedSessions); err != nil {
		return nil, err
	}
	return &updated, nil
}

// DeleteUser removes a dashboard user and their sessions/API keys via cascades.
func (s *Store) DeleteUser(ctx context.Context, id pgtype.UUID) (*pgdb.User, error) {
	if !id.Valid {
		return nil, ErrInvalidUserInput
	}
	current, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if current.Role == "owner" {
		if err := s.ensureAnotherOwner(ctx); err != nil {
			return nil, err
		}
	}
	deleted, err := s.q.DeleteUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete user: %w", err)
	}
	return &deleted, nil
}

// ListUserSessionAccess returns session IDs explicitly linked to userID.
func (s *Store) ListUserSessionAccess(ctx context.Context, userID pgtype.UUID) ([]string, error) {
	if !userID.Valid {
		return nil, ErrInvalidUserInput
	}
	return s.q.ListUserSessionAccess(ctx, userID)
}

// ReplaceUserSessionAccess replaces all explicit session links for userID.
func (s *Store) ReplaceUserSessionAccess(ctx context.Context, userID pgtype.UUID, sessions []string) error {
	if !userID.Valid {
		return ErrInvalidUserInput
	}
	if err := s.q.DeleteUserSessionAccess(ctx, userID); err != nil {
		return fmt.Errorf("delete user session access: %w", err)
	}
	seen := map[string]struct{}{}
	for _, sessionID := range sessions {
		sessionID = normalizeSessionAccessID(sessionID)
		if sessionID == "" {
			continue
		}
		if _, ok := seen[sessionID]; ok {
			continue
		}
		seen[sessionID] = struct{}{}
		if err := s.q.InsertUserSessionAccess(ctx, pgdb.InsertUserSessionAccessParams{
			UserID:    userID,
			SessionID: sessionID,
		}); err != nil {
			return fmt.Errorf("insert user session access: %w", err)
		}
	}
	return nil
}

// UserCanAccessSession reports whether userID has an explicit link to sessionID.
func (s *Store) UserCanAccessSession(ctx context.Context, userID pgtype.UUID, sessionID string) (bool, error) {
	if !userID.Valid {
		return false, ErrInvalidUserInput
	}
	sessionID = normalizeSessionAccessID(sessionID)
	if sessionID == "" {
		return false, nil
	}
	return s.q.UserCanAccessSession(ctx, pgdb.UserCanAccessSessionParams{
		UserID:    userID,
		SessionID: sessionID,
	})
}

func (s *Store) ensureAnotherOwner(ctx context.Context) error {
	count, err := s.q.CountOwners(ctx)
	if err != nil {
		return fmt.Errorf("count owners: %w", err)
	}
	if count <= 1 {
		return ErrLastOwner
	}
	return nil
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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRole(role string) string {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		return "user"
	}
	return role
}

func validRole(role string) bool {
	return role == "owner" || role == "admin" || role == "user"
}

func normalizeSessionAccessID(sessionID string) string {
	return strings.TrimSpace(sessionID)
}
