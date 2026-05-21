package provider

import (
	"encoding/json"
	"testing"

	goE2E "go.mau.fi/whatsmeow/proto/waE2E"
)

func TestBuildOrderDetailsMessageWithPIXAndBoleto(t *testing.T) {
	wire, extra, err := buildOrderDetailsMessage(&PaymentContent{
		ReferenceID: "PED-123",
		Body:        "Segue a cobrança do pedido.",
		TotalAmount: PaymentAmount{Value: 455400, Offset: 100},
		ItemName:    "Pedido 123",
		PIXStaticCode: &PIXStaticCodeInfo{
			MerchantName: "Loja Exemplo",
			Key:          "pix@example.com",
			KeyType:      "email",
		},
		Boleto: &BoletoInfo{
			DigitableLine: "34191790010104351004791020150008291070026000",
		},
	}, &goE2E.DocumentMessage{FileName: strptr("boleto.pdf")}, nil)
	if err != nil {
		t.Fatalf("buildOrderDetailsMessage() error = %v", err)
	}
	if wire.GetInteractiveMessage() == nil {
		t.Fatal("expected interactive message")
	}
	if len(wire.GetMessageContextInfo().GetMessageSecret()) != 32 {
		t.Fatalf("message secret len = %d", len(wire.GetMessageContextInfo().GetMessageSecret()))
	}
	if wire.GetInteractiveMessage().GetHeader() == nil || !wire.GetInteractiveMessage().GetHeader().GetHasMediaAttachment() {
		t.Fatal("expected PDF header attachment")
	}

	native := wire.GetInteractiveMessage().GetNativeFlowMessage()
	if native == nil || len(native.GetButtons()) != 1 {
		t.Fatalf("buttons = %d", len(native.GetButtons()))
	}
	if got := native.GetButtons()[0].GetName(); got != "review_and_pay" {
		t.Fatalf("button name = %q", got)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(native.GetButtons()[0].GetButtonParamsJSON()), &payload); err != nil {
		t.Fatalf("unmarshal button params: %v", err)
	}
	if got := payload["reference_id"]; got != "PED-123" {
		t.Fatalf("reference_id = %#v", got)
	}
	if got := payload["currency"]; got != "BRL" {
		t.Fatalf("currency = %#v", got)
	}

	settings, ok := payload["payment_settings"].([]any)
	if !ok {
		t.Fatalf("payment_settings type = %T", payload["payment_settings"])
	}
	if len(settings) != 2 {
		t.Fatalf("payment_settings len = %d", len(settings))
	}

	pix, ok := settings[0].(map[string]any)
	if !ok {
		t.Fatalf("pix setting type = %T", settings[0])
	}
	pixPayload, ok := pix["pix_static_code"].(map[string]any)
	if !ok {
		t.Fatalf("pix_static_code type = %T", pix["pix_static_code"])
	}
	if got := pixPayload["key_type"]; got != "EMAIL" {
		t.Fatalf("pix key_type = %#v", got)
	}

	boleto, ok := settings[1].(map[string]any)
	if !ok {
		t.Fatalf("boleto setting type = %T", settings[1])
	}
	boletoPayload, ok := boleto["boleto"].(map[string]any)
	if !ok {
		t.Fatalf("boleto payload type = %T", boleto["boleto"])
	}
	if got := boletoPayload["digitable_line"]; got != "34191790010104351004791020150008291070026000" {
		t.Fatalf("digitable_line = %#v", got)
	}

	if extra == nil || extra.AdditionalNodes == nil || len(*extra.AdditionalNodes) != 1 {
		t.Fatal("expected biz additional node")
	}
	if got := (*extra.AdditionalNodes)[0].Attrs["native_flow_name"]; got != "order_details" {
		t.Fatalf("native_flow_name = %#v", got)
	}
}

func TestBuildOrderDetailsMessageDefaults(t *testing.T) {
	wire, _, err := buildOrderDetailsMessage(&PaymentContent{
		TotalAmount: PaymentAmount{Value: 10990},
		PIXStaticCode: &PIXStaticCodeInfo{
			Key: "8e6a4a11-1234-5678-90ab-abcdef123456",
		},
	}, nil, nil)
	if err != nil {
		t.Fatalf("buildOrderDetailsMessage() error = %v", err)
	}

	if got := wire.GetInteractiveMessage().GetBody().GetText(); got != "Segue a cobrança. Você pode pagar pelo botão abaixo." {
		t.Fatalf("body = %q", got)
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(wire.GetInteractiveMessage().GetNativeFlowMessage().GetButtons()[0].GetButtonParamsJSON()), &payload); err != nil {
		t.Fatalf("unmarshal button params: %v", err)
	}

	order, ok := payload["order"].(map[string]any)
	if !ok {
		t.Fatalf("order type = %T", payload["order"])
	}
	if got := order["status"]; got != "pending" {
		t.Fatalf("order status = %#v", got)
	}
	items, ok := order["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("order items = %#v", order["items"])
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item type = %T", items[0])
	}
	if got := item["name"]; got != "Cobrança" {
		t.Fatalf("item name = %#v", got)
	}

	settings := payload["payment_settings"].([]any)
	pix := settings[0].(map[string]any)["pix_static_code"].(map[string]any)
	if got := pix["merchant_name"]; got != "Cobrança" {
		t.Fatalf("merchant_name = %#v", got)
	}
	if got := pix["key_type"]; got != "EVP" {
		t.Fatalf("key_type = %#v", got)
	}
}

func strptr(v string) *string {
	return &v
}
