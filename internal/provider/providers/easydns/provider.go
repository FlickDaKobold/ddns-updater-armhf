package easydns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strings"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/models"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/constants"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/errors"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/headers"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/utils"
	"github.com/FlickDaKobold/ddns-updater-armhf/pkg/publicip/ipversion"
)

type Provider struct {
	domain     string
	owner      string
	ipVersion  ipversion.IPVersion
	ipv6Suffix netip.Prefix
	username   string
	token      string
}

func New(data json.RawMessage, domain, owner string,
	ipVersion ipversion.IPVersion, ipv6Suffix netip.Prefix) (
	p *Provider, err error,
) {
	extraSettings := struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}{}
	err = json.Unmarshal(data, &extraSettings)
	if err != nil {
		return nil, err
	}

	err = validateSettings(domain, extraSettings.Username, extraSettings.Token)
	if err != nil {
		return nil, fmt.Errorf("validating provider specific settings: %w", err)
	}

	return &Provider{
		domain:     domain,
		owner:      owner,
		ipVersion:  ipVersion,
		ipv6Suffix: ipv6Suffix,
		username:   extraSettings.Username,
		token:      extraSettings.Token,
	}, nil
}

func validateSettings(domain, username, token string) (err error) {
	err = utils.CheckDomain(domain)
	if err != nil {
		return fmt.Errorf("%w: %w", errors.ErrDomainNotValid, err)
	}

	switch {
	case username == "":
		return fmt.Errorf("%w", errors.ErrUsernameNotSet)
	case token == "":
		return fmt.Errorf("%w", errors.ErrTokenNotSet)
	}
	return nil
}

func (p *Provider) String() string {
	return utils.ToString(p.domain, p.owner, constants.EasyDNS, p.ipVersion)
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
		Provider:  "<a href=\"https://easydns.com\">EasyDNS</a>",
		IPVersion: p.ipVersion.String(),
	}
}

func (p *Provider) Update(ctx context.Context, client *http.Client, ip netip.Addr) (
	newIP netip.Addr, err error,
) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.cp.easydns.com",
		Path:   "dyn/generic.php",
		User:   url.UserPassword(p.username, p.token),
	}
	values := url.Values{}
	values.Set("hostname", utils.BuildURLQueryHostname(p.owner, p.domain))
	values.Set("myip", ip.String())
	if p.owner == "*" {
		values.Set("wildcard", "ON")
	}
	u.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("creating http request: %w", err)
	}
	headers.SetUserAgent(request)

	response, err := client.Do(request)
	if err != nil {
		return netip.Addr{}, err
	}
	defer response.Body.Close()

	s, err := utils.ReadAndCleanBody(response.Body)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("reading response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return netip.Addr{}, fmt.Errorf("%w: %d: %s", errors.ErrHTTPStatusNotValid,
			response.StatusCode, utils.ToSingleLine(s))
	}

	switch {
	case s == "":
		return netip.Addr{}, fmt.Errorf("%w", errors.ErrReceivedNoResult)
	case strings.Contains(s, "no_service"):
		return netip.Addr{}, fmt.Errorf("%w", errors.ErrNoService)
	case strings.Contains(s, "no_access"):
		return netip.Addr{}, fmt.Errorf("%w", errors.ErrAuth)
	case strings.Contains(s, "illegal_input"), strings.Contains(s, "too_soon"):
		return netip.Addr{}, fmt.Errorf("%w", errors.ErrBannedAbuse)
	case strings.Contains(s, "no_error"), strings.Contains(s, "ok"):
		return ip, nil
	default:
		return netip.Addr{}, fmt.Errorf("%w: %s", errors.ErrUnknownResponse, utils.ToSingleLine(s))
	}
}
