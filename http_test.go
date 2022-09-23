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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestHTTP(t *testing.T) {
	t.Parallel()

	req := &HTTPPayload{}
	field := HTTP(req)

	if diff := cmp.Diff(field, zap.Object("httpRequest", req)); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}

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

func TestHTTPPayload_MarshalLogObject(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader("12345"))
	req.Header.Set("User-Agent", "test-user-agent")
	req.Header.Set("Referer", "test-referer")
	resp := &http.Response{
		Body:       io.NopCloser(strings.NewReader("6789101112")),
		StatusCode: 200,
	}
	data := NewHTTPRequest(req, resp)

	enc := zapcore.NewMapObjectEncoder()
	if err := data.MarshalLogObject(enc); err != nil {
		t.Fatal(err)
	}

	if gotMethod, want := enc.Fields["requestMethod"], http.MethodPost; gotMethod != want {
		t.Fatalf("got %s but want %s", gotMethod, want)
	}

	if gotRequestUrl, want := enc.Fields["requestUrl"], "/"; gotRequestUrl != want {
		t.Fatalf("got %s but want %s", gotRequestUrl, want)
	}

	if gotRequestSize, want := enc.Fields["requestSize"], int64(5); gotRequestSize != want {
		t.Fatalf("got %d but want %d", gotRequestSize, want)
	}

	if gotResponseSize, want := enc.Fields["responseSize"], int64(10); gotResponseSize != want {
		t.Fatalf("got %d but want %d", gotResponseSize, want)
	}

	if gotUserAgent, want := enc.Fields["userAgent"], "test-user-agent"; gotUserAgent != want {
		t.Fatalf("got %s but want %s", gotUserAgent, want)
	}

	if gotRemoteIp, want := enc.Fields["remoteIp"], "192.0.2.1:1234"; gotRemoteIp != want {
		t.Fatalf("got %s but want %s", gotRemoteIp, want)
	}

	if gotServerIp, want := enc.Fields["serverIp"], ""; gotServerIp != want {
		t.Fatalf("got %s but want %s", gotServerIp, want)
	}

	if gotReferer, want := enc.Fields["referer"], "test-referer"; gotReferer != want {
		t.Fatalf("got %s but want %s", gotReferer, want)
	}

	if gotProtocol, want := enc.Fields["protocol"], "HTTP/1.1"; gotProtocol != want {
		t.Fatalf("got %s but want %s", gotProtocol, want)
	}

	if gotStatus, want := enc.Fields["status"], int32(http.StatusOK); gotStatus != want {
		t.Fatalf("got %d but want %d", gotStatus, want)
	}
}
