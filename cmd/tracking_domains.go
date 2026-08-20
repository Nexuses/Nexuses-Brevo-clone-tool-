package main

import (
	"context"
	"net/http"

	"github.com/knadh/listmonk/internal/auth"
	"github.com/knadh/listmonk/internal/trackingdomain"
	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

type trackingDomainVerifyResp struct {
	models.TrackingDomain
	Verification trackingdomain.VerifyResult `json:"verification"`
	Message      string                      `json:"message"`
	Verified     bool                        `json:"verified"`
}

// GetTrackingDomains lists custom tracking domains for the authenticated user.
func (a *App) GetTrackingDomains(c echo.Context) error {
	user := auth.GetUser(c)
	out, err := a.core.GetTrackingDomains(user.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{out})
}

// CreateTrackingDomain registers a base domain and optional branded tracking subdomain.
// Any user_id in the request body is ignored.
func (a *App) CreateTrackingDomain(c echo.Context) error {
	var req struct {
		Domain       string `json:"domain"`
		TrackingHost string `json:"tracking_host"`
		UserID       int    `json:"user_id"` // ignored; ownership comes from auth
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	baseDomain, err := trackingdomain.NormalizeDomain(req.Domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	trackingHost := baseDomain
	if req.TrackingHost != "" {
		trackingHost, err = trackingdomain.NormalizeDomain(req.TrackingHost)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if !trackingdomain.IsUnderBase(trackingHost, baseDomain) {
			return echo.NewHTTPError(http.StatusBadRequest, "tracking host must belong to the registered domain")
		}
	}

	dnsName := trackingdomain.DNSHostLabel(trackingHost, baseDomain)

	user := auth.GetUser(c)
	ownerID := trackingdomain.OwnerIDFromAuth(user.ID, req.UserID)

	expected := trackingdomain.ExpectedCNAMETarget(a.urlCfg.TrackingURL, a.urlCfg.RootURL)
	if expected == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"app.tracking_url (or app.root_url) must be configured as the CNAME target")
	}

	out, err := a.core.CreateTrackingDomain(ownerID, trackingHost, baseDomain, dnsName, expected)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{out})
}

// VerifyTrackingDomain performs DNS CNAME verification for a user-owned domain.
// A tracking hostname can be supplied at verify time for base-only registrations.
func (a *App) VerifyTrackingDomain(c echo.Context) error {
	id := getID(c)
	user := auth.GetUser(c)

	var req struct {
		Domain string `json:"domain"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	td, err := a.core.GetTrackingDomain(id, user.ID)
	if err != nil {
		return err
	}

	baseDomain := td.BaseDomain
	if baseDomain == "" {
		baseDomain = td.Domain
	}

	verifyHost := td.Domain
	if req.Domain != "" && td.Status != models.TrackingDomainStatusVerified {
		newHost, err := trackingdomain.NormalizeDomain(req.Domain)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if !trackingdomain.IsUnderBase(newHost, baseDomain) {
			return echo.NewHTTPError(http.StatusBadRequest, "tracking host must belong to the registered domain")
		}
		if newHost == baseDomain {
			return echo.NewHTTPError(http.StatusBadRequest,
				"enter a tracking subdomain such as emailtrack."+baseDomain)
		}
		dnsName := trackingdomain.DNSHostLabel(newHost, baseDomain)
		td, err = a.core.UpdateTrackingDomainHost(id, user.ID, newHost, dnsName)
		if err != nil {
			return err
		}
		verifyHost = newHost
	} else if verifyHost == baseDomain && td.Status != models.TrackingDomainStatusVerified {
		return echo.NewHTTPError(http.StatusBadRequest,
			"enter a tracking subdomain such as emailtrack."+baseDomain+" before verifying DNS")
	}

	expected := trackingdomain.ExpectedCNAMETarget(a.urlCfg.TrackingURL, a.urlCfg.RootURL)
	res := trackingdomain.VerifyCNAME(context.Background(), a.trackingDNS, verifyHost, expected)

	lastErr := ""
	if res.Status != trackingdomain.StatusVerified {
		lastErr = res.Message
	}

	out, err := a.core.UpdateTrackingDomainVerification(id, user.ID, res.Status, lastErr, expected)
	if err != nil {
		return err
	}

	// Newly verified CTDs should apply to draft/scheduled campaigns that
	// were created before authentication completed.
	if res.Status == trackingdomain.StatusVerified {
		_ = a.core.AttachTrackingDomainToPendingCampaigns(out.ID)
	}

	return c.JSON(http.StatusOK, okResp{trackingDomainVerifyResp{
		TrackingDomain: out,
		Verification:   res,
		Message:        res.Message,
		Verified:       res.Status == trackingdomain.StatusVerified,
	}})
}

// DeleteTrackingDomain deletes a user-owned CTD when safe.
func (a *App) DeleteTrackingDomain(c echo.Context) error {
	id := getID(c)
	user := auth.GetUser(c)
	if err := a.core.DeleteTrackingDomain(id, user.ID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{true})
}
