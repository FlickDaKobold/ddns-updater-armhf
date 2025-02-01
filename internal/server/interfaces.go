package server

import (
	"context"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/records"
)

type Database interface {
	SelectAll() (records []records.Record)
}

type UpdateForcer interface {
	ForceUpdate(ctx context.Context) (errors []error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
