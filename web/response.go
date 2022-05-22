// Package web provides commons for interacting with 3rd party web services such as those in the broker package.
package web

import (
	"net/http"
	"time"
)

// Response wraps a http response returned by a 3rd party service.
type Response struct {
	Resp *http.Response
	Meta ResponseMetadata
}

// ResponseMetadata contains metadata about the response.
type ResponseMetadata struct {
	Page Page
	Rate Rate
}

// Page contains pagination information about the response.
type Page struct {
	PageSize   int
	PageNum    int
	PagesTotal int
}

// Rate contains information about the current rate limit status.
type Rate struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
}

// ListOpts is used to specify pagination options in a web call.
type ListOpts struct {
	PageSize int
	PageNum  int
}
