package health

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/constants"
)

func MakeIsHealthy(db AllSelecter, resolver LookupIPer) func(ctx context.Context) error {
	return func(ctx context.Context) (err error) {
		return isHealthy(ctx, db, resolver)
	}
}

var (
	ErrRecordUpdateFailed = errors.New("record update failed")
	ErrRecordIPNotSet     = errors.New("record IP not set")
	ErrLookupMismatch     = errors.New("lookup IP addresses do not match")
)

// isHealthy checks all the records were updated successfully and returns an error if not.
func isHealthy(ctx context.Context, db AllSelecter, resolver LookupIPer) (err error) {
	records := db.SelectAll()
	for _, record := range records {
		if record.Status == constants.FAIL {
			return fmt.Errorf("%w: %s", ErrRecordUpdateFailed, record.String())
		} else if record.Provider.Proxied() {
			continue
		}

		hostname := record.Provider.BuildDomainName()

		currentIP := record.History.GetCurrentIP()
		if !currentIP.IsValid() {
			return fmt.Errorf("%w: for hostname %s", ErrRecordIPNotSet, hostname)
		}

		lookedUpNetIPs, err := resolver.LookupIP(ctx, "ip", hostname)
		if err != nil {
			return err
		}

		found := false
		lookedUpIPsString := make([]string, len(lookedUpNetIPs))
		for i, netIP := range lookedUpNetIPs {
			var ip netip.Addr
			switch {
			case netIP == nil:
			case netIP.To4() != nil:
				ip = netip.AddrFrom4([4]byte(netIP.To4()))
			default: // IPv6
				ip = netip.AddrFrom16([16]byte(netIP.To16()))
			}
			if ip.Compare(currentIP) == 0 {
				found = true
				break
			}
			lookedUpIPsString[i] = ip.String()
		}
		if !found {
			return fmt.Errorf("%w: %s instead of %s for %s",
				ErrLookupMismatch, strings.Join(lookedUpIPsString, ","), currentIP, hostname)
		}
	}
	return nil
}
