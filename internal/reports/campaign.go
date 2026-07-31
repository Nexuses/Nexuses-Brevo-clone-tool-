// Package reports generates downloadable campaign result reports.
package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/knadh/listmonk/internal/core"
	"github.com/knadh/listmonk/models"
	"github.com/xuri/excelize/v2"
)

type subscriberEngagement struct {
	ID            int
	Name          string
	Email         string
	Company       string
	Domain        string
	Opens         int
	Clicks        int
	FirstOpenAt   time.Time
	FirstClickURL string
}

// GenerateCampaignXLSX builds a multi-sheet XLSX report modeled on campaign-report-template.xlsx.
func GenerateCampaignXLSX(data core.CampaignReportData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	eng := buildEngagement(data)
	reportDate := time.Now()

	if err := writeSummarySheet(f, data, reportDate); err != nil {
		return nil, err
	}
	if err := writeActivitySheet(f, data, eng, reportDate); err != nil {
		return nil, err
	}
	if err := writeBouncesSheet(f, data, reportDate); err != nil {
		return nil, err
	}

	f.DeleteSheet("Sheet1")

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSummarySheet(f *excelize.File, data core.CampaignReportData, reportDate time.Time) error {
	const sheet = "Summary"
	if _, err := f.NewSheet(sheet); err != nil {
		return err
	}

	camp := data.Campaign
	sent := camp.Sent
	if sent == 0 {
		sent = len(data.Audience)
	}
	openRate := pct(data.UniqueCounts.UniqueOpens, sent)
	clickRate := pct(data.UniqueCounts.UniqueClicks, sent)

	title := fmt.Sprintf("%s · Campaign Summary", camp.Name)
	rows := [][]string{
		{title},
		{fmt.Sprintf("Generated · %s", reportDate.Format("2 January 2006"))},
		{"Metric", "Value"},
		{"Campaign", camp.Name},
		{"Report Date", reportDate.Format("2 January 2006")},
		{"Total Contacts (audience)", fmt.Sprintf("%d", len(data.Audience))},
		{"Eligible Recipients", fmt.Sprintf("%d", camp.ToSend)},
		{"Unsubscribed at Launch", fmt.Sprintf("%d", data.UnsubCount)},
		{"Total Sent", fmt.Sprintf("%d", camp.Sent)},
		{"Unique Opens", fmt.Sprintf("%d", data.UniqueCounts.UniqueOpens)},
		{"Unique Clicks", fmt.Sprintf("%d", data.UniqueCounts.UniqueClicks)},
		{"Bounces", fmt.Sprintf("%d", data.UniqueCounts.Bounces)},
		{"Open Rate", openRate},
		{"Click Rate", clickRate},
	}

	for i, row := range rows {
		if err := f.SetSheetRow(sheet, fmt.Sprintf("A%d", i+1), &row); err != nil {
			return err
		}
	}
	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	_ = f.SetCellStyle(sheet, "A3", "B3", style)
	return nil
}

func writeActivitySheet(f *excelize.File, data core.CampaignReportData, eng map[int]subscriberEngagement, reportDate time.Time) error {
	const sheet = "Sent, Open, Click"
	if _, err := f.NewSheet(sheet); err != nil {
		return err
	}

	camp := data.Campaign
	sendAt := camp.StartedAt.Time
	if sendAt.IsZero() {
		sendAt = camp.UpdatedAt.Time
	}
	if sendAt.IsZero() {
		sendAt = camp.CreatedAt.Time
	}

	subtitle := fmt.Sprintf("%s · %s · %d Total Sent · %d Opens · %d Clicks",
		camp.Name,
		formatReportDate(reportDate),
		camp.Sent,
		data.UniqueCounts.UniqueOpens,
		data.UniqueCounts.UniqueClicks,
	)

	row := 1
	row = writeRow(f, sheet, row, []string{fmt.Sprintf("%s · Sent, Open & Click Tracker", camp.Name)})
	row = writeRow(f, sheet, row, []string{subtitle})
	row++
	row = writeRow(f, sheet, row, []string{"SENT — Campaign Audience"})
	row = writeRow(f, sheet, row, []string{"#", "Person Name", "Email Address", "Company", "Domain", "Step", "Send Date", "Send Time"})

	for i, sub := range data.Audience {
		e := eng[sub.ID]
		row = writeRow(f, sheet, row, []string{
			fmt.Sprintf("%d", i+1),
			pickName(sub.Name, sub.Email),
			sub.Email,
			pickCompany(sub.Attribs, e.Company),
			pickDomain(sub.Email, e.Domain),
			"Email",
			formatDate(sendAt),
			formatTime(sendAt),
		})
	}

	row++
	row = writeRow(f, sheet, row, []string{"CLICKERS — Opened & Clicked a Link"})
	row = writeRow(f, sheet, row, []string{"#", "Person Name", "Email Address", "Company", "Domain", "Open Date", "Open Time", "Opens", "Clicks", "CTA Clicked", "Status"})

	clickers := filterEngagement(eng, func(e subscriberEngagement) bool { return e.Clicks > 0 })
	if len(clickers) == 0 {
		row = writeRow(f, sheet, row, []string{"—", "No records", "", "", "", "", "", "0", "0", "—", "—"})
	} else {
		for i, e := range clickers {
			row = writeRow(f, sheet, row, []string{
				fmt.Sprintf("%d", i+1),
				e.Name,
				e.Email,
				e.Company,
				e.Domain,
				formatDate(e.FirstOpenAt),
				formatTime(e.FirstOpenAt),
				fmt.Sprintf("%d", e.Opens),
				fmt.Sprintf("%d", e.Clicks),
				shortURL(e.FirstClickURL),
				"Clicked",
			})
		}
	}

	row++
	row = writeRow(f, sheet, row, []string{"OPENS ONLY — Opened Email, No Click"})
	row = writeRow(f, sheet, row, []string{"#", "Person Name", "Email Address", "Company", "Domain", "Open Date", "Open Time", "Opens", "Clicks", "CTA Clicked", "Status"})

	opensOnly := filterEngagement(eng, func(e subscriberEngagement) bool { return e.Opens > 0 && e.Clicks == 0 })
	if len(opensOnly) == 0 {
		row = writeRow(f, sheet, row, []string{"—", "No records", "", "", "", "", "", "0", "0", "—", "—"})
	} else {
		for i, e := range opensOnly {
			row = writeRow(f, sheet, row, []string{
				fmt.Sprintf("%d", i+1),
				e.Name,
				e.Email,
				e.Company,
				e.Domain,
				formatDate(e.FirstOpenAt),
				formatTime(e.FirstOpenAt),
				fmt.Sprintf("%d", e.Opens),
				"0",
				"—",
				"Opened",
			})
		}
	}

	return nil
}

func writeBouncesSheet(f *excelize.File, data core.CampaignReportData, reportDate time.Time) error {
	const sheet = "Bounces & Unsubscribes"
	if _, err := f.NewSheet(sheet); err != nil {
		return err
	}

	camp := data.Campaign
	row := 1
	row = writeRow(f, sheet, row, []string{fmt.Sprintf("%s · Bounces & Unsubscribes", camp.Name)})
	row = writeRow(f, sheet, row, []string{
		fmt.Sprintf("%s · %s · %d Bounces · %d Unsubscribed (at launch)",
			camp.Name, formatReportDate(reportDate), data.UniqueCounts.Bounces, data.UnsubCount),
	})
	row++
	row = writeRow(f, sheet, row, []string{"BOUNCES"})
	row = writeRow(f, sheet, row, []string{"#", "Person Name", "Email Address", "Company", "Domain", "Reason", "Bounced Date", "Bounced Time"})

	if len(data.Bounces) == 0 {
		row = writeRow(f, sheet, row, []string{"—", "No records", "", "", "", "", "", ""})
	} else {
		for i, b := range data.Bounces {
			row = writeRow(f, sheet, row, []string{
				fmt.Sprintf("%d", i+1),
				pickName(b.Name, b.Email),
				b.Email,
				attribString(b.Attribs, "company"),
				emailDomain(b.Email),
				bounceReason(b),
				formatDate(b.CreatedAt),
				formatTime(b.CreatedAt),
			})
		}
	}

	row++
	row = writeRow(f, sheet, row, []string{"UNSUBSCRIBES (suppressed at launch)"})
	row = writeRow(f, sheet, row, []string{"Count", fmt.Sprintf("%d", data.UnsubCount)})
	row = writeRow(f, sheet, row, []string{"Contacts on the unsubscribe list were excluded when the campaign launched and were not emailed."})

	return nil
}

func buildEngagement(data core.CampaignReportData) map[int]subscriberEngagement {
	out := make(map[int]subscriberEngagement, len(data.Audience))
	for _, sub := range data.Audience {
		out[sub.ID] = subscriberEngagement{
			ID:      sub.ID,
			Name:    pickName(sub.Name, sub.Email),
			Email:   sub.Email,
			Company: attribString(sub.Attribs, "company"),
			Domain:  emailDomain(sub.Email),
		}
	}
	for _, v := range data.ViewStats {
		e := out[v.SubscriberID]
		e.Opens = v.Opens
		e.FirstOpenAt = v.FirstOpenAt
		out[v.SubscriberID] = e
	}
	for _, c := range data.ClickStats {
		e := out[c.SubscriberID]
		e.Clicks = c.Clicks
		if e.FirstOpenAt.IsZero() {
			e.FirstOpenAt = c.FirstClickAt
		}
		e.FirstClickURL = c.FirstClickURL
		out[c.SubscriberID] = e
	}
	return out
}

func filterEngagement(eng map[int]subscriberEngagement, fn func(subscriberEngagement) bool) []subscriberEngagement {
	var out []subscriberEngagement
	for _, e := range eng {
		if fn(e) {
			out = append(out, e)
		}
	}
	return out
}

func writeRow(f *excelize.File, sheet string, row int, values []string) int {
	_ = f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &values)
	return row + 1
}

