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

// CreateTrackingDomain registers a pending custom tracking domain for the authenticated user.
// Any user_id in the request body is ignored.
func (a *App) CreateTrackingDomain(c echo.Context) error {
	var req struct {
		Domain string `json:"domain"`
		UserID int    `json:"user_id"` // ignored; ownership comes from auth
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	domain, err := trackingdomain.NormalizeDomain(req.Domain)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := auth.GetUser(c)
	ownerID := trackingdomain.OwnerIDFromAuth(user.ID, req.UserID)

	expected := trackingdomain.ExpectedCNAMETarget(a.urlCfg.TrackingURL, a.urlCfg.RootURL)
	if expected == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"app.tracking_url (or app.root_url) must be configured as the CNAME target")
	}

	out, err := a.core.CreateTrackingDomain(ownerID, domain, domain, expected)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, okResp{out})
}

// VerifyTrackingDomain performs DNS CNAME verification for a user-owned domain.
func (a *App) VerifyTrackingDomain(c echo.Context) error {
	id := getID(c)
	user := auth.GetUser(c)

	td, err := a.core.GetTrackingDomain(id, user.ID)
	if err != nil {
		return err
	}

	expected := trackingdomain.ExpectedCNAMETarget(a.urlCfg.TrackingURL, a.urlCfg.RootURL)
	res := trackingdomain.VerifyCNAME(context.Background(), a.trackingDNS, td.Domain, expected)

	lastErr := ""
	if res.Status != trackingdomain.StatusVerified {
		lastErr = res.Message
	}

	out, err := a.core.UpdateTrackingDomainVerification(id, user.ID, res.Status, lastErr, expected)
	if err != nil {
		return err
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
