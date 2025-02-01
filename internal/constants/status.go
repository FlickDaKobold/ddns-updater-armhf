package constants

import "github.com/FlickDaKobold/ddns-updater-armhf/internal/models"

const (
	FAIL     models.Status = "failure"
	SUCCESS  models.Status = "success"
	UPTODATE models.Status = "up to date"
	UPDATING models.Status = "updating"
	UNSET    models.Status = "unset"
)
