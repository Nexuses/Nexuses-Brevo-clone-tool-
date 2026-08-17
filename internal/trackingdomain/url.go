package trackingdomain

import (
	"errors"
	"net/url"
	"strings"

	"github.com/knadh/listmonk/models"
)

var ErrInvalidTrackingURL = errors.New("invalid platform tracking URL")

// NormalizeTrackingURL validates a platform tracking base. A bare hostname is
// accepted and promoted to HTTPS. Paths, credentials, queries, and fragments are
// rejected because tracking paths are appended by the sender.
func NormalizeTrackingURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}

	u, err := url.Parse(value)
	if err != nil || u.Hostname() == "" || u.User != nil {
		return "", ErrInvalidTrackingURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", ErrInvalidTrackingURL
	}
	if (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return "", ErrInvalidTrackingURL
	}

	return strings.TrimRight(u.Scheme+"://"+u.Host, "/"), nil
}

// ResolveTrackingBase returns the tracking base URL (no trailing slash).
// Priority: verified custom domain (https only) → app.tracking_url → app.root_url.
func ResolveTrackingBase(verifiedDomain, trackingURL, rootURL string) string {
	if d := CanonicalHost(verifiedDomain); d != "" {
		return "https://" + d
	}
	if t := strings.TrimSpace(trackingURL); t != "" {
		return strings.TrimRight(t, "/")
	}
	return strings.TrimRight(strings.TrimSpace(rootURL), "/")
}

// LinkTrackPath formats a click-tracking URL under base.
func LinkTrackPath(base, linkUUID, campUUID, subUUID string) string {
	return strings.TrimRight(base, "/") + "/link/" + linkUUID + "/" + campUUID + "/" + subUUID
}

// ViewTrackPath formats an open-tracking pixel URL under base.
func ViewTrackPath(base, campUUID, subUUID string) string {
	return strings.TrimRight(base, "/") + "/campaign/" + campUUID + "/" + subUUID + "/px.png"
}

// BlocksDelete reports whether any campaign status prevents CTD deletion.
func BlocksDelete(statuses []string) bool {
	for _, s := range statuses {
		switch s {
		case models.CampaignStatusRunning, models.CampaignStatusScheduled, models.CampaignStatusPaused:
			return true
		}
	}
	return false
}

// OwnerIDFromAuth returns the authenticated user id and ignores any body-supplied id.
// Callers must pass the session/API user id as authUserID.
func OwnerIDFromAuth(authUserID, _bodyUserID int) int {
	return authUserID
}
