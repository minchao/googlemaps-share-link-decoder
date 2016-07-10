package googlemaps_share_link_decoder_test

import (
	"testing"

	decoder "github.com/minchao/googlemaps-share-link-decoder"
)

func TestDecode(t *testing.T) {
	svc := decoder.ShareLinkService{}
	if _, e := svc.Decode(&decoder.Request{"https://goo.gl/maps/yvhHAqiQKfs"}); e != nil {
		t.Errorf("decode failed: %s", e)
	}
}

func TestDataNotFound(t *testing.T) {
	svc := decoder.ShareLinkService{}
	// Route planner
	if _, e := svc.Decode(&decoder.Request{"https://goo.gl/maps/4MwturqawME2"}); e != nil {
		if e.Error() != "place data not found" {
			t.Errorf("decode failed: %s", e)
		}
	} else {
		t.Error("error is nil !")
	}
	if _, e := svc.Decode(&decoder.Request{"https://goo.gl/maps/CWjrVFgAomM2"}); e != nil {
		if e.Error() != "location data not found" {
			t.Errorf("decode failed: %s", e)
		}
	} else {
		t.Error("error is nil !")
	}
}
