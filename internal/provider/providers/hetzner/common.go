package hetzner

import (
	"net/http"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/headers"
)

func (p *Provider) setHeaders(request *http.Request) {
	headers.SetUserAgent(request)
	headers.SetContentType(request, "application/json")
	headers.SetAccept(request, "application/json")
	request.Header.Set("Auth-Api-Token", p.token)
}
