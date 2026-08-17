package manager

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"github.com/knadh/listmonk/models"
	"golang.org/x/net/html"
)

// rewriteMessageLinks converts ordinary HTTP(S) anchors into the existing
// listmonk tracking-link format after the template has been rendered for a
// subscriber. Explicit {{ TrackLink }} URLs are already internal tracking URLs
// at this point and are intentionally left untouched.
func (m *Manager) rewriteMessageLinks(msg *CampaignMessage) {
	if msg == nil || msg.Campaign == nil ||
		msg.Campaign.ContentType == models.CampaignContentTypePlain ||
		m.cfg.DisableTracking || !msg.Campaign.TrackClicks(m.cfg.DisableTracking) {
		return
	}

	subUUID := msg.Subscriber.UUID
	if !msg.Campaign.UseIndividualTracking(m.cfg.IndividualTracking) {
		subUUID = dummyUUID
	}

	body, err := rewriteHTMLAnchors(msg.body, func(raw string) string {
		if !m.isTrackableHref(raw, msg.Campaign) {
			return raw
		}
		target := msg.Campaign.ApplyTrackingParams(raw)
		return m.trackLink(target, msg.Campaign, subUUID)
	})
	if err != nil {
		m.log.Printf("error rewriting campaign links: %v", err)
		return
	}
	msg.body = body
}

// rewriteHTMLAnchors preserves the original HTML byte-for-byte except for start
// tags of anchors whose href is replaced. Using the tokenizer avoids rewriting
// CSS, text, images, or unrelated attributes.
func rewriteHTMLAnchors(body []byte, rewrite func(string) string) ([]byte, error) {
	z := html.NewTokenizer(bytes.NewReader(body))
	var out bytes.Buffer
	out.Grow(len(body))

	for {
		switch z.Next() {
		case html.ErrorToken:
			if err := z.Err(); err != nil && err != io.EOF {
				return nil, err
			}
			return out.Bytes(), nil
		case html.StartTagToken:
			tok := z.Token()
			if tok.Data != "a" {
				out.Write(z.Raw())
				continue
			}
			for i := range tok.Attr {
				if tok.Attr[i].Namespace == "" && tok.Attr[i].Key == "href" {
					tok.Attr[i].Val = rewrite(tok.Attr[i].Val)
				}
			}
			out.WriteString(tok.String())
		default:
			out.Write(z.Raw())
		}
	}
}

func (m *Manager) isTrackableHref(raw string, campaign *models.Campaign) bool {
	value := strings.TrimSpace(raw)
	if value == "" || strings.HasPrefix(value, "#") {
		return false
	}

	u, err := url.Parse(value)
	if err != nil || u.Host == "" || u.User != nil ||
		(u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	// Never wrap existing tracking, open, archive, or subscription URLs.
	switch {
	case strings.HasPrefix(u.Path, "/link/"),
		strings.HasPrefix(u.Path, "/campaign/"),
		strings.HasPrefix(u.Path, "/subscription/"):
		return false
	}

	internalHosts := []string{m.trackingBase(campaign), m.cfg.TrackingURL, m.cfg.RootURL}
	for _, base := range internalHosts {
		parsed, parseErr := url.Parse(base)
		if parseErr == nil && parsed.Hostname() != "" &&
			strings.EqualFold(parsed.Hostname(), u.Hostname()) &&
			strings.HasPrefix(u.Path, "/archive") {
			return false
		}
	}

	return true
}
