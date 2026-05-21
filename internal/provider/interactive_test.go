package provider

import (
	"encoding/json"
	"testing"

	waBinary "go.mau.fi/whatsmeow/binary"
)

func TestBuildInteractiveListMessageFromCloudPayload(t *testing.T) {
	wire, err := buildInteractiveListMessage(&InteractiveContent{
		Type: "list",
		Body: &InteractiveText{Text: "Opção invalida, tente novamente."},
		Action: &InteractiveAction{
			Button: "Ver opções",
			Sections: []InteractiveSection{
				{
					Title: "Opções disponíveis",
					Rows: []InteractiveRow{
						{ID: "option_1", Title: "Aro 13"},
						{ID: "option_2", Title: "Aro 14"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("buildInteractiveListMessage() error = %v", err)
	}
	list := wire.GetListMessage()
	if list == nil {
		t.Fatal("expected list message")
	}
	if got := list.GetTitle(); got != "Opção invalida, tente novamente." {
		t.Fatalf("title = %q", got)
	}
	if got := list.GetButtonText(); got != "Ver opções" {
		t.Fatalf("button text = %q", got)
	}
	if got := list.GetFooterText(); got != "Ver opções" {
		t.Fatalf("footer = %q", got)
	}
	if got := list.GetSections()[0].GetRows()[1].GetRowID(); got != "option_2" {
		t.Fatalf("row id = %q", got)
	}
}

func TestBuildInteractiveButtonMessageFromCloudPayload(t *testing.T) {
	wire, err := buildInteractiveButtonMessage(&InteractiveContent{
		Type: "button",
		Body: &InteractiveText{Text: "Opção invalida, tente novamente."},
		Action: &InteractiveAction{
			Buttons: []InteractiveButton{
				{
					Type: "reply",
					Reply: InteractiveReply{
						ID:    "option_1",
						Title: "Sim",
					},
				},
				{
					Type: "reply",
					Reply: InteractiveReply{
						ID:    "option_2",
						Title: "Não, pode encerrar.",
					},
				},
			},
		},
	}, nil)
	if err != nil {
		t.Fatalf("buildInteractiveButtonMessage() error = %v", err)
	}

	interactive := wire.GetInteractiveMessage()
	if interactive == nil {
		t.Fatal("expected interactive message")
	}
	if got := interactive.GetBody().GetText(); got != "Opção invalida, tente novamente." {
		t.Fatalf("body = %q", got)
	}

	native := interactive.GetNativeFlowMessage()
	if native == nil {
		t.Fatal("expected native flow message")
	}
	if got := len(native.GetButtons()); got != 2 {
		t.Fatalf("buttons len = %d", got)
	}
	if got := native.GetButtons()[0].GetName(); got != "quick_reply" {
		t.Fatalf("button type = %q", got)
	}

	var payload map[string]string
	if err := json.Unmarshal([]byte(native.GetButtons()[1].GetButtonParamsJSON()), &payload); err != nil {
		t.Fatalf("unmarshal button params: %v", err)
	}
	if got := payload["id"]; got != "option_2" {
		t.Fatalf("button id = %q", got)
	}
	if got := payload["display_text"]; got != "Não, pode encerrar." {
		t.Fatalf("button title = %q", got)
	}
}

func TestInteractiveButtonSendExtraUsesMixedNativeFlowBizNode(t *testing.T) {
	extra := interactiveButtonSendExtra()
	if extra == nil || extra.AdditionalNodes == nil || len(*extra.AdditionalNodes) != 1 {
		t.Fatalf("unexpected send extra: %#v", extra)
	}

	biz := (*extra.AdditionalNodes)[0]
	if biz.Tag != "biz" {
		t.Fatalf("biz tag = %q", biz.Tag)
	}
	children, ok := biz.Content.([]waBinary.Node)
	if !ok || len(children) != 1 {
		t.Fatalf("unexpected biz content: %#v", biz.Content)
	}
	if got := children[0].Tag; got != "interactive" {
		t.Fatalf("interactive tag = %q", got)
	}
	if got := children[0].Attrs["type"]; got != "native_flow" {
		t.Fatalf("interactive type = %#v", got)
	}
}
