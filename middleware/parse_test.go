package middleware

import (
	"net/http"
	"net/url"
	"testing"
)

type querySet struct {
	URL    string
	Size   int
	Offset int
	Code   int
	OK     bool
}

func Test_RangeQueryParams(test *testing.T) {
	var request *http.Request = new(http.Request)
	var q_defaults map[string]int = defaults()

	var set querySet
	var sets []querySet = []querySet{
		querySet{
			URL:    "http://imonke.io/?offset=10&size=10",
			Size:   10,
			Offset: 10,
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?offset=100000&size=-3",
			Size:   q_defaults["size"],
			Offset: 100000,
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?offset=-3&size=-3",
			Size:   q_defaults["size"],
			Offset: q_defaults["offset"],
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?&size=300",
			Size:   limits["size"],
			Offset: q_defaults["offset"],
			Code:   0,
			OK:     true,
		},
		querySet{
			URL:    "http://imonke.io/?&size=lol",
			Size:   q_defaults["size"],
			Offset: q_defaults["offset"],
			Code:   400,
			OK:     false,
		},
		querySet{
			URL:    "http://imonke.io/?&offset=42069",
			Size:   q_defaults["size"],
			Offset: 42069,
			Code:   0,
			OK:     true,
		},
	}

	var parsed map[string]int

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

		parsed = modified.Context().Value("parsed_query").(map[string]int)

		if parsed["size"] != set.Size {
			test.Errorf("size mismatch! have: %d, want: %d", parsed["size"], set.Size)
		}

		if parsed["offset"] != set.Offset {
			test.Errorf("offset mismatch! have: %d, want: %d", parsed["offset"], set.Offset)
		}
	}
}
