package state

var _httpHeaders httpHeaders
var _httpContentTypeValues httpContentTypeValues

func init() {
	_httpHeaders = newHttpHeaders()
	_httpContentTypeValues = newHttpContentTypeValues()
}

type HttpContentWrapper struct {
	ApplicationJson string
}

type httpContentTypeValues struct {
	ApplicationJson string
}

func newHttpContentTypeValues() httpContentTypeValues {
	return httpContentTypeValues{
		ApplicationJson: "application/json",
	}
}

func HttpContentTypeValues() httpContentTypeValues {
	return _httpContentTypeValues
}

type httpHeader string

func (h httpHeader) String() string {
	return string(h)
}

type httpHeaders struct {
	Authorization httpHeader
	ContentType   httpHeader
	StartTime     httpHeader
	RequestId     httpHeader
	PlatformType  httpHeader
	Platform      httpHeader
	Version       httpHeader
	CacheControl  httpHeader
	Accept        httpHeader
}

func newHttpHeaders() httpHeaders {
	return httpHeaders{
		Authorization: "Authorization",
		ContentType:   "Content-Type",
		PlatformType:  "Platform-Type",
		Platform:      "Platform",
		Version:       "Version",
		StartTime:     "Start-Time",
		RequestId:     "Request-Id",
		CacheControl:  "Cache-Control",
		Accept:        "Accept",
	}
}

func HttpHeaders() httpHeaders {
	return _httpHeaders
}
