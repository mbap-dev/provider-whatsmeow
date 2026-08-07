package provider

import (
	"testing"

	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
)

func TestApplyLegacyListBizType(t *testing.T) {
	names := map[int32]string{
		int32(waE2E.ListMessage_SINGLE_SELECT): "SINGLE_SELECT",
	}

	if err := applyLegacyListBizType(names); err != nil {
		t.Fatalf("applyLegacyListBizType() error = %v", err)
	}
	if got := names[int32(waE2E.ListMessage_SINGLE_SELECT)]; got != "PRODUCT_LIST" {
		t.Fatalf("single-select enum name = %q, want PRODUCT_LIST", got)
	}
	if err := applyLegacyListBizType(names); err != nil {
		t.Fatalf("second applyLegacyListBizType() error = %v", err)
	}
}

func TestApplyLegacyListBizTypeRejectsUnexpectedUpstreamShape(t *testing.T) {
	names := map[int32]string{
		int32(waE2E.ListMessage_SINGLE_SELECT): "CHANGED_UPSTREAM",
	}

	if err := applyLegacyListBizType(names); err == nil {
		t.Fatal("applyLegacyListBizType() error = nil, want unsupported-shape error")
	}
}

func TestOfficialWhatsmeowUsesLegacyListBizType(t *testing.T) {
	if got := waE2E.ListMessage_ListType_name[int32(waE2E.ListMessage_SINGLE_SELECT)]; got != "PRODUCT_LIST" {
		t.Fatalf("runtime single-select enum name = %q, want PRODUCT_LIST", got)
	}
}
