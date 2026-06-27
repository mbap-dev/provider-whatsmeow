package provider

import (
	"reflect"
	"testing"
)

func TestCandidatesBRAddsNinthDigitForEightDigitMobile(t *testing.T) {
	got := candidatesBR("+556296303132")
	want := []string{"+556296303132", "+5562996303132"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("candidatesBR() = %#v, want %#v", got, want)
	}
}

func TestCandidatesBRKeepsNineDigitMobileAndAddsLegacyCandidate(t *testing.T) {
	got := candidatesBR("+5562996303132")
	want := []string{"+5562996303132", "+556296303132"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("candidatesBR() = %#v, want %#v", got, want)
	}
}
