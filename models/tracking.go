package models

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

// CampaignTrackingConfig holds per-campaign tracking overrides and UTM parameters.
// Nil pointer fields inherit site-wide privacy settings.
type CampaignTrackingConfig struct {
	Enabled            *bool             `json:"enabled"`
	TrackViews         *bool             `json:"track_views"`
	TrackClicks        *bool             `json:"track_clicks"`
	IndividualTracking *bool             `json:"individual_tracking"`
	Notes              string            `json:"notes"`
	UTMSource          string            `json:"utm_source"`
	UTMMedium          string            `json:"utm_medium"`
	UTMCampaign        string            `json:"utm_campaign"`
	UTMContent         string            `json:"utm_content"`
	UTMTerm            string            `json:"utm_term"`
	CustomParams       map[string]string `json:"custom_params"`
}

// CampaignAnalyticsRecord is a single view or click event with subscriber details.
type CampaignAnalyticsRecord struct {
	CampaignID     int    `db:"campaign_id" json:"campaign_id"`
	CampaignName   string `db:"campaign_name" json:"campaign_name"`
	SubscriberID   int    `db:"subscriber_id" json:"subscriber_id"`
	SubscriberUUID string `db:"subscriber_uuid" json:"subscriber_uuid"`
	Email          string `db:"email" json:"email"`
	Name           string `db:"name" json:"name"`
	URL            string `db:"url" json:"url,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// ParseTrackingConfig unmarshals the campaign tracking_config JSON column.
func (c *Campaign) ParseTrackingConfig() CampaignTrackingConfig {
	var tc CampaignTrackingConfig
	if len(c.TrackingConfig) == 0 {
		return tc
	}
	_ = json.Unmarshal(c.TrackingConfig, &tc)
	return tc
}

// TrackingDisabled returns true if tracking should not run for this campaign.
func (c *Campaign) TrackingDisabled(globalDisabled bool) bool {
	if globalDisabled {
		return true
	}
	tc := c.ParseTrackingConfig()
	if tc.Enabled != nil && !*tc.Enabled {
		return true
	}
	return false
}

// TrackViews returns whether open/view tracking is enabled for this campaign.
func (c *Campaign) TrackViews(globalDisabled bool) bool {
	if c.TrackingDisabled(globalDisabled) {
		return false
	}
	tc := c.ParseTrackingConfig()
	if tc.TrackViews != nil {
		return *tc.TrackViews
	}
	return true
}

// TrackClicks returns whether link click tracking is enabled for this campaign.
func (c *Campaign) TrackClicks(globalDisabled bool) bool {
	if c.TrackingDisabled(globalDisabled) {
		return false
	}
	tc := c.ParseTrackingConfig()
	if tc.TrackClicks != nil {
		return *tc.TrackClicks
	}
	return true
}

// UseIndividualTracking resolves per-campaign vs site-wide individual tracking.
func (c *Campaign) UseIndividualTracking(globalIndividual bool) bool {
	tc := c.ParseTrackingConfig()
	if tc.IndividualTracking != nil {
		return *tc.IndividualTracking
	}
	return globalIndividual
}

// ApplyTrackingParams appends UTM and custom query params from the campaign config to a URL.
func (c *Campaign) ApplyTrackingParams(rawURL string) string {
	tc := c.ParseTrackingConfig()
	if rawURL == "" {
		return rawURL
	}

	u, err := url.Parse(strings.ReplaceAll(rawURL, "&amp;", "&"))
	if err != nil {
		return rawURL
	}

	q := u.Query()
	setIf := func(k, v string) {
		if v != "" && q.Get(k) == "" {
			q.Set(k, v)
		}
	}
	setIf("utm_source", tc.UTMSource)
	setIf("utm_medium", tc.UTMMedium)
	setIf("utm_campaign", tc.UTMCampaign)
	setIf("utm_content", tc.UTMContent)
	setIf("utm_term", tc.UTMTerm)
	for k, v := range tc.CustomParams {
		if k != "" && v != "" && q.Get(k) == "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
