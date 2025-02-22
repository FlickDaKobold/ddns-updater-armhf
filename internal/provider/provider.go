package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/FlickDaKobold/ddns-updater-armhf/internal/models"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/constants"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/aliyun"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/allinkl"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/changeip"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/cloudflare"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/custom"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dd24"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/ddnss"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/desec"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/digitalocean"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dnsomatic"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dnspod"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/domeneshop"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dondominio"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dreamhost"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/duckdns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dyn"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dynu"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/dynv6"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/easydns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/example"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/freedns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/gandi"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/gcp"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/godaddy"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/goip"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/he"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/hetzner"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/infomaniak"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/inwx"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/ionos"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/linode"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/loopia"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/luadns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/myaddr"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/namecheap"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/namecom"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/namesilo"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/netcup"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/njalla"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/noip"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/nowdns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/opendns"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/ovh"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/porkbun"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/route53"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/selfhostde"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/servercow"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/spdyn"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/strato"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/variomedia"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/vultr"
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/provider/providers/zoneedit"
	"github.com/FlickDaKobold/ddns-updater-armhf/pkg/publicip/ipversion"
)

type Provider interface {
	String() string
	Domain() string
	Owner() string
	BuildDomainName() string
	HTML() models.HTMLRow
	Proxied() bool
	IPVersion() ipversion.IPVersion
	IPv6Suffix() netip.Prefix
	Update(ctx context.Context, client *http.Client, ip netip.Addr) (newIP netip.Addr, err error)
}

var ErrProviderUnknown = errors.New("unknown provider")

//nolint:gocyclo
func New(providerName models.Provider, data json.RawMessage, domain, owner string, //nolint:ireturn
	ipVersion ipversion.IPVersion, ipv6Suffix netip.Prefix,
) (provider Provider, err error) {
	switch providerName {
	case constants.Aliyun:
		return aliyun.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.AllInkl:
		return allinkl.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Changeip:
		return changeip.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Cloudflare:
		return cloudflare.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Custom:
		return custom.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Dd24:
		return dd24.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DdnssDe:
		return ddnss.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DeSEC:
		return desec.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DigitalOcean:
		return digitalocean.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DNSOMatic:
		return dnsomatic.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DNSPod:
		return dnspod.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Domeneshop:
		return domeneshop.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DonDominio:
		return dondominio.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Dreamhost:
		return dreamhost.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DuckDNS:
		return duckdns.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Dyn:
		return dyn.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Dynu:
		return dynu.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.DynV6:
		return dynv6.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.EasyDNS:
		return easydns.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Example:
		return example.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.FreeDNS:
		return freedns.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Gandi:
		return gandi.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.GCP:
		return gcp.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.GoDaddy:
		return godaddy.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.GoIP:
		return goip.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.HE:
		return he.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Hetzner:
		return hetzner.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Infomaniak:
		return infomaniak.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.INWX:
		return inwx.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Ionos:
		return ionos.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Linode:
		return linode.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Loopia:
		return loopia.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.LuaDNS:
		return luadns.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Myaddr:
		return myaddr.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Namecheap:
		return namecheap.New(data, domain, owner)
	case constants.NameCom:
		return namecom.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.NameSilo:
		return namesilo.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Netcup:
		return netcup.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Njalla:
		return njalla.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.NoIP:
		return noip.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.NowDNS:
		return nowdns.New(data, domain, ipVersion, ipv6Suffix)
	case constants.OpenDNS:
		return opendns.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.OVH:
		return ovh.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Porkbun:
		return porkbun.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Route53:
		return route53.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.SelfhostDe:
		return selfhostde.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Servercow:
		return servercow.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Spdyn:
		return spdyn.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Strato:
		return strato.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Variomedia:
		return variomedia.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Vultr:
		return vultr.New(data, domain, owner, ipVersion, ipv6Suffix)
	case constants.Zoneedit:
		return zoneedit.New(data, domain, owner, ipVersion, ipv6Suffix)
	default:
		return nil, fmt.Errorf("%w: %s", ErrProviderUnknown, providerName)
	}
}
