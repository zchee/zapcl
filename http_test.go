// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHTTPRequestField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req  *http.Request
		res  *http.Response
		want *HTTPPayload
	}{
		"Empty": {
			nil,
			nil,
			&HTTPPayload{},
		},

		"RequestMethod": {
			&http.Request{
				Method: "GET",
			},
			nil,
			&HTTPPayload{
				RequestMethod: "GET",
			},
		},

		"Status": {
			nil,
			&http.Response{StatusCode: 404},
			&HTTPPayload{Status: 404},
		},

		"UserAgent": {
			&http.Request{Header: http.Header{"User-Agent": []string{"hello world"}}},
			nil,
			&HTTPPayload{UserAgent: "hello world"},
		},

		"RemoteIP": {
			&http.Request{RemoteAddr: "127.0.0.1"},
			nil,
			&HTTPPayload{RemoteIP: "127.0.0.1"},
		},

		"Referrer": {
			&http.Request{Header: http.Header{"Referer": []string{"hello universe"}}},
			nil,
			&HTTPPayload{Referer: "hello universe"},
		},

		"Protocol": {
			&http.Request{Proto: "HTTP/1.1"},
			nil,
			&HTTPPayload{Protocol: "HTTP/1.1"},
		},

		"RequestURL": {
			&http.Request{URL: &url.URL{Host: "example.com", Scheme: "https"}},
			nil,
			&HTTPPayload{RequestURL: "https://example.com"},
		},

		"RequestSize": {
			&http.Request{Body: io.NopCloser(strings.NewReader("12345"))},
			nil,
			&HTTPPayload{RequestSize: "5"},
		},

		"ResponseSize": {
			nil,
			&http.Response{Body: io.NopCloser(strings.NewReader("12345"))},
			&HTTPPayload{ResponseSize: "5"},
		},

		"simple request": {
			httptest.NewRequest("POST", "/", strings.NewReader("12345")),
			nil,
			&HTTPPayload{
				RequestSize:   "5",
				RequestMethod: "POST",
				RemoteIP:      "192.0.2.1:1234",
				Protocol:      "HTTP/1.1",
				RequestURL:    "/",
			},
		},

		"simple response": {
			nil,
			&http.Response{Body: io.NopCloser(strings.NewReader("12345")), StatusCode: 404},
			&HTTPPayload{ResponseSize: "5", Status: 404},
		},

		"request & response": {
			&http.Request{Method: "POST", Proto: "HTTP/1.1"},
			&http.Response{StatusCode: 200},
			&HTTPPayload{RequestMethod: "POST", Protocol: "HTTP/1.1", Status: 200},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tt.want, NewHTTPRequest(tt.req, tt.res)); diff != "" {
				t.Fatalf("(-want, +got)\n%s\n", diff)
			}
		})
	}
}
