package utils

import (
	"github.com/FlickDaKobold/ddns-updater-armhf/internal/models"
	"github.com/FlickDaKobold/ddns-updater-armhf/pkg/publicip/ipversion"
)

func ToString(domain, owner string, provider models.Provider, ipVersion ipversion.IPVersion) string {
	return "[domain: " + domain + " | owner: " + owner + " | provider: " +
		string(provider) + " | ip: " + ipVersion.String() + "]"
}
