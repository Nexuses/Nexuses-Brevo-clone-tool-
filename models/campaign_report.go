package models

import (
	"encoding/json"
	"time"
)

// CampaignReportAudience is a subscriber targeted by a campaign's lists.
type CampaignReportAudience struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Email   string `db:"email" json:"email"`
	Attribs JSON   `db:"attribs" json:"attribs"`
}

// CampaignReportViewStats aggregates opens per subscriber for a campaign.
type CampaignReportViewStats struct {
	SubscriberID int       `db:"subscriber_id" json:"subscriber_id"`
	Opens        int       `db:"opens" json:"opens"`
	FirstOpenAt  time.Time `db:"first_open_at" json:"first_open_at"`
}

// CampaignReportClickStats aggregates clicks per subscriber for a campaign.
type CampaignReportClickStats struct {
	SubscriberID  int       `db:"subscriber_id" json:"subscriber_id"`
	Clicks        int       `db:"clicks" json:"clicks"`
	FirstClickAt  time.Time `db:"first_click_at" json:"first_click_at"`
	FirstClickURL string    `db:"first_click_url" json:"first_click_url"`
}

// CampaignReportBounce is a bounce row for the campaign report.
type CampaignReportBounce struct {
	Type      string          `db:"type" json:"type"`
	Source    string          `db:"source" json:"source"`
	Meta      json.RawMessage `db:"meta" json:"meta"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	Name      string          `db:"name" json:"name"`
	Email     string          `db:"email" json:"email"`
	Attribs   JSON            `db:"attribs" json:"attribs"`
}

// CampaignReportUniqueCounts holds unique engagement totals for summary metrics.
type CampaignReportUniqueCounts struct {
	UniqueOpens  int `db:"unique_opens" json:"unique_opens"`
	UniqueClicks int `db:"unique_clicks" json:"unique_clicks"`
	Bounces      int `db:"bounces" json:"bounces"`
}

// CampaignReportSummary is the JSON payload for the campaign report UI.
type CampaignReportSummary struct {
	UniqueOpens  int `db:"unique_opens" json:"unique_opens"`
	TotalOpens   int `db:"total_opens" json:"total_opens"`
	UniqueClicks int `db:"unique_clicks" json:"unique_clicks"`
	TotalClicks  int `db:"total_clicks" json:"total_clicks"`
	Unsubscribes int `db:"unsubscribes" json:"unsubscribes"`
	SoftBounces  int `db:"soft_bounces" json:"soft_bounces"`
	HardBounces  int `db:"hard_bounces" json:"hard_bounces"`
	Complaints   int `db:"complaints" json:"complaints"`
	Bounces      int `db:"bounces" json:"bounces"`
}

// CampaignReportBounceReason is an aggregated bounce reason row.
type CampaignReportBounceReason struct {
	Type   string `db:"type" json:"type"`
	Reason string `db:"reason" json:"reason"`
	Count  int    `db:"count" json:"count"`
}
