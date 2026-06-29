package auth

import (
	"context"
	"net/http"
	"strings"

	"your.org/provider-whatsmeow/internal/store/pg/pgdb"
)

type contextKey int

const userCtxKey contextKey = 0

const sessionCookieName = "wm_session"

// UserFromContext retrieves the authenticated user injected by Middleware.
func UserFromContext(ctx context.Context) (*pgdb.User, bool) {
	u, ok := ctx.Value(userCtxKey).(*pgdb.User)
	return u, ok && u != nil
}

func withUser(ctx context.Context, u *pgdb.User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// Middleware populates the user in context when a valid session cookie or
// "wm_"-prefixed Bearer token is present. Unauthenticated requests pass through.
func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if u := h.extractUser(r); u != nil {
			ctx = withUser(ctx, u)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth wraps a handler and returns 401 when no authenticated user is in context.
func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := UserFromContext(r.Context()); !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireOwner wraps a handler and returns 403 unless the user owns auth settings.
func (h *Handler) RequireOwner(next http.Handler) http.Handler {
	return h.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := UserFromContext(r.Context())
		if u.Role != "owner" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "owner role required"})
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// RequireAPIKeyManager allows users that can create machine-to-machine API keys.
func (h *Handler) RequireAPIKeyManager(next http.Handler) http.Handler {
	return h.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := UserFromContext(r.Context())
		if u.Role != "owner" && u.Role != "admin" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "owner or admin role required"})
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// UserFromRequest extracts an authenticated user from a request via cookie or Bearer wm_* key.
// Returns nil when no valid credentials are found.
func (h *Handler) UserFromRequest(r *http.Request) *pgdb.User {
	return h.extractUser(r)
}

// UserCanAccessSession reports whether u may access sessionID.
func (h *Handler) UserCanAccessSession(ctx context.Context, u *pgdb.User, sessionID string) bool {
	if h == nil || h.store == nil || u == nil {
		return true
	}
	if u.Role == "owner" {
		return true
	}
	ok, err := h.store.UserCanAccessSession(ctx, u.ID, sessionID)
	return err == nil && ok
}

// ListUserSessionAccess returns the sessions explicitly linked to u.
func (h *Handler) ListUserSessionAccess(ctx context.Context, u *pgdb.User) []string {
	if h == nil || h.store == nil || u == nil {
		return nil
	}
	items, err := h.store.ListUserSessionAccess(ctx, u.ID)
	if err != nil {
		return nil
	}
	return items
}

func (h *Handler) extractUser(r *http.Request) *pgdb.User {
	ctx := r.Context()
	// Cookie session takes priority.
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		if u, err := h.store.ValidateSession(ctx, c.Value); err == nil {
			return u
		}
	}
	// Bearer API key — only accept "wm_" prefix to avoid confusion with the
	// existing static bearer token (APIAuthToken) which has no such prefix.
	if bearer := bearerToken(r); strings.HasPrefix(bearer, apiKeyPrefix) {
		if u, _, err := h.store.ValidateAPIKey(ctx, bearer); err == nil {
			return u
		}
	}
	return nil
}

func bearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
