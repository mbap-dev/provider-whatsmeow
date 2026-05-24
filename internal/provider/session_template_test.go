package provider

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestSessionTemplateRenderSupportsNamedAndPositionalVariables(t *testing.T) {
	manager := NewClientManager(ManagerOptions{
		SessionStore:        t.TempDir(),
		SessionTemplateFile: filepath.Join(t.TempDir(), "templates.json"),
	})
	_, err := manager.UpsertSessionTemplate("session-a", SessionTemplate{
		Name:     "boleto_disponivel",
		Language: "pt_BR",
		Header: SessionTemplatePart{
			Type: "text",
			Text: "Olá {{nome}}",
		},
		Body:   "Seu boleto {{1}} vence em {{vencimento}}.",
		Footer: "Protocolo {{protocolo}}",
	})
	if err != nil {
		t.Fatalf("UpsertSessionTemplate() error = %v", err)
	}

	rendered, err := manager.renderTemplateMessage("session-a", &TemplateContent{
		Name:     "boleto_disponivel",
		Language: TemplateLanguage{Code: "pt_BR"},
		Variables: map[string]string{
			"nome":       "Maria",
			"vencimento": "10/06/2026",
			"protocolo":  "ABC-123",
		},
		Components: []TemplateComponent{
			{
				Type: "body",
				Parameters: []TemplateParameter{
					{Type: "text", Text: "12345"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("renderTemplateMessage() error = %v", err)
	}
	if got := rendered.HeaderText; got != "Olá Maria" {
		t.Fatalf("header = %q", got)
	}
	if got := rendered.Body; got != "Seu boleto 12345 vence em 10/06/2026." {
		t.Fatalf("body = %q", got)
	}
	if got := rendered.Footer; got != "Protocolo ABC-123" {
		t.Fatalf("footer = %q", got)
	}
}

func TestSessionTemplateRenderUsesCloudAPIParameterName(t *testing.T) {
	manager := NewClientManager(ManagerOptions{SessionStore: t.TempDir()})
	_, err := manager.UpsertSessionTemplate("session-a", SessionTemplate{
		Name: "aviso",
		Body: "Olá {{cliente}}, status {{status}}.",
	})
	if err != nil {
		t.Fatalf("UpsertSessionTemplate() error = %v", err)
	}

	rendered, err := manager.renderTemplateMessage("session-a", &TemplateContent{
		Name: "aviso",
		Components: []TemplateComponent{
			{
				Type: "body",
				Parameters: []TemplateParameter{
					{Type: "text", ParameterName: "cliente", Text: "João"},
					{Type: "text", ParameterName: "status", Text: "aprovado"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("renderTemplateMessage() error = %v", err)
	}
	if got := rendered.Body; got != "Olá João, status aprovado." {
		t.Fatalf("body = %q", got)
	}
}

func TestBuildStructuredTemplateMessageKeepsHeaderBodyAndFooter(t *testing.T) {
	manager := NewClientManager(ManagerOptions{SessionStore: t.TempDir()})
	wire, err := manager.buildStructuredTemplateMessage(nil, nil, "session-a", &renderedTemplate{
		TemplateID: "aviso__pt_BR",
		HeaderType: "text",
		HeaderText: "Olá Maria",
		Body:       "Seu boleto está disponível.",
		Footer:     "Protocolo ABC-123",
	}, nil)
	if err != nil {
		t.Fatalf("buildStructuredTemplateMessage() error = %v", err)
	}
	interactive := wire.GetInteractiveMessage()
	if interactive == nil {
		t.Fatal("expected interactive message")
	}
	if got := interactive.GetHeader().GetTitle(); got != "Olá Maria" {
		t.Fatalf("header = %q", got)
	}
	if got := interactive.GetBody().GetText(); got != "Seu boleto está disponível." {
		t.Fatalf("body = %q", got)
	}
	if got := interactive.GetFooter().GetText(); got != "Protocolo ABC-123" {
		t.Fatalf("footer = %q", got)
	}
	if got := interactive.GetHeader().GetHasMediaAttachment(); got {
		t.Fatalf("header media attachment = %t", got)
	}
	native := interactive.GetNativeFlowMessage()
	if native == nil {
		t.Fatal("expected native flow message")
	}
	if got := native.GetMessageVersion(); got != 1 {
		t.Fatalf("message version = %d", got)
	}
}

func TestTemplateOrderDetailsUsesCloudAPIActionPayload(t *testing.T) {
	manager := NewClientManager(ManagerOptions{SessionStore: t.TempDir()})
	_, err := manager.UpsertSessionTemplate("session-a", SessionTemplate{
		Name: "envio_boletos",
		Body: "Olá {{1}}, total {{2}}.",
	})
	if err != nil {
		t.Fatalf("UpsertSessionTemplate() error = %v", err)
	}

	rendered, err := manager.renderTemplateMessage("session-a", &TemplateContent{
		Name: "envio_boletos",
		Components: []TemplateComponent{
			{
				Type: "body",
				Parameters: []TemplateParameter{
					{Type: "text", Text: "Maria"},
					{Type: "text", Text: "R$ 10,00"},
				},
			},
			{
				Type:    "button",
				SubType: "order_details",
				Index:   0,
				Parameters: []TemplateParameter{
					{
						Type: "action",
						Action: &TemplateAction{
							OrderDetails: map[string]any{
								"reference_id": "PED-1",
								"type":         "digital-goods",
								"payment_type": "br",
								"currency":     "BRL",
								"total_amount": map[string]any{"value": "1000", "offset": 100},
								"payment_settings": []any{
									map[string]any{
										"type": "pix_dynamic_code",
										"pix_dynamic_code": map[string]any{
											"code":          "000201",
											"merchant_name": "EnvolveNEXT",
											"key":           "57401868000110",
											"key_type":      "CNPJ",
										},
									},
									map[string]any{
										"type":   "boleto",
										"boleto": map[string]any{"digitable_line": "34191790010104351004791020150008291070026000"},
									},
								},
								"order": map[string]any{
									"status":   "pending",
									"subtotal": map[string]any{"value": "1000", "offset": 100},
									"items": []any{
										map[string]any{
											"retailer_id": "1234567",
											"name":        "Plataforma EnvolveNEXT",
											"amount":      map[string]any{"value": "1000", "offset": 100},
											"quantity":    "1",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("renderTemplateMessage() error = %v", err)
	}
	if got := rendered.Body; got != "Olá Maria, total R$ 10,00." {
		t.Fatalf("body = %q", got)
	}

	wire, err := manager.buildStructuredTemplateMessage(nil, nil, "session-a", rendered, nil)
	if err != nil {
		t.Fatalf("buildStructuredTemplateMessage() error = %v", err)
	}
	buttons := wire.GetInteractiveMessage().GetNativeFlowMessage().GetButtons()
	if len(buttons) != 1 {
		t.Fatalf("buttons len = %d", len(buttons))
	}
	if got := buttons[0].GetName(); got != "review_and_pay" {
		t.Fatalf("button name = %q", got)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(buttons[0].GetButtonParamsJSON()), &payload); err != nil {
		t.Fatalf("unmarshal button params: %v", err)
	}
	settings := payload["payment_settings"].([]any)
	if len(settings) != 2 {
		t.Fatalf("payment_settings len = %d", len(settings))
	}
	if got := payload["total_amount"].(map[string]any)["value"]; got != float64(1000) {
		t.Fatalf("total value = %#v", got)
	}
}

func TestSessionTemplatePersistsBySession(t *testing.T) {
	root := t.TempDir()
	templateFile := filepath.Join(root, "templates.json")
	manager := NewClientManager(ManagerOptions{
		SessionStore:        root,
		SessionTemplateFile: templateFile,
	})
	if _, err := manager.UpsertSessionTemplate("session-a", SessionTemplate{Name: "um", Body: "Olá"}); err != nil {
		t.Fatalf("UpsertSessionTemplate() error = %v", err)
	}

	reloaded := NewClientManager(ManagerOptions{
		SessionStore:        root,
		SessionTemplateFile: templateFile,
	})
	if got := len(reloaded.ListSessionTemplates("session-a")); got != 1 {
		t.Fatalf("templates len = %d", got)
	}
	if got := len(reloaded.ListSessionTemplates("session-b")); got != 0 {
		t.Fatalf("other session templates len = %d", got)
	}
}
