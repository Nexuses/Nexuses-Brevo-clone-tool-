package core

import (
	"net/http"

	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

// CampaignReportData is the aggregated dataset for a campaign XLSX report.
type CampaignReportData struct {
	Campaign     models.Campaign
	UniqueCounts models.CampaignReportUniqueCounts
	UnsubCount   int
	Audience     []models.CampaignReportAudience
	ViewStats    []models.CampaignReportViewStats
	ClickStats   []models.CampaignReportClickStats
	Bounces      []models.CampaignReportBounce
}

// GetCampaignReportData loads all rows required to build a campaign report.
func (c *Core) GetCampaignReportData(id int) (CampaignReportData, error) {
	var out CampaignReportData

	camp, err := c.GetCampaign(id, "", "")
	if err != nil {
		return out, err
	}
	out.Campaign = camp

	if err := c.q.GetCampaignReportUniqueCounts.Get(&out.UniqueCounts, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}
	if err := c.q.GetCampaignReportUnsubCount.Get(&out.UnsubCount, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}
	if err := c.q.GetCampaignReportAudience.Select(&out.Audience, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}
	if err := c.q.GetCampaignReportViewStats.Select(&out.ViewStats, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}
	if err := c.q.GetCampaignReportClickStats.Select(&out.ClickStats, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}
	if err := c.q.GetCampaignReportBounces.Select(&out.Bounces, id); err != nil {
		return out, echo.NewHTTPError(http.StatusInternalServerError,
			c.i18n.Ts("globals.messages.errorFetching", "name", "{globals.terms.campaign}", "error", pqErrMsg(err)))
	}

	return out, nil
}