func pickName(name, email string) string {
	if strings.TrimSpace(name) != "" {
		return name
	}
	return email
}

func pickCompany(attribs models.JSON, fallback string) string {
	if c := attribString(attribs, "company"); c != "" {
		return c
	}
	return fallback
}

func pickDomain(email, fallback string) string {
	if d := emailDomain(email); d != "" {
		return d
	}
	return fallback
}

func attribString(attribs models.JSON, key string) string {
	if attribs == nil {
		return ""
	}
	v, ok := attribs[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprintf("%v", t)
	}
}

func emailDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

func formatReportDate(t time.Time) string {
	return t.Format("2 January 2006")
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02-01-2006")
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("15:04:05")
}

func pct(num, den int) string {
	if den <= 0 {
		return "0%"
	}
	return fmt.Sprintf("%d%%", int(math.Round(float64(num)/float64(den)*100)))
}

func shortURL(u string) string {
	u = strings.TrimSpace(u)
	if u == "" {
		return "—"
	}
	if len(u) > 80 {
		return u[:77] + "..."
	}
	return u
}

func bounceReason(b models.CampaignReportBounce) string {
	if len(b.Meta) > 0 {
		var meta map[string]any
		if err := json.Unmarshal(b.Meta, &meta); err == nil {
			for _, key := range []string{"reason", "message", "description", "diagnostic", "error"} {
				if v, ok := meta[key]; ok {
					if s := strings.TrimSpace(fmt.Sprintf("%v", v)); s != "" {
						return s
					}
				}
			}
		}
	}
	if b.Source != "" {
		return b.Source
	}
	if b.Type != "" {
		return b.Type
	}
	return "bounce"
}
