package utils

import "testing"

func TestIsUthenticatedRPC(t *testing.T) {
	fullMethod := "/proto.SponsorService/LoginAdmin"
	if err := isUnauthenticatedRPC(fullMethod); err != nil {
		t.Fatalf("%v", err)
	}
}
