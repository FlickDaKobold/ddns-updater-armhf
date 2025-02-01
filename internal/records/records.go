package records

import (
	"fmt"
	"time"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/constants"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/models"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider"
)

// Record contains all the information to update and display a DNS record.
type Record struct { // internal
	Provider provider.Provider // fixed
	History  models.History    // past information
	Status   models.Status
	Message  string
	Time     time.Time
	LastBan  *time.Time // nil means no last ban
}

// New returns a new Record with provider and some history.
func New(provider provider.Provider, events []models.HistoryEvent) Record {
	return Record{
		Provider: provider,
		History:  events,
		Status:   constants.UNSET,
	}
}

func (r *Record) String() string {
	status := string(r.Status)
	if r.Message != "" {
		status += " (" + r.Message + ")"
	}
	return fmt.Sprintf("%s: %s %s; %s",
		r.Provider, status, r.Time.Format("2006-01-02 15:04:05 MST"), r.History)
}
