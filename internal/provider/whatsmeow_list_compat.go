package provider

import (
	"fmt"

	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
)

const (
	officialSingleSelectName = "SINGLE_SELECT"
	legacyListBizTypeName    = "PRODUCT_LIST"
)

func init() {
	if err := applyLegacyListBizType(waE2E.ListMessage_ListType_name); err != nil {
		panic(err)
	}
}

// Whatsmeow derives the biz/list type from this generated enum-name map. Only
// change that lookup: the protobuf and reverse enum map must stay SINGLE_SELECT.
func applyLegacyListBizType(names map[int32]string) error {
	key := int32(waE2E.ListMessage_SINGLE_SELECT)
	switch current := names[key]; current {
	case officialSingleSelectName:
		names[key] = legacyListBizTypeName
		return nil
	case legacyListBizTypeName:
		return nil
	default:
		return fmt.Errorf("unsupported whatsmeow ListMessage_SINGLE_SELECT enum name %q", current)
	}
}
