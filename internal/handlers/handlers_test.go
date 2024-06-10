package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, 200},
	{"about", "/about", "GET", []postData{}, 200},
	{"sb", "/double-bed", "GET", []postData{}, 200},
	{"db", "/single-bed", "GET", []postData{}, 200},
	{"sa", "/search-availability", "GET", []postData{}, 200},
	{"contact", "/contact", "GET", []postData{}, 200},
	{"mr", "/make-reservation", "GET", []postData{}, 200},

	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "01-01-2020"},
		{key: "end", value: "02-01-2020"},
	}, http.StatusOK},
	{"post-make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Parth"},
		{key: "last_name", value: "Vinchhi"},
		{key: "email", value: "Parth@gmail.com"},
		{key: "phone", value: "2323423453"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
