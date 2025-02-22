package hetzner

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/models"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/constants"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/errors"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/utils"
	"github.com/FlickDaKobold/ddns-updater-armhf/pkg/publicip/ipversion"
)

type Provider struct {
	domain         string
	owner          string
	ipVersion      ipversion.IPVersion
	ipv6Suffix     netip.Prefix
	token          string
	zoneIdentifier string
	ttl            uint32
}

func New(data json.RawMessage, domain, owner string,
	ipVersion ipversion.IPVersion, ipv6Suffix netip.Prefix) (
	p *Provider, err error,
) {
	extraSettings := struct {
		Token          string `json:"token"`
		ZoneIdentifier string `json:"zone_identifier"`
		TTL            uint32 `json:"ttl"`
	}{}
	err = json.Unmarshal(data, &extraSettings)
	if err != nil {
		return nil, err
	}

	ttl := uint32(1)
	if extraSettings.TTL > 0 {
		ttl = extraSettings.TTL
	}

	err = validateSettings(domain, extraSettings.ZoneIdentifier, extraSettings.Token)
	if err != nil {
		return nil, fmt.Errorf("validating provider specific settings: %w", err)
	}

	return &Provider{
		domain:         domain,
		owner:          owner,
		ipVersion:      ipVersion,
		ipv6Suffix:     ipv6Suffix,
		token:          extraSettings.Token,
		zoneIdentifier: extraSettings.ZoneIdentifier,
		ttl:            ttl,
	}, nil
}

func validateSettings(domain, zoneIdentifier, token string) (err error) {
	err = utils.CheckDomain(domain)
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrDomainNotValid, err)
	}

	switch {
	case zoneIdentifier == "":
		return fmt.Errorf("%w", errors.ErrZoneIdentifierNotSet)
	case token == "":
		return fmt.Errorf("%w", errors.ErrTokenNotSet)
	}
	return nil
}

func (p *Provider) String() string {
	return utils.ToString(p.domain, p.owner, constants.Hetzner, p.ipVersion)
}

func (p *Provider) Domain() string {
	return p.domain
}

func (p *Provider) Owner() string {
	return p.owner
}

func (p *Provider) IPVersion() ipversion.IPVersion {
	return p.ipVersion
}

func (p *Provider) IPv6Suffix() netip.Prefix {
	return p.ipv6Suffix
}

func (p *Provider) Proxied() bool {
	return false
}

func (p *Provider) BuildDomainName() string {
	return utils.BuildDomainName(p.owner, p.domain)
}

func (p *Provider) HTML() models.HTMLRow {
	return models.HTMLRow{
		Domain:    fmt.Sprintf("<a href=\"http://%s\">%s</a>", p.BuildDomainName(), p.BuildDomainName()),
		Owner:     p.Owner(),
		Provider:  "<a href=\"https://www.hetzner.com\">Hetzner</a>",
		IPVersion: p.ipVersion.String(),
	}
}

func (p *Provider) Update(ctx context.Context, client *http.Client, ip netip.Addr) (newIP netip.Addr, err error) {
	recordID, upToDate, err := p.getRecordID(ctx, client, ip)
	switch {
	case stderrors.Is(err, errors.ErrReceivedNoResult):
		err = p.createRecord(ctx, client, ip)
		if err != nil {
			return netip.Addr{}, fmt.Errorf("creating record: %w", err)
		}
		return ip, nil
	case err != nil:
		return netip.Addr{}, fmt.Errorf("getting record id: %w", err)
	case upToDate:
		return ip, nil
	}

	ip, err = p.updateRecord(ctx, client, recordID, ip)
	if err != nil {
		return newIP, fmt.Errorf("updating record: %w", err)
	}

	return ip, nil
}
