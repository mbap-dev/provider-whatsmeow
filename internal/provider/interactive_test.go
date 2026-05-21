package provider

import "testing"

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
