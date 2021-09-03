package zapdriver

// "Broker: Request timed out"
// https://console.cloud.google.com/logs/viewer?project=bnl-blendle&minLogLevel=
// 0&expandAll=false&timestamp=2018-05-23T22:21:56.142000000Z&customFacets=&limi
// tCustomFacetWidth=true&dateRangeEnd=2018-05-23T22:21:52.545Z&interval=PT1H&re
// source=container%2Fcluster_name%2Fblendle-2%2Fnamespace_id%2Fstream-
// composition-analytic-events-
// backfill&scrollTimestamp=2018-05-23T05:29:33.000000000Z&logName=projects
// %2Fbnl-blendle%2Flogs%2Fstream-composition-analytic-events-
// pipe-1&dateRangeUnbound=backwardInTime

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTP adds the correct Stackdriver "HTTP" field.
//
// see: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func HTTP(req *HTTPPayload) zap.Field {
	return zap.Object("httpRequest", req)
}

// HTTPPayload is the complete payload that can be interpreted by
// Stackdriver as a HTTP request.
type HTTPPayload struct {
	// The request method. Examples: "GET", "HEAD", "PUT", "POST".
	RequestMethod string `json:"requestMethod,omitempty"`

	// The scheme (http, https), the host name, the path and the query portion of
	// the URL that was requested.
	//
	// Example: "http://example.com/some/info?color=red".
	RequestURL string `json:"requestUrl,omitempty"`

	// The size of the HTTP request message in bytes, including the request
	// headers and the request body.
	RequestSize string `json:"requestSize,omitempty"`

	// The response code indicating the status of response.
	//
	// Examples: 200, 404.
	Status int `json:"status,omitempty"`

	// The size of the HTTP response message sent back to the client, in bytes,
	// including the response headers and the response body.
	ResponseSize string `json:"responseSize,omitempty"`

	// The user agent sent by the client.
	//
	// Example: "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98; Q312461; .NET CLR 1.0.3705)".
	UserAgent string `json:"userAgent,omitempty"`

	// The IP address (IPv4 or IPv6) of the client that issued the HTTP request.
	//
	// Examples: "192.168.1.1", "FE80::0202:B3FF:FE1E:8329".
	RemoteIP string `json:"remoteIp,omitempty"`

	// The IP address (IPv4 or IPv6) of the origin server that the request was
	// sent to.
	ServerIP string `json:"serverIp,omitempty"`

	// The referrer URL of the request, as defined in HTTP/1.1 Header Field
	// Definitions.
	Referer string `json:"referer,omitempty"`

	// The request processing latency on the server, from the time the request was
	// received until the response was sent.
	//
	// A duration in seconds with up to nine fractional digits, terminated by 's'.
	//
	// Example: "3.5s".
	Latency string `json:"latency,omitempty"`

	// Whether or not a cache lookup was attempted.
	CacheLookup bool `json:"cacheLookup,omitempty"`

	// Whether or not an entity was served from cache (with or without
	// validation).
	CacheHit bool `json:"cacheHit,omitempty"`

	// Whether or not the response was validated with the origin server before
	// being served from cache. This field is only meaningful if cacheHit is True.
	CacheValidatedWithOriginServer bool `json:"cacheValidatedWithOriginServer,omitempty"`

	// The number of HTTP response bytes inserted into cache. Set only when a
	// cache fill was attempted.
	CacheFillBytes string `json:"cacheFillBytes,omitempty"`

	// Protocol used for the request.
	//
	// Examples: "HTTP/1.1", "HTTP/2", "websocket"
	Protocol string `json:"protocol,omitempty"`
}

// NewHTTP returns a new HTTPPayload struct, based on the passed
// in http.Request and http.Response objects.
func NewHTTP(req *http.Request, res *http.Response) *HTTPPayload {
	if req == nil {
		req = &http.Request{}
	}

	if res == nil {
		res = &http.Response{}
	}

	sdreq := &HTTPPayload{
		RequestMethod: req.Method,
		Status:        res.StatusCode,
		UserAgent:     req.UserAgent(),
		RemoteIP:      req.RemoteAddr,
		Referer:       req.Referer(),
		Protocol:      req.Proto,
	}

	if req.URL != nil {
		sdreq.RequestURL = req.URL.String()
	}

	buf := &bytes.Buffer{}
	if req.Body != nil {
		n, _ := io.Copy(buf, req.Body) // nolint: gas
		sdreq.RequestSize = strconv.FormatInt(n, 10)
	}

	if res.Body != nil {
		buf.Reset()
		n, _ := io.Copy(buf, res.Body) // nolint: gas
		sdreq.ResponseSize = strconv.FormatInt(n, 10)
	}

	return sdreq
}

// MarshalLogObject implements zapcore.ObjectMarshaller interface.
func (req HTTPPayload) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if req.RequestMethod != "" {
		enc.AddString("requestMethod", req.RequestMethod)
	}
	if req.RequestURL != "" {
		enc.AddString("requestUrl", req.RequestURL)
	}
	if req.RequestSize != "" {
		enc.AddString("requestSize", req.RequestSize)
	}
	if req.Status != 0 {
		enc.AddInt("status", req.Status)
	}
	if req.ResponseSize != "" {
		enc.AddString("responseSize", req.ResponseSize)
	}
	if req.UserAgent != "" {
		enc.AddString("userAgent", req.UserAgent)
	}
	if req.RemoteIP != "" {
		enc.AddString("remoteIp", req.RemoteIP)
	}
	if req.ServerIP != "" {
		enc.AddString("serverIp", req.ServerIP)
	}
	if req.Referer != "" {
		enc.AddString("referer", req.Referer)
	}
	if req.Latency != "" {
		enc.AddString("latency", req.Latency)
	}
	if req.CacheFillBytes != "" {
		enc.AddString("cacheFillBytes", req.CacheFillBytes)
	}
	if req.Protocol != "" {
		enc.AddString("protocol", req.Protocol)
	}
	enc.AddBool("cacheLookup", req.CacheLookup)
	enc.AddBool("cacheHit", req.CacheHit)
	enc.AddBool("cacheValidatedWithOriginServer", req.CacheValidatedWithOriginServer)
	return nil
}
