package googlemaps_share_link_decoder_test

import (
	decoder "github.com/minchao/googlemaps-share-link-decoder"
	"testing"
)

func TestDecode(t *testing.T) {
	svc := decoder.ShareLinkService{}
	if _, e := svc.Decode(&decoder.Request{"https://goo.gl/maps/yvhHAqiQKfs"}); e != nil {
		t.Errorf("decode failed: %s", e)
	} else {
		t.Log("decode succeeded")
	}
}
