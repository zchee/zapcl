// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"bytes"
	"io"
	"net/http"
	"unicode/utf8"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"
)

const (
	httpRequestKey = "httpRequest"
)

// HTTPPayload represents a Cloud Logging httpRequest fields.
type HTTPPayload struct {
	*logtypepb.HttpRequest
}

var _ zapcore.ObjectMarshaler = (*HTTPPayload)(nil)

// MarshalLogObject implements zapcore.ObjectMarshaler.
func (p *HTTPPayload) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("requestMethod", p.GetRequestMethod())
	enc.AddString("requestUrl", p.GetRequestUrl())
	enc.AddInt64("requestSize", p.GetRequestSize())
	enc.AddInt64("responseSize", p.GetResponseSize())
	enc.AddString("userAgent", p.GetUserAgent())
	enc.AddString("remoteIp", p.GetRemoteIp())
	enc.AddString("serverIp", p.GetServerIp())
	enc.AddString("referer", p.GetReferer())
	enc.AddDuration("latency", p.GetLatency().AsDuration())
	enc.AddInt64("cacheFillBytes", p.GetCacheFillBytes())
	enc.AddString("protocol", p.GetProtocol())
	enc.AddInt32("status", p.GetStatus())
	enc.AddBool("cacheLookup", p.GetCacheLookup())
	enc.AddBool("cacheHit", p.GetCacheHit())
	enc.AddBool("cacheValidatedWithOriginServer", p.GetCacheValidatedWithOriginServer())

	return nil
}

// NewHTTPRequest returns a new HTTPPayload struct, based on the passed
// in http.Request and http.Response objects.
func NewHTTPRequest(r *http.Request, res *http.Response) *HTTPPayload {
	if r == nil {
		r = &http.Request{}
	}
	if res == nil {
		res = &http.Response{}
	}

	req := &HTTPPayload{
		HttpRequest: &logtypepb.HttpRequest{
			RequestMethod: r.Method,
			Status:        int32(res.StatusCode),
			UserAgent:     r.UserAgent(),
			RemoteIp:      r.RemoteAddr,
			Referer:       r.Referer(),
			Protocol:      r.Proto,
		},
	}

	if url := r.URL; url != nil {
		u := *r.URL
		u.Fragment = ""
		req.RequestUrl = fixUTF8(u.String())
	}

	buf := new(bytes.Buffer)
	if body := r.Body; body != nil {
		n, _ := io.Copy(buf, body)
		req.RequestSize = n
	}

	if body := res.Body; body != nil {
		buf.Reset()
		n, _ := io.Copy(buf, body)
		req.ResponseSize = n
	}

	return req
}

// fixUTF8 is a helper that fixes an invalid UTF-8 string by replacing
// invalid UTF-8 runes with the Unicode replacement character (U+FFFD).
// See Issue https://github.com/googleapis/google-cloud-go/issues/1383.
func fixUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Otherwise time to build the sequence.
	buf := new(bytes.Buffer)
	buf.Grow(len(s))
	for _, r := range s {
		if utf8.ValidRune(r) {
			buf.WriteRune(r)
		} else {
			buf.WriteRune('\uFFFD')
		}
	}
	return buf.String()
}

// HTTP adds the correct Stackdriver "httpRequest" field.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func HTTP(req *HTTPPayload) zap.Field {
	return zap.Object(httpRequestKey, req)
}
