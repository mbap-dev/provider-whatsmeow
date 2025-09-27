package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"your.org/provider-whatsmeow/internal/amqp"
	"your.org/provider-whatsmeow/internal/config"
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
	// Session management endpoints
	router.HandleFunc("/sessions/{id}/connect", s.handleConnect).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/disconnect", s.handleDisconnect).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/reload", s.handleReload).Methods(http.MethodPost)
	router.HandleFunc("/sessions/{id}/qr", s.handleQR).Methods(http.MethodGet)
	router.HandleFunc("/sessions/{id}/status", s.handleStatus).Methods(http.MethodGet)
	// Debug: resolve destination JID applying overrides and LID lookup
	router.HandleFunc("/sessions/{id}/resolve", s.handleResolve).Methods(http.MethodGet)

	// Health and readiness probes
	router.HandleFunc("/healthz", s.handleHealth).Methods(http.MethodGet)
	router.HandleFunc("/readyz", s.handleReady).Methods(http.MethodGet)
	s.httpSrv = &http.Server{Addr: cfg.HTTPAddr, Handler: router}
	return s
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
	// UNO expects 204 No Content; connect proceeds asynchronously
	w.WriteHeader(http.StatusNoContent)
}

// handleDisconnect tears down an existing session.  It always
// returns HTTP 200.
func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = s.provider.Disconnect(id)
	// UNO expects 204 No Content
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
	// UNO expects 204 No Content
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
	// Provider stores QR as base64 PNG; UNO expects raw bytes image/png
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

// decodeJSON is a helper that decodes a JSON payload from the
// request body into the provided target.  It returns an error if
// decoding fails.  Not used currently but retained for future
// extensions (e.g. sending messages directly via HTTP).
func decodeJSON(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	return dec.Decode(target)
}
