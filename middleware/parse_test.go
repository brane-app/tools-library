package middleware

import (
	"net/http"
	"net/url"
	"testing"
)

type querySet struct {
	URL   string
	Size  int
	After string
	Code  int
	OK    bool
}

func Test_RangeQueryParams(test *testing.T) {
	var request *http.Request = new(http.Request)
	var q_defaults map[string]interface{} = defaults()

	var set querySet
	var sets []querySet = []querySet{
		querySet{
			URL:   "http://imonke.io/?&size=10",
			Size:  10,
			After: "",
			Code:  0,
			OK:    true,
		},
		querySet{
			URL:   "http://imonke.io/?after=foobar_baz&size=50",
			Size:  q_defaults["size"].(int),
			After: "foobar_baz",
			Code:  0,
			OK:    true,
		},
		querySet{
			URL:   "http://imonke.io/?after=jol",
			Size:  q_defaults["size"].(int),
			After: "jol",
			Code:  0,
			OK:    true,
		},
		querySet{
			URL:   "http://imonke.io/?&size=300",
			Size:  RANGE_SIZE_LIMIT,
			After: q_defaults["after"].(string),
			Code:  0,
			OK:    true,
		},
		querySet{
			URL:   "http://imonke.io/?&size=lol",
			Size:  q_defaults["size"].(int),
			After: q_defaults["after"].(string),
			Code:  400,
			OK:    false,
		},
		querySet{
			URL:   "http://imonke.io/?&size=-2",
			Size:  q_defaults["size"].(int),
			After: q_defaults["after"].(string),
			Code:  400,
			OK:    false,
		},
	}

	var parsed map[string]interface{}
	var modified *http.Request
	var ok bool
	var code int
	var err error
	for _, set = range sets {
		if request.URL, err = url.Parse(set.URL); err != nil {
			test.Fatal(err)
		}

		if modified, ok, code, _, err = RangeQueryParams(request); err != nil {
			test.Fatal(err)
		}

		if code != set.Code {
			test.Errorf("got code %d", code)
		}

		if ok != set.OK {
			test.Errorf("got ok %t", ok)
		}

		if !ok {
			continue
		}

		parsed = modified.Context().Value("query").(map[string]interface{})

		if parsed["size"].(int) != set.Size {
			test.Errorf("size mismatch! have: %d, want: %d", parsed["size"], set.Size)
		}

		if parsed["after"].(string) != set.After {
			test.Errorf("after mismatch! have: %s, want: %s", parsed["after"], set.After)
		}
	}
}
