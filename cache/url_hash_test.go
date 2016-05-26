package cache

import (
	"net/url"
	"testing"
)

func TestURLHash(t *testing.T) {
	testCases := map[string]string{
		"http://example.com/foo":     "283f730bd92870df91cdded3819af2b1",
		"https://example.com/foo":    "283f730bd92870df91cdded3819af2b1",
		"http://example.com/FOO":     "283f730bd92870df91cdded3819af2b1",
		"http://example.com/foo?bar": "c0a036ac0556c603bf32298cb59993a6",
	}

	for rawURL, v := range testCases {
		u, _ := url.Parse(rawURL)
		h := URLHash(u)

		if u.String() != rawURL {
			t.Errorf("The URL got modified to %s", u.String())
		}

		if h != v {
			t.Errorf("Asking for %s, should have yielded %s, but returned %s", rawURL, v, h)
		}
	}
}
