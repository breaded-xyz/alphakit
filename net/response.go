package net

import (
	"net/http"
	"time"
)

// Response wraps a http response returned by a 3rd party service.
type Response struct {
	NetResponse *http.Response
	Meta        ResponseMetadata
}

type ResponseMetadata struct {
	Page Page
	Rate Rate
}

type Page struct {
	PageSize   int
	PageNum    int
	PagesTotal int
}

type Rate struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
}

type ListOpts struct {
	PageSize int
	PageNum  int
}
