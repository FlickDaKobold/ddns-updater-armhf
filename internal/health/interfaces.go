package health

import (
	"context"
	"net"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/records"
)

type AllSelecter interface {
	SelectAll() (records []records.Record)
}

type LookupIPer interface {
	LookupIP(ctx context.Context, network, host string) (ips []net.IP, err error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
