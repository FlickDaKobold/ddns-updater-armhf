package namecom

import (
	"net/http"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/headers"
)

func setHeaders(request *http.Request) {
	headers.SetAccept(request, "application/json")
	headers.SetUserAgent(request)
}
