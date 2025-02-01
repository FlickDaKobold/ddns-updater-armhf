package dreamhost

import (
	"net/http"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/headers"
)

func setHeaders(request *http.Request) {
	headers.SetUserAgent(request)
	headers.SetAccept(request, "application/json")
}
