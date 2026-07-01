package provider

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types"
	log "your.org/provider-whatsmeow/internal/log"
)

func (m *ClientManager) setPasskeyRequest(sessionID string, publicKey *types.WebAuthnPublicKey) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return
	}
	requestOptionsJSON := ""
	if publicKey != nil {
		if data, err := json.Marshal(publicKey); err == nil {
			requestOptionsJSON = string(data)
		}
	}
	m.passkeyMu.Lock()
	m.passkeyRequests[sessionID] = PasskeyRequest{
		Required:           true,
		RequestOptionsJSON: strings.TrimSpace(requestOptionsJSON),
		UpdatedAt:          time.Now(),
	}
	m.passkeyMu.Unlock()
}

func (m *ClientManager) setPasskeyConfirmation(sessionID, code string, skipHandoffUX bool) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return
	}
	m.passkeyMu.Lock()
	req := m.passkeyRequests[sessionID]
	req.Required = true
	req.ConfirmationRequired = !skipHandoffUX
	req.ConfirmationCode = strings.TrimSpace(code)
	req.SkipHandoffUX = skipHandoffUX
	req.UpdatedAt = time.Now()
	m.passkeyRequests[sessionID] = req
	m.passkeyMu.Unlock()
}

func (m *ClientManager) clearPasskeyRequest(sessionID string) {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return
	}
	m.passkeyMu.Lock()
	delete(m.passkeyRequests, sessionID)
	m.passkeyMu.Unlock()
}

func (m *ClientManager) PasskeyRequest(sessionID string) PasskeyRequest {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return PasskeyRequest{}
	}
	m.passkeyMu.RLock()
	req := m.passkeyRequests[sessionID]
	m.passkeyMu.RUnlock()
	return req
}

type SubmitPasskeyPrologueRequest struct {
	CredentialID        string          `json:"credential_id"`
	WebAuthnAssertion   json.RawMessage `json:"webauthn_assertion"`
	PairingHandoffProof []byte          `json:"pairing_handoff_proof,omitempty"`
}

func (m *ClientManager) SubmitPasskeyPrologue(ctx context.Context, sessionID string, req SubmitPasskeyPrologueRequest) error {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	req.CredentialID = strings.TrimSpace(req.CredentialID)
	if len(req.WebAuthnAssertion) == 0 || !json.Valid(req.WebAuthnAssertion) {
		return ErrMissingPasskey
	}
	var assertion types.WebAuthnResponse
	if err := json.Unmarshal(req.WebAuthnAssertion, &assertion); err != nil {
		return ErrMissingPasskey
	}
	if assertion.ID == "" || len(assertion.RawID) == 0 {
		return ErrMissingPasskey
	}
	client, ok := m.GetClient(sessionID)
	if !ok || client == nil {
		return ErrPairingNotReady
	}
	log.WithSession(sessionID).Info("evt=passkey_response_submit credential_id_len=%d raw_id_len=%d assertion_len=%d", len(assertion.ID), len(assertion.RawID), len(req.WebAuthnAssertion))
	err := client.SendPasskeyResponse(ctx, &assertion)
	if err != nil {
		log.WithSession(sessionID).Error("evt=passkey_response_submit_failed err=%v", err)
		return err
	}
	log.WithSession(sessionID).Info("evt=passkey_response_submit_ok")
	return nil
}

func (m *ClientManager) ConfirmPasskey(ctx context.Context, sessionID string) error {
	sessionID = normalizeSessionID(sessionID)
	if sessionID == "" {
		return ErrInvalidSessionID
	}
	client, ok := m.GetClient(sessionID)
	if !ok || client == nil {
		return ErrPairingNotReady
	}
	log.WithSession(sessionID).Info("evt=passkey_confirmation_submit")
	if err := client.SendPasskeyConfirmation(ctx); err != nil {
		log.WithSession(sessionID).Error("evt=passkey_confirmation_submit_failed err=%v", err)
		return err
	}
	log.WithSession(sessionID).Info("evt=passkey_confirmation_submit_ok")
	return nil
}
