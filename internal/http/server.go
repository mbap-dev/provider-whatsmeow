package http

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"your.org/provider-whatsmeow/internal/amqp"
	"your.org/provider-whatsmeow/internal/config"
	"your.org/provider-whatsmeow/internal/dashboard"
	ilog "your.org/provider-whatsmeow/internal/log"
	"your.org/provider-whatsmeow/internal/provider"
)

// Server encapsulates the HTTP API surface exposed by the
// adapter.  It holds references to the configuration, the
// client manager and the AMQP consumer so that it can report
// readiness.  When Start is called the server begins listening on
// cfg.HTTPAddr.  Shutdown gracefully stops the listener.
type Server struct {
	cfg      *config.Config
	provider *provider.ClientManager
	consumer *amqp.Consumer
	assets   http.Handler
	httpSrv  *http.Server
	ready    bool
}

// NewServer constructs a new HTTP server.  It wires up all routes
// using Gorilla mux.  The server will report itself as ready once
// Start has been invoked.
func NewServer(cfg *config.Config, provider *provider.ClientManager, consumer *amqp.Consumer) *Server {
	s := &Server{
		cfg:      cfg,
		provider: provider,
		consumer: consumer,
	}
	router := mux.NewRouter()
	router.Use(s.uiAuthMiddleware)
	router.Use(s.apiTokenMiddleware)
	if cfg.UIEnabled {
		if assets, err := dashboard.AssetsHandler(); err == nil {
			s.assets = assets
			router.HandleFunc("/", s.handleDashboard).Methods(http.MethodGet)
			router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", s.assets))
		}
	}

	// Session management endpoints
	router.HandleFunc("/sessions/{id}/connect", s.handleConnect).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/disconnect", s.handleDisconnect).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/reload", s.handleReload).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/qr", s.handleQR).Methods(http.MethodGet)
	router.HandleFunc("/sessions/{id}/status", s.handleStatus).Methods(http.MethodGet)
	s.registerVersionedAPIRoutes(router, "/api/v1")
	s.registerVersionedAPIRoutes(router, "/api/{version:[^/]+}")
	// Meta-style compatibility: /vXX.X/{phone_number_id}/messages
	router.HandleFunc("/{version:[^/]+}/{id}/messages", s.handleSendMessage).Methods(http.MethodPost)
	router.HandleFunc("/{version:[^/]+}/{id}/messages/", s.handleSendMessage).Methods(http.MethodPost)
	router.HandleFunc("/{version:[^/]+}/{id}/media", s.handleUploadMedia).Methods(http.MethodPost)
	router.HandleFunc("/{version:[^/]+}/{id}/media/", s.handleUploadMedia).Methods(http.MethodPost)
	router.HandleFunc("/{version:[^/]+}/download/{id}/{file:.+}", s.handleDownloadMedia).Methods(http.MethodGet)
	router.HandleFunc("/{version:[^/]+}/download/{id}/{file:.+}/", s.handleDownloadMedia).Methods(http.MethodGet)
	router.HandleFunc("/{version:[^/]+}/{id}/{media_id}", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/{version:[^/]+}/{id}/{media_id}/", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/{version:[^/]+}/{media_id}", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/{version:[^/]+}/{media_id}/", s.handleGetMedia).Methods(http.MethodGet)
	// Compatibility when a proxy/base URL adds /api prefix.
	router.HandleFunc("/api/{version:[^/]+}/{id}/messages", s.handleSendMessage).Methods(http.MethodPost)
	router.HandleFunc("/api/{version:[^/]+}/{id}/messages/", s.handleSendMessage).Methods(http.MethodPost)
	router.HandleFunc("/api/{version:[^/]+}/{id}/media", s.handleUploadMedia).Methods(http.MethodPost)
	router.HandleFunc("/api/{version:[^/]+}/{id}/media/", s.handleUploadMedia).Methods(http.MethodPost)
	router.HandleFunc("/api/{version:[^/]+}/download/{id}/{file:.+}", s.handleDownloadMedia).Methods(http.MethodGet)
	router.HandleFunc("/api/{version:[^/]+}/download/{id}/{file:.+}/", s.handleDownloadMedia).Methods(http.MethodGet)
	router.HandleFunc("/api/{version:[^/]+}/{id}/{media_id}", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/api/{version:[^/]+}/{id}/{media_id}/", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/api/{version:[^/]+}/{media_id}", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc("/api/{version:[^/]+}/{media_id}/", s.handleGetMedia).Methods(http.MethodGet)
	// Calls (ring-only fake call)
	router.HandleFunc("/sessions/{id}/fakecall", s.handleFakeCall).Methods(http.MethodPost)
	// Debug: resolve destination JID applying overrides and LID lookup
	router.HandleFunc("/sessions/{id}/resolve", s.handleResolve).Methods(http.MethodGet)

	// Health and readiness probes
	router.HandleFunc("/healthz", s.handleHealth).Methods(http.MethodGet)
	router.HandleFunc("/readyz", s.handleReady).Methods(http.MethodGet)
	s.httpSrv = &http.Server{Addr: cfg.HTTPAddr, Handler: router}
	return s
}

func (s *Server) uiAuthMiddleware(next http.Handler) http.Handler {
	if s == nil || s.cfg == nil || !s.cfg.UIAuthEnabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !isUIProtectedPath(path) {
			next.ServeHTTP(w, r)
			return
		}
		user, pass, ok := r.BasicAuth()
		if !ok || !secureEqual(user, s.cfg.UIAuthUser) || !secureEqual(pass, s.cfg.UIAuthPass) {
			if isUITokenCompatiblePath(path) && s.hasValidAPIToken(r) {
				next.ServeHTTP(w, r)
				return
			}
			realm := strings.TrimSpace(s.cfg.UIAuthRealm)
			if realm == "" {
				realm = "Protected"
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="`+strings.ReplaceAll(realm, `"`, "")+`", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) apiTokenMiddleware(next http.Handler) http.Handler {
	if s == nil || s.cfg == nil || !s.cfg.APIAuthEnabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAPITokenProtectedPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		// UI users authenticated with Basic Auth can use the frontend without
		// manually injecting Bearer token in browser fetch calls.
		if s.hasValidUIBasicAuth(r) {
			next.ServeHTTP(w, r)
			return
		}
		if s.hasValidAPIToken(r) {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Bearer realm="WhatsMeow API"`)
		writeJSON(w, http.StatusUnauthorized, map[string]any{
			"error": map[string]any{
				"message": "Invalid OAuth access token.",
				"type":    "OAuthException",
				"code":    190,
			},
		})
	})
}

func (s *Server) hasValidUIBasicAuth(r *http.Request) bool {
	if s == nil || s.cfg == nil {
		return false
	}
	userCfg := strings.TrimSpace(s.cfg.UIAuthUser)
	passCfg := strings.TrimSpace(s.cfg.UIAuthPass)
	if userCfg == "" || passCfg == "" {
		return false
	}
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return secureEqual(user, userCfg) && secureEqual(pass, passCfg)
}

func (s *Server) hasValidAPIToken(r *http.Request) bool {
	if s == nil || s.cfg == nil {
		return false
	}
	expected := strings.TrimSpace(s.cfg.APIAuthToken)
	if expected == "" {
		return false
	}
	token := extractBearerToken(r.Header.Get("Authorization"))
	if token == "" && s.cfg.APIAuthAllowQuery {
		token = strings.TrimSpace(r.URL.Query().Get("access_token"))
	}
	if token == "" {
		return false
	}
	return secureEqual(token, expected)
}

func extractBearerToken(headerValue string) string {
	headerValue = strings.TrimSpace(headerValue)
	if headerValue == "" {
		return ""
	}
	parts := strings.Fields(headerValue)
	if len(parts) == 1 {
		// Accept raw token value for compatibility with some clients.
		return strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func isUIProtectedPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "/" {
		return true
	}
	if strings.HasPrefix(path, "/assets/") {
		return true
	}
	if strings.HasPrefix(path, "/api/") {
		return true
	}
	if strings.HasPrefix(path, "/sessions/") {
		return true
	}
	return false
}

func isAPITokenProtectedPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	switch path {
	case "/", "/healthz", "/readyz":
		return false
	}
	if strings.HasPrefix(path, "/assets/") {
		return false
	}
	return true
}

func isUITokenCompatiblePath(path string) bool {
	path = strings.TrimSpace(path)
	return strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/sessions/")
}

func secureEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	html, err := dashboard.Index()
	if err != nil {
		http.Error(w, "dashboard unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(html)
}

func (s *Server) registerVersionedAPIRoutes(router *mux.Router, base string) {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		return
	}
	router.HandleFunc(base+"/sessions/{id}/messages", s.handleSendMessage).Methods(http.MethodPost)
	router.HandleFunc(base+"/sessions/{id}/media", s.handleUploadMedia).Methods(http.MethodPost)
	router.HandleFunc(base+"/sessions/{id}/media/{media_id}", s.handleGetMedia).Methods(http.MethodGet)
	router.HandleFunc(base+"/download/{id}/{file:.+}", s.handleDownloadMedia).Methods(http.MethodGet)
	router.HandleFunc(base+"/sessions/{id}/config", s.handleGetSessionConfig).Methods(http.MethodGet)
	router.HandleFunc(base+"/sessions/{id}/config", s.handleUpdateSessionConfig).Methods(http.MethodPut)
	router.HandleFunc(base+"/sessions/{id}/config", s.handleClearSessionConfig).Methods(http.MethodDelete)
	router.HandleFunc(base+"/sessions/{id}/webhooks", s.handleListSessionWebhooks).Methods(http.MethodGet)
	router.HandleFunc(base+"/sessions/{id}/webhooks", s.handleCreateSessionWebhook).Methods(http.MethodPost)
	router.HandleFunc(base+"/sessions/{id}/webhooks/{webhook_id}", s.handleUpdateSessionWebhook).Methods(http.MethodPut)
	router.HandleFunc(base+"/sessions/{id}/webhooks/{webhook_id}", s.handleDeleteSessionWebhook).Methods(http.MethodDelete)
	router.HandleFunc(base+"/sessions", s.handleListSessions).Methods(http.MethodGet)
	router.HandleFunc(base+"/sessions", s.handleCreateSession).Methods(http.MethodPost)
	router.HandleFunc(base+"/worker", s.handleWorker).Methods(http.MethodGet)
	router.HandleFunc(base+"/worker/assignments", s.handleAssignSession).Methods(http.MethodPost)
	router.HandleFunc(base+"/worker/assignments/{id}", s.handleUnassignSession).Methods(http.MethodDelete)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	st, err := s.provider.Status(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(st)
}

func (s *Server) handleListSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"sessions": s.provider.ListSessions(),
	})
}

type createSessionRequest struct {
	ID      string `json:"id"`
	Connect *bool  `json:"connect,omitempty"`
	Assign  *bool  `json:"assign,omitempty"`
}

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	sessionID := strings.TrimSpace(req.ID)
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing id")
		return
	}
	shouldConnect := true
	if req.Connect != nil {
		shouldConnect = *req.Connect
	}
	shouldAssign := true
	if req.Assign != nil {
		shouldAssign = *req.Assign
	}
	if shouldAssign {
		if _, err := s.provider.AssignSession(sessionID); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if shouldConnect {
		if _, err := s.provider.Connect(sessionID); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"connected":  shouldConnect,
		"assigned":   shouldAssign,
	})
}

type assignSessionRequest struct {
	SessionID string `json:"session_id"`
	Connect   *bool  `json:"connect,omitempty"`
}

func (s *Server) handleAssignSession(w http.ResponseWriter, r *http.Request) {
	var req assignSessionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session_id")
		return
	}
	norm, err := s.provider.AssignSession(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	shouldConnect := false
	if req.Connect != nil {
		shouldConnect = *req.Connect
	}
	if shouldConnect {
		if _, err := s.provider.Connect(norm); err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": norm,
		"assigned":   true,
		"connected":  shouldConnect,
	})
}

func (s *Server) handleUnassignSession(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	removed, err := s.provider.UnassignSession(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"removed":    removed,
	})
}

type workerResponse struct {
	WorkerID          string   `json:"worker_id"`
	WorkerStrict      bool     `json:"worker_strict"`
	AssignedSessions  []string `json:"assigned_sessions"`
	AMQPEnabled       bool     `json:"amqp_enabled"`
	AMQPExchange      string   `json:"amqp_exchange"`
	AMQPQueue         string   `json:"amqp_queue"`
	AMQPBindings      []string `json:"amqp_bindings"`
	WorkerAutoConnect bool     `json:"worker_auto_connect_assigned"`
	UIEnabled         bool     `json:"ui_enabled"`
	UITitle           string   `json:"ui_title"`
}

func (s *Server) handleWorker(w http.ResponseWriter, r *http.Request) {
	bindings := s.cfg.AMQPBindings
	if len(bindings) == 0 && s.cfg.AMQPBinding != "" {
		bindings = []string{s.cfg.AMQPBinding}
	}
	resp := workerResponse{
		WorkerID:          s.provider.WorkerID(),
		WorkerStrict:      s.provider.WorkerStrict(),
		AssignedSessions:  s.provider.AssignedSessions(),
		AMQPEnabled:       s.cfg.AMQPConsumerEnabled,
		AMQPExchange:      s.cfg.AMQPExchange,
		AMQPQueue:         s.cfg.AMQPQueue,
		AMQPBindings:      bindings,
		WorkerAutoConnect: s.cfg.WorkerAutoConnectAssigned,
		UIEnabled:         s.cfg.UIEnabled,
		UITitle:           s.cfg.UITitle,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	to := strings.TrimSpace(r.URL.Query().Get("to"))
	if to == "" {
		http.Error(w, "missing to query param", http.StatusBadRequest)
		return
	}
	res, err := s.provider.ResolveDest(id, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["id"]
	var msg provider.OutgoingMessage
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	_ = r.Body.Close()
	if err := json.Unmarshal(rawBody, &msg); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	var raw map[string]any
	dec := json.NewDecoder(bytes.NewReader(rawBody))
	dec.UseNumber()
	_ = dec.Decode(&raw)
	if handled := s.handleCloudStatusUpdate(w, raw); handled {
		return
	}
	if strings.TrimSpace(msg.To) == "" {
		msg.To = extractSendRecipient(raw)
	}
	if strings.TrimSpace(msg.To) == "" {
		msg.To = strings.TrimSpace(r.URL.Query().Get("to"))
	}
	if strings.TrimSpace(msg.To) == "" {
		logSendRecipientDebug(sessionID, r, raw)
		writeJSONError(w, http.StatusBadRequest, "missing recipient: use 'to' (fallbacks: recipient_id/source_id)")
		return
	}
	if strings.TrimSpace(msg.SessionID) == "" {
		msg.SessionID = sessionID
	}
	msg.SessionID = strings.TrimSpace(msg.SessionID)
	if strings.TrimSpace(msg.Type) == "" && msg.Text != nil && strings.TrimSpace(msg.Text.Body) != "" {
		msg.Type = "text"
	}
	if strings.TrimSpace(msg.MessageID) == "" {
		msg.MessageID = generateAPIMessageID("INC")
	}
	if err := s.provider.Send(r.Context(), msg); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"messaging_product": "whatsapp",
		"contacts": []map[string]string{
			{
				"wa_id": normalizePhoneID(msg.To),
			},
		},
		"messages": []map[string]string{
			{
				"id": apiResponseMessageID(sessionID, msg.MessageID),
			},
		},
	})
}

func apiResponseMessageID(sessionID, messageID string) string {
	return provider.WithInstanceMessageID(sessionID, messageID)
}

func (s *Server) handleCloudStatusUpdate(w http.ResponseWriter, raw map[string]any) bool {
	status := strings.ToLower(strings.TrimSpace(anyStringValue(raw["status"])))
	if status == "" {
		return false
	}
	switch status {
	case "sent", "delivered", "read", "failed", "progress", "deleted":
		writeJSON(w, http.StatusOK, map[string]any{
			"success":    true,
			"status":     status,
			"message_id": strings.TrimSpace(anyStringValue(raw["message_id"])),
		})
		return true
	default:
		writeJSONError(w, http.StatusBadRequest, "unsupported status value")
		return true
	}
}

func (s *Server) handleUploadMedia(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		ilog.WithSession("").Error("evt=media_upload_api failed reason=missing_phone_number_id")
		writeJSONError(w, http.StatusBadRequest, "missing phone number id")
		return
	}

	const maxMediaUpload = 64 << 20 // 64 MiB
	if err := r.ParseMultipartForm(maxMediaUpload); err != nil {
		ilog.WithSession(sessionID).Error("evt=media_upload_api failed step=parse_multipart err=%v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid multipart payload")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		file, header, err = r.FormFile("media")
		if err != nil {
			ilog.WithSession(sessionID).Error("evt=media_upload_api failed reason=missing_file_field err=%v", err)
			writeJSONError(w, http.StatusBadRequest, "missing multipart file field 'file'")
			return
		}
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxMediaUpload+1))
	if err != nil {
		ilog.WithSession(sessionID).Error("evt=media_upload_api failed step=read_upload filename=%q err=%v", strings.TrimSpace(header.Filename), err)
		writeJSONError(w, http.StatusBadRequest, "failed to read upload")
		return
	}
	if len(data) == 0 {
		ilog.WithSession(sessionID).Error("evt=media_upload_api failed reason=empty_file filename=%q", strings.TrimSpace(header.Filename))
		writeJSONError(w, http.StatusBadRequest, "empty file")
		return
	}
	if len(data) > maxMediaUpload {
		ilog.WithSession(sessionID).Error(
			"evt=media_upload_api failed reason=file_too_large filename=%q bytes=%d limit=%d",
			strings.TrimSpace(header.Filename),
			len(data),
			maxMediaUpload,
		)
		writeJSONError(w, http.StatusRequestEntityTooLarge, "file too large")
		return
	}

	mimeType := strings.TrimSpace(header.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = strings.TrimSpace(r.FormValue("type"))
	}
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}
	media, err := s.provider.StoreUploadedMedia(r.Context(), sessionID, strings.TrimSpace(header.Filename), mimeType, data)
	if err != nil {
		ilog.WithSession(sessionID).Error(
			"evt=media_upload_api failed step=store_upload filename=%q mime=%s bytes=%d mode=%s err=%v",
			strings.TrimSpace(header.Filename),
			mimeType,
			len(data),
			s.provider.MediaMode(),
			err,
		)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ilog.WithSession(sessionID).WithMessageID(media.MediaID).Info(
		"evt=media_upload_api success filename=%q mime=%s bytes=%d public_id=%s mode=%s",
		strings.TrimSpace(header.Filename),
		mimeType,
		len(data),
		media.ID,
		s.provider.MediaMode(),
	)

	writeJSON(w, http.StatusOK, map[string]any{
		"id": media.MediaID,
	})
}

func (s *Server) handleGetMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := mediaRequestSessionID(r, vars)
	mediaID := strings.TrimSpace(vars["media_id"])
	version := strings.Trim(strings.TrimSpace(vars["version"]), "/")
	if sessionID == "" {
		ilog.WithSession("").Error("evt=media_get_api failed reason=missing_phone_number_id media_ref=%s", mediaID)
		writeJSONError(w, http.StatusBadRequest, "missing phone number id")
		return
	}
	if mediaID == "" {
		ilog.WithSession(sessionID).Error("evt=media_get_api failed reason=missing_media_id")
		writeJSONError(w, http.StatusBadRequest, "missing media id")
		return
	}

	meta, ok := s.provider.GetMediaMetadata(sessionID, mediaID)
	if !ok {
		ilog.WithSession(sessionID).WithMessageID(mediaID).Error("evt=media_get_api failed reason=metadata_not_found")
		writeJSONError(w, http.StatusNotFound, "media not found")
		return
	}
	baseURL := requestBaseURL(r)
	url := s.provider.MediaURL(r.Context(), baseURL, version, meta)
	if url == "" {
		ilog.WithSession(sessionID).WithMessageID(meta.MediaID).Error("evt=media_get_api failed reason=empty_media_url public_id=%s", meta.ID)
		writeJSONError(w, http.StatusInternalServerError, "failed to resolve media URL")
		return
	}
	ilog.WithSession(sessionID).WithMessageID(meta.MediaID).Info(
		"evt=media_get_api success requested_ref=%s public_id=%s mime=%s file_size=%d mode=%s",
		mediaID,
		meta.ID,
		meta.MIMEType,
		meta.FileSize,
		s.provider.MediaMode(),
	)

	resp := map[string]any{
		"messaging_product": "whatsapp",
		"url":               url,
		"mime_type":         meta.MIMEType,
		"sha256":            meta.SHA256,
		"file_size":         meta.FileSize,
		"id":                provider.CompatibleMediaID(meta),
	}
	if strings.TrimSpace(meta.MessagingProduct) != "" {
		resp["messaging_product"] = meta.MessagingProduct
	}
	if strings.TrimSpace(meta.Filename) != "" {
		resp["filename"] = meta.Filename
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleDownloadMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := strings.TrimSpace(vars["id"])
	filePath := strings.TrimSpace(vars["file"])
	if sessionID == "" || filePath == "" {
		ilog.WithSession(sessionID).Error("evt=media_download_api failed reason=missing_download_path file=%q", filePath)
		writeJSONError(w, http.StatusBadRequest, "missing download path")
		return
	}

	reader, meta, err := s.provider.OpenMediaForDownload(r.Context(), sessionID, filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ilog.WithSession(sessionID).WithMessageID(filePath).Error("evt=media_download_api failed reason=not_found file=%q", filePath)
			writeJSONError(w, http.StatusNotFound, "media not found")
			return
		}
		ilog.WithSession(sessionID).WithMessageID(filePath).Error("evt=media_download_api failed file=%q err=%v", filePath, err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	if strings.TrimSpace(meta.MIMEType) != "" {
		w.Header().Set("Content-Type", meta.MIMEType)
	}
	if strings.TrimSpace(meta.Filename) != "" {
		w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(meta.Filename, `"`, "")+`"`)
	}
	w.WriteHeader(http.StatusOK)
	ilog.WithSession(sessionID).WithMessageID(meta.MediaID).Info(
		"evt=media_download_api success file=%q public_id=%s mime=%s filename=%q mode=%s",
		filePath,
		meta.ID,
		meta.MIMEType,
		meta.Filename,
		s.provider.MediaMode(),
	)
	_, _ = io.Copy(w, reader)
}

type updateSessionConfigRequest struct {
	RejectCalls                 *bool   `json:"reject_calls,omitempty"`
	RejectCallsMessage          *string `json:"reject_calls_message,omitempty"`
	MarkReadOnMessage           *bool   `json:"mark_read_on_message,omitempty"`
	IgnoreStatusBroadcast       *bool   `json:"ignore_status_broadcast,omitempty"`
	IgnoreNewsletters           *bool   `json:"ignore_newsletters,omitempty"`
	OutgoingMessagesCoex        *bool   `json:"outgoing_messages_coex,omitempty"`
	AlwaysOnline                *bool   `json:"always_online,omitempty"`
	AlwaysOnlineIntervalSeconds *int    `json:"always_online_interval_seconds,omitempty"`
}

func (s *Server) handleGetSessionConfig(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	snapshot, err := s.provider.GetSessionConfig(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) handleUpdateSessionConfig(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	var req updateSessionConfigRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	snapshot, err := s.provider.UpdateSessionConfig(sessionID, provider.SessionConfigPatch{
		RejectCalls:                 req.RejectCalls,
		RejectCallsMessage:          req.RejectCallsMessage,
		MarkReadOnMessage:           req.MarkReadOnMessage,
		IgnoreStatusBroadcast:       req.IgnoreStatusBroadcast,
		IgnoreNewsletters:           req.IgnoreNewsletters,
		OutgoingMessagesCoex:        req.OutgoingMessagesCoex,
		AlwaysOnline:                req.AlwaysOnline,
		AlwaysOnlineIntervalSeconds: req.AlwaysOnlineIntervalSeconds,
	})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) handleClearSessionConfig(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	removed, snapshot, err := s.provider.ClearSessionConfig(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"removed":    removed,
		"config":     snapshot,
	})
}

type sessionWebhookRequest struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	URL                    string `json:"url,omitempty"`
	URLAbsolute            string `json:"url_absolute,omitempty"`
	Header                 string `json:"header,omitempty"`
	Token                  string `json:"token,omitempty"`
	TimeoutMS              *int   `json:"timeout_ms,omitempty"`
	Enabled                *bool  `json:"enabled,omitempty"`
	AppendSession          *bool  `json:"append_session,omitempty"`
	SendIncomingMessages   *bool  `json:"send_incoming_messages,omitempty"`
	SendOutgoingMessages   *bool  `json:"send_outgoing_messages,omitempty"`
	SendUpdateMessages     *bool  `json:"send_update_messages,omitempty"`
	SendGroupMessages      *bool  `json:"send_group_messages,omitempty"`
	SendNewsletterMessages *bool  `json:"send_newsletter_messages,omitempty"`
}

func (r sessionWebhookRequest) toProviderWebhook() provider.SessionWebhook {
	timeout := 5000
	if r.TimeoutMS != nil {
		timeout = *r.TimeoutMS
	}
	enabled := true
	if r.Enabled != nil {
		enabled = *r.Enabled
	}
	appendSession := true
	if r.AppendSession != nil {
		appendSession = *r.AppendSession
	}
	sendIncoming := true
	if r.SendIncomingMessages != nil {
		sendIncoming = *r.SendIncomingMessages
	}
	sendOutgoing := true
	if r.SendOutgoingMessages != nil {
		sendOutgoing = *r.SendOutgoingMessages
	}
	sendUpdate := true
	if r.SendUpdateMessages != nil {
		sendUpdate = *r.SendUpdateMessages
	}
	sendGroup := true
	if r.SendGroupMessages != nil {
		sendGroup = *r.SendGroupMessages
	}
	sendNewsletter := false
	if r.SendNewsletterMessages != nil {
		sendNewsletter = *r.SendNewsletterMessages
	}
	return provider.SessionWebhook{
		ID:                     strings.TrimSpace(r.ID),
		Name:                   strings.TrimSpace(r.Name),
		URL:                    strings.TrimSpace(r.URL),
		URLAbsolute:            strings.TrimSpace(r.URLAbsolute),
		Header:                 strings.TrimSpace(r.Header),
		Token:                  strings.TrimSpace(r.Token),
		TimeoutMS:              timeout,
		Enabled:                enabled,
		AppendSession:          appendSession,
		SendIncomingMessages:   sendIncoming,
		SendOutgoingMessages:   sendOutgoing,
		SendUpdateMessages:     sendUpdate,
		SendGroupMessages:      sendGroup,
		SendNewsletterMessages: sendNewsletter,
	}
}

func (s *Server) handleListSessionWebhooks(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	webhooks := s.provider.ListSessionWebhooks(sessionID)
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"webhooks":   webhooks,
	})
}

func (s *Server) handleCreateSessionWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	var req sessionWebhookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	webhook, err := s.provider.UpsertSessionWebhook(sessionID, req.toProviderWebhook())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"webhook":    webhook,
	})
}

func (s *Server) handleUpdateSessionWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	webhookID := strings.TrimSpace(mux.Vars(r)["webhook_id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	if webhookID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing webhook id")
		return
	}
	var req sessionWebhookRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	if req.ID != "" && strings.TrimSpace(req.ID) != webhookID {
		writeJSONError(w, http.StatusBadRequest, "webhook id mismatch")
		return
	}
	req.ID = webhookID
	webhook, err := s.provider.UpsertSessionWebhook(sessionID, req.toProviderWebhook())
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"webhook":    webhook,
	})
}

func (s *Server) handleDeleteSessionWebhook(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(mux.Vars(r)["id"])
	webhookID := strings.TrimSpace(mux.Vars(r)["webhook_id"])
	if sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing session id")
		return
	}
	if webhookID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing webhook id")
		return
	}
	removed, err := s.provider.DeleteSessionWebhook(sessionID, webhookID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"session_id": sessionID,
		"webhook_id": webhookID,
		"removed":    removed,
	})
}

type fakeCallRequest struct {
	To     string `json:"to"`
	RingMS int    `json:"ring_ms"`
}

// handleFakeCall triggers a ring-only fake voice call to the destination.
// Body: {"to": "<phone or jid>", "ring_ms": 15000}
func (s *Server) handleFakeCall(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req fakeCallRequest
	if err := decodeJSON(r, &req); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}
	to := strings.TrimSpace(req.To)
	if to == "" {
		http.Error(w, "missing to", http.StatusBadRequest)
		return
	}
	ringMS := req.RingMS
	if ringMS <= 0 {
		ringMS = 15000
	}

	callID, err := s.provider.FakeCall(id, to, ringMS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"call_id": callID,
		"ring_ms": ringMS,
		"to":      to,
	})
}

// Start begins serving HTTP requests.  It sets the readiness flag
// which causes the readyz endpoint to return HTTP 200.  If the
// underlying http.Server exits with an error other than
// http.ErrServerClosed it will be returned to the caller.
func (s *Server) Start() error {
	s.ready = true
	return s.httpSrv.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.  After shutdown the
// readyz endpoint will return HTTP 503.
func (s *Server) Shutdown(ctx context.Context) error {
	s.ready = false
	return s.httpSrv.Shutdown(ctx)
}

// handleConnect initialises or retrieves a session.  It always
// returns HTTP 200, optionally with an error message in the body.
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := s.provider.Connect(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, err.Error())
		return
	}
	// Connect proceeds asynchronously
	w.WriteHeader(http.StatusNoContent)
}

// handleDisconnect tears down an existing session.  It always
// returns HTTP 200.
func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = s.provider.Disconnect(id)
	// Keep response empty
	w.WriteHeader(http.StatusNoContent)
}

// handleReload triggers a session reload.  It returns 200 on success
// or 500 if reload is not implemented or fails.
func (s *Server) handleReload(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := s.provider.Reload(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, err.Error())
		return
	}
	// Keep response empty
	w.WriteHeader(http.StatusNoContent)
}

// handleQR returns a base64 encoded PNG image of the current QR code
// for the session.  On success the response has content type
// text/plain and the body contains only the base64 string.  If the
// session does not exist or QR generation fails then an error
// message is returned with status 404 or 500 respectively.
func (s *Server) handleQR(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	qr, err := s.provider.GetQR(id)
	if err != nil {
		if strings.Contains(err.Error(), "not connected") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, _ = io.WriteString(w, err.Error())
		return
	}
	// Provider stores QR as base64 PNG and returns raw bytes image/png
	buf, decErr := base64.StdEncoding.DecodeString(qr)
	if decErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, decErr.Error())
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf)
}

// handleHealth always returns HTTP 200 OK.  It can be used by
// Kubernetes or other orchestrators to determine if the process is
// alive.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleReady returns HTTP 200 if the service has started and the
// consumer has bound the queue.  Otherwise it returns 503 Service
// Unavailable.  Note that the skeleton does not track the consumer
// status so readiness simply reflects whether the HTTP server has
// started.
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if s.ready {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ready")
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if xfp := firstHeaderToken(r.Header.Get("X-Forwarded-Proto")); xfp != "" {
		scheme = xfp
	}
	host := strings.TrimSpace(r.Host)
	if xfh := firstHeaderToken(r.Header.Get("X-Forwarded-Host")); xfh != "" {
		host = xfh
	}
	if host == "" {
		host = "localhost"
	}
	return scheme + "://" + host
}

func mediaRequestSessionID(r *http.Request, vars map[string]string) string {
	if sessionID := strings.TrimSpace(vars["id"]); sessionID != "" {
		return sessionID
	}
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.URL.Query().Get("phone_number_id"))
}

func firstHeaderToken(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if i := strings.IndexByte(v, ','); i >= 0 {
		v = v[:i]
	}
	return strings.TrimSpace(v)
}

// decodeJSON is a helper that decodes a JSON payload from the
// request body into the provided target.  It returns an error if
// decoding fails.  Not used currently but retained for future
// extensions (e.g. sending messages directly via HTTP).
func decodeJSON(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	return dec.Decode(target)
}

func extractSendRecipient(raw map[string]any) string {
	if len(raw) == 0 {
		return ""
	}
	candidates := []string{
		anyStringValue(raw["to"]),
		anyStringValue(raw["recipient_id"]),
		anyStringValue(raw["source_id"]),
		anyStringValue(raw["wa_id"]),
		anyStringValue(raw["phone"]),
		anyStringValue(raw["number"]),
		anyStringValue(raw["recipient"]),
	}
	for _, c := range candidates {
		if c != "" {
			return c
		}
	}

	if conversation, ok := raw["conversation"].(map[string]any); ok {
		if v := anyStringValue(conversation["source_id"]); v != "" {
			return v
		}
		if inbox, ok := conversation["contact_inbox"].(map[string]any); ok {
			if v := anyStringValue(inbox["source_id"]); v != "" {
				return v
			}
		}
	}
	if attrs, ok := raw["content_attributes"].(map[string]any); ok {
		if v := anyStringValue(attrs["source_id"]); v != "" {
			return v
		}
		if v := anyStringValue(attrs["recipient_id"]); v != "" {
			return v
		}
	}

	if msgs, ok := raw["messages"].([]any); ok && len(msgs) > 0 {
		if first, ok := msgs[0].(map[string]any); ok {
			if v := anyStringValue(first["to"]); v != "" {
				return v
			}
			if v := anyStringValue(first["recipient_id"]); v != "" {
				return v
			}
		}
	}

	if contacts, ok := raw["contacts"].([]any); ok && len(contacts) > 0 {
		if first, ok := contacts[0].(map[string]any); ok {
			if v := anyStringValue(first["wa_id"]); v != "" {
				return v
			}
		}
	}

	entryAny, ok := raw["entry"].([]any)
	if !ok || len(entryAny) == 0 {
		return ""
	}
	entry, ok := entryAny[0].(map[string]any)
	if !ok {
		return ""
	}
	changesAny, ok := entry["changes"].([]any)
	if !ok || len(changesAny) == 0 {
		return ""
	}
	change, ok := changesAny[0].(map[string]any)
	if !ok {
		return ""
	}
	value, ok := change["value"].(map[string]any)
	if !ok {
		return ""
	}
	if statuses, ok := value["statuses"].([]any); ok && len(statuses) > 0 {
		if st, ok := statuses[0].(map[string]any); ok {
			if v := anyStringValue(st["recipient_id"]); v != "" {
				return v
			}
		}
	}
	if messages, ok := value["messages"].([]any); ok && len(messages) > 0 {
		if first, ok := messages[0].(map[string]any); ok {
			if v := anyStringValue(first["to"]); v != "" {
				return v
			}
		}
	}
	if contacts, ok := value["contacts"].([]any); ok && len(contacts) > 0 {
		if first, ok := contacts[0].(map[string]any); ok {
			if v := anyStringValue(first["wa_id"]); v != "" {
				return v
			}
		}
	}
	return ""
}

func anyStringValue(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case json.Number:
		return strings.TrimSpace(t.String())
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strings.TrimSpace(strconv.FormatFloat(t, 'f', -1, 64))
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	default:
		return ""
	}
}

func logSendRecipientDebug(sessionID string, r *http.Request, raw map[string]any) {
	entry := ilog.WithSession(sessionID)
	contentType := strings.TrimSpace(r.Header.Get("Content-Type"))
	keys := payloadKeys(raw)
	recipientID := anyStringValue(raw["recipient_id"])
	sourceID := anyStringValue(raw["source_id"])
	to := anyStringValue(raw["to"])
	conversationSourceID := ""
	if conv, ok := raw["conversation"].(map[string]any); ok {
		conversationSourceID = anyStringValue(conv["source_id"])
		if conversationSourceID == "" {
			if inbox, ok := conv["contact_inbox"].(map[string]any); ok {
				conversationSourceID = anyStringValue(inbox["source_id"])
			}
		}
	}
	entry.Error(
		"send_http missing recipient content_type=%q keys=%v to=%q recipient_id=%q source_id=%q conversation_source_id=%q",
		contentType,
		keys,
		to,
		recipientID,
		sourceID,
		conversationSourceID,
	)
}

func payloadKeys(raw map[string]any) []string {
	if len(raw) == 0 {
		return nil
	}
	out := make([]string, 0, len(raw))
	for k := range raw {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"status":  "error",
		"message": strings.TrimSpace(message),
	})
}

func normalizePhoneID(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range v {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	if b.Len() > 0 {
		return b.String()
	}
	return v
}

func generateAPIMessageID(prefix string) string {
	prefix = strings.TrimSpace(strings.ToUpper(prefix))
	if prefix == "" {
		prefix = "INC"
	}
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return prefix + "-" + strings.ToUpper(hex.EncodeToString([]byte(time.Now().Format("150405.000000000"))))
	}
	return prefix + "-" + strings.ToUpper(hex.EncodeToString(buf))
}
