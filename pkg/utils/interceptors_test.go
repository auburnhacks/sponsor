package utils

import "testing"

func TestTrueUnauthenticatedRPCCall(t *testing.T) {
	fullMethod := "/proto.SponsorService/LoginAdmin"
	if err := isUnauthenticatedRPC(fullMethod); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestInvalidRPCCall(t *testing.T) {
	fullMethod := "/proto.SponsorService/FooBar"
	err := isUnauthenticatedRPC(fullMethod)
	if err == nil {
		t.Fail()
	}
}
