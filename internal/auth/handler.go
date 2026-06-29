package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"your.org/provider-whatsmeow/internal/store/pg/pgdb"
)

// Handler holds dependencies for auth HTTP endpoints.
type Handler struct {
	store        *Store
	secureCookie bool
	sessionTTL   time.Duration
}

// NewHandler creates an auth Handler.
func NewHandler(store *Store, secureCookie bool, sessionTTL time.Duration) *Handler {
	if sessionTTL <= 0 {
		sessionTTL = defaultSessionTTL
	}
	return &Handler{store: store, secureCookie: secureCookie, sessionTTL: sessionTTL}
}

// Register attaches auth routes to r.
// /auth/login and /auth/logout are public; all others require a valid session.
func (h *Handler) Register(r *mux.Router) {
	sub := r.PathPrefix("/auth").Subrouter()
	sub.Use(h.Middleware)

	sub.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	sub.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	sub.Handle("/me", h.RequireAuth(http.HandlerFunc(h.Me))).Methods(http.MethodGet)
	sub.Handle("/api-keys", h.RequireAuth(http.HandlerFunc(h.ListAPIKeys))).Methods(http.MethodGet)
	sub.Handle("/api-keys", h.RequireAPIKeyManager(http.HandlerFunc(h.CreateAPIKey))).Methods(http.MethodPost)
	sub.Handle("/api-keys/{id}", h.RequireAuth(http.HandlerFunc(h.DeleteAPIKey))).Methods(http.MethodDelete)
	sub.Handle("/users", h.RequireOwner(http.HandlerFunc(h.ListUsers))).Methods(http.MethodGet)
	sub.Handle("/users", h.RequireOwner(http.HandlerFunc(h.CreateUser))).Methods(http.MethodPost)
	sub.Handle("/users/{id}", h.RequireOwner(http.HandlerFunc(h.UpdateUser))).Methods(http.MethodPut)
	sub.Handle("/users/{id}", h.RequireOwner(http.HandlerFunc(h.DeleteUser))).Methods(http.MethodDelete)
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	u, err := h.store.Authenticate(r.Context(), body.Email, body.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = fwd
	}
	rawToken, err := h.store.CreateSession(r.Context(), u.ID, h.sessionTTL, ip, r.UserAgent())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create session"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    rawToken,
		Path:     "/",
		Expires:  time.Now().Add(h.sessionTTL),
		HttpOnly: true,
		Secure:   h.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, h.userResponse(r.Context(), u))
}

// Logout handles POST /auth/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		_ = h.store.DeleteSession(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})
	writeJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// Me handles GET /auth/me.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	u, _ := UserFromContext(r.Context())
	writeJSON(w, http.StatusOK, h.userResponse(r.Context(), u))
}

// ListAPIKeys handles GET /auth/api-keys.
func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	u, _ := UserFromContext(r.Context())
	keys, err := h.store.ListAPIKeys(r.Context(), u.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not list api keys"})
		return
	}
	out := make([]map[string]any, len(keys))
	for i := range keys {
		out[i] = apiKeyResponse(&keys[i], "")
	}
	writeJSON(w, http.StatusOK, out)
}

// CreateAPIKey handles POST /auth/api-keys.
func (h *Handler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name      string  `json:"name"`
		ExpiresAt *string `json:"expires_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	var expiresAt *time.Time
	if body.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *body.ExpiresAt)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "expires_at must be RFC3339"})
			return
		}
		expiresAt = &t
	}
	u, _ := UserFromContext(r.Context())
	rawKey, id, err := h.store.CreateAPIKey(r.Context(), u.ID, body.Name, expiresAt)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create api key"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":   uuidStr(id),
		"name": body.Name,
		"key":  rawKey,
	})
}

// DeleteAPIKey handles DELETE /auth/api-keys/{id}.
func (h *Handler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	var keyID pgtype.UUID
	if err := keyID.Scan(mux.Vars(r)["id"]); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid key id"})
		return
	}
	u, _ := UserFromContext(r.Context())
	if err := h.store.DeleteAPIKey(r.Context(), keyID, u.ID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "api key not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers handles GET /auth/users.
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListUsers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not list users"})
		return
	}
	out := make([]map[string]any, len(users))
	for i := range users {
		out[i] = h.userResponse(r.Context(), &users[i])
	}
	writeJSON(w, http.StatusOK, out)
}

// CreateUser handles POST /auth/users.
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email            string   `json:"email"`
		Password         string   `json:"password"`
		Role             string   `json:"role"`
		AssignedSessions []string `json:"assigned_sessions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	u, err := h.store.CreateUser(r.Context(), body.Email, body.Password, body.Role, body.AssignedSessions)
	if err != nil {
		writeUserMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, h.userResponse(r.Context(), u))
}

// UpdateUser handles PUT /auth/users/{id}.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUIDParam(w, mux.Vars(r)["id"], "invalid user id")
	if !ok {
		return
	}
	var body struct {
		Email            string   `json:"email"`
		Password         string   `json:"password"`
		Role             string   `json:"role"`
		AssignedSessions []string `json:"assigned_sessions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	u, err := h.store.UpdateUser(r.Context(), userID, body.Email, body.Password, body.Role, body.AssignedSessions)
	if err != nil {
		writeUserMutationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h.userResponse(r.Context(), u))
}

// DeleteUser handles DELETE /auth/users/{id}.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUIDParam(w, mux.Vars(r)["id"], "invalid user id")
	if !ok {
		return
	}
	current, _ := UserFromContext(r.Context())
	if current != nil && current.ID == userID {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot delete current user"})
		return
	}
	if _, err := h.store.DeleteUser(r.Context(), userID); err != nil {
		writeUserMutationError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- helpers -----------------------------------------------------------------

func (h *Handler) userResponse(ctx context.Context, u *pgdb.User) map[string]any {
	m := map[string]any{
		"id":         uuidStr(u.ID),
		"email":      u.Email,
		"role":       u.Role,
		"created_at": u.CreatedAt.Time,
	}
	if h != nil && h.store != nil {
		sessions, err := h.store.ListUserSessionAccess(ctx, u.ID)
		if err == nil {
			m["assigned_sessions"] = sessions
		}
	}
	return m
}

func apiKeyResponse(k *pgdb.ApiKey, rawKey string) map[string]any {
	m := map[string]any{
		"id":         uuidStr(k.ID),
		"name":       k.Name,
		"created_at": k.CreatedAt.Time,
	}
	if k.LastUsedAt.Valid {
		m["last_used_at"] = k.LastUsedAt.Time
	}
	if k.ExpiresAt.Valid {
		m["expires_at"] = k.ExpiresAt.Time
	}
	if rawKey != "" {
		m["key"] = rawKey
	}
	return m
}

func parseUUIDParam(w http.ResponseWriter, value, message string) (pgtype.UUID, bool) {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": message})
		return id, false
	}
	return id, true
}

func writeUserMutationError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError
	switch {
	case errors.Is(err, ErrInvalidUserInput):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email, role and password for new users are required"})
	case errors.Is(err, ErrLastOwner):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot remove the last owner"})
	case errors.As(err, &pgErr) && pgErr.Code == "23505":
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email already exists"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save user"})
	}
}

func uuidStr(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	b := u.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
