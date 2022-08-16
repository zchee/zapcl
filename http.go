// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	httpRequestKey = "httpRequest"
)

// HTTPPayload represents a Cloud Logging httpRequest fields.
type HTTPPayload struct {
	// RequestMethod is the request method. Examples: "GET", "HEAD", "PUT", "POST".
	RequestMethod string `json:"requestMethod"`

	// RequestUrl is the scheme (http, https), the host name, the path and the query portion of the URL that was requested. Example: "http://example.com/some/info?color=red".
	RequestUrl string `json:"requestUrl"`

	// RequestSize is THE size of the HTTP request message in bytes, including the request headers and the request body.
	RequestSize string `json:"requestSize"` // int64 format

	// Status is the response code indicating the Status of response. Examples: 200, 404.
	Status int `json:"status"`

	// ResponseSize is the size of the HTTP response message sent back to the client, in bytes, including the response headers and the response body.
	ResponseSize string `json:"responseSize"` // int64 format

	// UserAgent is the user agent sent by the client. Example: "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)".
	UserAgent string `json:"userAgent"`

	// RemoteIP is the IP address (IPv4 or IPv6) of the client that issued the HTTP request. This field can include port information. Examples: "192.168.1.1", "10.0.0.1:80", "FE80::0202:B3FF:FE1E:8329".
	RemoteIP string `json:"remoteIp"`

	// ServerIP is the IP address (IPv4 or IPv6) of the origin server that the request was sent to. This field can include port information. Examples: "192.168.1.1", "10.0.0.1:80", "FE80::0202:B3FF:FE1E:8329".
	ServerIP string `json:"serverIp"`

	// Referer is the Referer URL of the request, as defined in HTTP/1.1 Header Field Definitions.
	Referer string `json:"referer"`

	// Latency is the request processing Latency on the server, from the time the request was received until the response was sent.
	// A duration in seconds with up to nine fractional digits, ending with 's'. Example: "3.5s".
	Latency string `json:"latency"` // Duration format

	// CacheLookup whether or not a cache lookup was attempted.
	CacheLookup bool `json:"cacheLookup"`

	// CacheHit whether or not an entity was served from cache (with or without validation).
	CacheHit bool `json:"cacheHit"`

	// CacheValidatedWithOriginServer whether or not the response was validated with the origin server before being served from cache. This field is only meaningful if cacheHit is True.
	CacheValidatedWithOriginServer bool `json:"cacheValidatedWithOriginServer"`

	// CacheFillBytes is the number of HTTP response bytes inserted into cache. Set only when a cache fill was attempted.
	CacheFillBytes string `json:"cacheFillBytes"` // int64 format

	// Protocol used for the request. Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `json:"protocol"`
}

// MarshalLogObject implements zapcore.ObjectMarshaler.
func (p *HTTPPayload) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("requestMethod", p.RequestMethod)
	enc.AddString("requestUrl", p.RequestUrl)
	enc.AddString("requestSize", p.RequestSize)
	enc.AddInt("status", p.Status)
	enc.AddString("responseSize", p.ResponseSize)
	enc.AddString("userAgent", p.UserAgent)
	enc.AddString("remoteIp", p.RemoteIP)
	enc.AddString("serverIp", p.ServerIP)
	enc.AddString("referer", p.Referer)
	enc.AddString("latency", p.Latency)
	enc.AddBool("cacheLookup", p.CacheLookup)
	enc.AddBool("cacheHit", p.CacheHit)
	enc.AddBool("cacheValidatedWithOriginServer", p.CacheValidatedWithOriginServer)
	enc.AddString("cacheFillBytes", p.CacheFillBytes)
	enc.AddString("protocol", p.Protocol)

	return nil
}

// HTTPRequestField adds the Cloud Logging "httpRequest" fields from req and resp.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func HTTPRequestField(r *http.Request, res *http.Response) zapcore.Field {
	req := &HTTPPayload{
		RequestMethod: r.Method,
		Status:        res.StatusCode,
		UserAgent:     r.UserAgent(),
		RemoteIP:      r.RemoteAddr,
		Referer:       r.Referer(),
		Protocol:      r.Proto,
	}

	if url := r.URL; url != nil {
		req.RequestUrl = url.String()
	}

	buf := new(bytes.Buffer)
	if body := r.Body; body != nil {
		n, _ := io.Copy(buf, body)
		req.RequestSize = strconv.FormatInt(n, 10)
	}

	if body := res.Body; body != nil {
		buf.Reset()
		n, _ := io.Copy(buf, body)
		req.ResponseSize = strconv.FormatInt(n, 10)
	}

	return zap.Object(httpRequestKey, req)
}
