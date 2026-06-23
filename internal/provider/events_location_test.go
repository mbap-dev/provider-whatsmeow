package provider

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestBuildInboundLocationPayload(t *testing.T) {
	payload := buildInboundLocationPayload(&waE2E.LocationMessage{
		DegreesLatitude:  proto.Float64(-23.55052),
		DegreesLongitude: proto.Float64(-46.633308),
	})

	if payload["latitude"] != -23.55052 {
		t.Fatalf("latitude = %v", payload["latitude"])
	}
	if payload["longitude"] != -46.633308 {
		t.Fatalf("longitude = %v", payload["longitude"])
	}
}

func TestBuildInboundLiveLocationPayloadMatchesCommonLocationShape(t *testing.T) {
	payload := buildInboundLiveLocationPayload(&waE2E.LiveLocationMessage{
		DegreesLatitude:  proto.Float64(-22.906847),
		DegreesLongitude: proto.Float64(-43.172897),
		AccuracyInMeters: proto.Uint32(12),
		SpeedInMps:       proto.Float32(1.5),
		Caption:          proto.String("indo"),
	})

	if len(payload) != 2 {
		t.Fatalf("payload keys = %v, want only latitude and longitude", payload)
	}
	if payload["latitude"] != -22.906847 {
		t.Fatalf("latitude = %v", payload["latitude"])
	}
	if payload["longitude"] != -43.172897 {
		t.Fatalf("longitude = %v", payload["longitude"])
	}
}
