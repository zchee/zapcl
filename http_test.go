// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestHTTPRequestField(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		r    *http.Request
		res  *http.Response
		want *HTTPPayload
	}{
		"Empty": {
			r:   nil,
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{},
			},
		},
		"RequestMethod": {
			r: &http.Request{
				Method: "GET",
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RequestMethod: "GET",
				},
			},
		},
		"Status": {
			r:   nil,
			res: &http.Response{StatusCode: 404},
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					Status: 404,
				},
			},
		},
		"UserAgent": {
			r: &http.Request{
				Header: http.Header{
					"User-Agent": []string{"hello world"},
				},
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					UserAgent: "hello world",
				},
			},
		},
		"RemoteIP": {
			r: &http.Request{
				RemoteAddr: "127.0.0.1",
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RemoteIp: "127.0.0.1",
				},
			},
		},
		"Referrer": {
			r: &http.Request{
				Header: http.Header{
					"Referer": []string{"hello universe"},
				},
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					Referer: "hello universe",
				},
			},
		},
		"Protocol": {
			r: &http.Request{
				Proto: "HTTP/1.1",
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					Protocol: "HTTP/1.1",
				},
			},
		},
		"RequestURL": {
			r: &http.Request{
				URL: &url.URL{
					Host:   "example.com",
					Scheme: "https",
				},
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RequestUrl: "https://example.com",
				},
			},
		},
		"RequestSize": {
			r: &http.Request{
				Body: io.NopCloser(strings.NewReader("12345")),
			},
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RequestSize: 5,
				},
			},
		},
		"ResponseSize": {
			r: nil,
			res: &http.Response{
				Body: io.NopCloser(strings.NewReader("12345")),
			},
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					ResponseSize: 5,
				},
			},
		},
		"SimpleRequest": {
			r:   httptest.NewRequest("POST", "/", strings.NewReader("12345")),
			res: nil,
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RequestSize:   5,
					RequestMethod: "POST",
					RemoteIp:      "192.0.2.1:1234",
					Protocol:      "HTTP/1.1",
					RequestUrl:    "/",
				},
			},
		},
		"SimpleResponse": {
			r: nil,
			res: &http.Response{
				Body:       io.NopCloser(strings.NewReader("12345")),
				StatusCode: 404,
			},
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					ResponseSize: 5,
					Status:       404,
				},
			},
		},
		"RequestAndResponse": {
			r: &http.Request{
				Method: "POST",
				Proto:  "HTTP/1.1",
			},
			res: &http.Response{StatusCode: 200},
			want: &HTTPPayload{
				HttpRequest: &logtypepb.HttpRequest{
					RequestMethod: "POST",
					Protocol:      "HTTP/1.1",
					Status:        200,
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tt.want, NewHTTPRequest(tt.r, tt.res),
				protocmp.Transform(),
			); diff != "" {
				t.Fatalf("(-want, +got)\n%s\n", diff)
			}
		})
	}
}
