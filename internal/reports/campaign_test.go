package reports

import (
	"testing"
	"time"

	"github.com/knadh/listmonk/internal/core"
	"github.com/knadh/listmonk/models"
	null "gopkg.in/volatiletech/null.v6"
)

func TestGenerateCampaignXLSX(t *testing.T) {
	started := time.Date(2026, 7, 10, 23, 21, 0, 0, time.UTC)
	data := core.CampaignReportData{
		Campaign: models.Campaign{
			Base: models.Base{ID: 1, CreatedAt: null.TimeFrom(started), UpdatedAt: null.TimeFrom(started)},
			CampaignMeta: models.CampaignMeta{
				StartedAt: null.TimeFrom(started),
				ToSend:    3,
				Sent:      3,
				Views:     3,
				Clicks:    0,
				Bounces:   1,
			},
			Name: "Test Campaign",
		},
		UniqueCounts: models.CampaignReportUniqueCounts{UniqueOpens: 3, UniqueClicks: 0, Bounces: 1},
		Audience: []models.CampaignReportAudience{
			{ID: 1, Name: "Alice", Email: "alice@nexuses.in", Attribs: models.JSON{"company": "Nex"}},
		},
		ViewStats: []models.CampaignReportViewStats{
			{SubscriberID: 1, Opens: 1, FirstOpenAt: started.Add(3 * time.Minute)},
		},
		Bounces: []models.CampaignReportBounce{
			{Type: "hard", Source: "smtp", CreatedAt: started, Name: "Bob", Email: "bob@nexuses.in"},
		},
	}

	out, err := GenerateCampaignXLSX(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) < 1000 {
		t.Fatalf("expected non-trivial xlsx output, got %d bytes", len(out))
	}
}
