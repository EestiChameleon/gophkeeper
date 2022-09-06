package responses

import (
	"net/http"
	"time"
)

const (
	// HEADERS ------------------------------------
	HeaderContentType     = "Content-Type"
	HeaderLocation        = "Location"
	HeaderXForwardedFor   = "X-Forwarded-For"
	HeaderXRealIP         = "X-Real-IP"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"

	// CONTENT TYPE -------------------------------------------------------------
	charsetUTF8 = "charset=utf-8"

	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationXML             = "application/xml"
	MIMEApplicationXMLCharsetUTF8  = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                    = "text/xml"
	MIMETextXMLCharsetUTF8         = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm            = "application/x-www-form-urlencoded"
	MIMETextPlain                  = "text/plain"
	MIMETextPlainCharsetUTF8       = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm              = "multipart/form-data"
)

// CreateCookie func provides a cookie "key=value" based on given params.
func CreateCookie(key string, value string) *http.Cookie {
	return &http.Cookie{
		Name:    key,
		Value:   value,
		Path:    "/",
		Expires: time.Now().Add(time.Second * 60 * 60),
	}
}
