package core

import (
	"database/sql"
	"net/http"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// GetTrackingDomains returns custom tracking domains owned by userID.
func (c *Core) GetTrackingDomains(userID int) ([]models.TrackingDomain, error) {
	out := []models.TrackingDomain{}
	if err := c.q.GetTrackingDomains.Select(&out, userID); err != nil {
		c.log.Printf("error fetching tracking domains: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return out, nil
}

// GetTrackingDomain returns a single tracking domain owned by userID.
func (c *Core) GetTrackingDomain(id, userID int) (models.TrackingDomain, error) {
	var out models.TrackingDomain
	if err := c.q.GetTrackingDomain.Get(&out, id, userID); err != nil {
		if err == sql.ErrNoRows {
			return out, echo.NewHTTPError(http.StatusNotFound,
				c.i18n.Ts("globals.messages.notFound", "name", "tracking domain"))
		}
		c.log.Printf("error fetching tracking domain: %v", err)
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return out, nil
}

// CreateTrackingDomain creates a pending CTD for userID. Body user_id must never be trusted.
func (c *Core) CreateTrackingDomain(userID int, domain, baseDomain, dnsName, dnsValue string) (models.TrackingDomain, error) {
	var out models.TrackingDomain
	if err := c.q.CreateTrackingDomain.Get(&out, userID, domain, baseDomain, dnsName, dnsValue); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Constraint == "custom_tracking_domains_domain_key" {
			return out, echo.NewHTTPError(http.StatusConflict, "tracking domain already exists")
		}
		c.log.Printf("error creating tracking domain: %v", err)
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorCreating", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return out, nil
}

// UpdateTrackingDomainHost sets the tracking hostname before DNS verification.
func (c *Core) UpdateTrackingDomainHost(id, userID int, domain, dnsName string) (models.TrackingDomain, error) {
	var out models.TrackingDomain
	if err := c.q.UpdateTrackingDomainHost.Get(&out, id, userID, domain, dnsName); err != nil {
		if err == sql.ErrNoRows {
			return out, echo.NewHTTPError(http.StatusNotFound,
				c.i18n.Ts("globals.messages.notFound", "name", "tracking domain"))
		}
		c.log.Printf("error updating tracking domain host: %v", err)
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorUpdating", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return out, nil
}

// UpdateTrackingDomainVerification updates verification status for a user-owned CTD.
func (c *Core) UpdateTrackingDomainVerification(id, userID int, status, lastError, dnsValue string) (models.TrackingDomain, error) {
	var out models.TrackingDomain
	if err := c.q.UpdateTrackingDomainVerification.Get(&out, id, userID, status, lastError, dnsValue); err != nil {
		if err == sql.ErrNoRows {
			return out, echo.NewHTTPError(http.StatusNotFound,
				c.i18n.Ts("globals.messages.notFound", "name", "tracking domain"))
		}
		c.log.Printf("error updating tracking domain verification: %v", err)
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorUpdating", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return out, nil
}

// DeleteTrackingDomain deletes a user-owned CTD if not referenced by active campaigns.
func (c *Core) DeleteTrackingDomain(id, userID int) error {
	var statuses []string
	if err := c.q.GetTrackingDomainBlockingCampaignStatuses.Select(&statuses, id, userID); err != nil {
		c.log.Printf("error checking tracking domain campaign refs: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	if len(statuses) > 0 {
		return echo.NewHTTPError(http.StatusConflict,
			"tracking domain is used by a running, scheduled, or paused campaign")
	}

	res, err := c.q.DeleteTrackingDomain.Exec(id, userID)
	if err != nil {
		c.log.Printf("error deleting tracking domain: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorDeleting", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return echo.NewHTTPError(http.StatusNotFound,
			c.i18n.Ts("globals.messages.notFound", "name", "tracking domain"))
	}
	return nil
}

// GetLatestVerifiedTrackingDomainID returns the user's most recently verified CTD id, or 0.
func (c *Core) GetLatestVerifiedTrackingDomainID(userID int) (int, error) {
	var id int
	if err := c.q.GetLatestVerifiedTrackingDomainID.Get(&id, userID); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		c.log.Printf("error fetching latest verified tracking domain: %v", err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "tracking domain", "error", pqErrMsg(err)))
	}
	return id, nil
}
