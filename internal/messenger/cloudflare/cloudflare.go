// Package cloudflare implements a messenger that sends e-mail via the
// Cloudflare Email Sending REST API.
package cloudflare

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/knadh/listmonk/models"
)

const (
	// MessengerName is the unique messenger identifier campaigns / tx messages use.
	MessengerName = "cloudflare"

	apiBase = "https://api.cloudflare.com/client/v4/accounts"
)

// Options represents Cloudflare Email Sending credentials and HTTP client options.
type Options struct {
	AccountID     string        `json:"account_id"`
	APIToken      string        `json:"api_token"`
	MaxConns      int           `json:"max_conns"`
	Timeout       time.Duration `json:"timeout"`
	MaxMsgRetries int           `json:"max_msg_retries"`
}

// Cloudflare is a messenger that sends mail via Cloudflare Email Sending.
type Cloudflare struct {
	o   Options
	c   *http.Client
	url string
}

// emailAddr is Cloudflare's REST address object (uses "address", not "email").
type emailAddr struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}

type sendRequest struct {
	From        any               `json:"from"`
	To          any               `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []attachment      `json:"attachments,omitempty"`
}

type attachment struct {
	Content     string `json:"content"`
	Filename    string `json:"filename"`
	Type        string `json:"type,omitempty"`
	Disposition string `json:"disposition,omitempty"`
}

type apiResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Result *struct {
		Delivered         []string `json:"delivered"`
		PermanentBounces  []string `json:"permanent_bounces"`
		Queued            []string `json:"queued"`
	} `json:"result"`
}

// New returns a new Cloudflare Email messenger.
func New(o Options) (*Cloudflare, error) {
	o.AccountID = strings.TrimSpace(o.AccountID)
	o.APIToken = strings.TrimSpace(o.APIToken)
	if o.AccountID == "" {
		return nil, fmt.Errorf("cloudflare account ID is required")
	}
	if o.APIToken == "" {
		return nil, fmt.Errorf("cloudflare API token is required")
	}
	if o.MaxConns < 1 {
		o.MaxConns = 10
	}
	if o.Timeout < time.Second {
		o.Timeout = 10 * time.Second
	}

	return &Cloudflare{
		o:   o,
		url: fmt.Sprintf("%s/%s/email/sending/send", apiBase, o.AccountID),
		c: &http.Client{
			Timeout: o.Timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   o.MaxConns,
				MaxConnsPerHost:       o.MaxConns,
				ResponseHeaderTimeout: o.Timeout,
				IdleConnTimeout:       o.Timeout,
			},
		},
	}, nil
}

// Name returns the messenger's name.
func (cf *Cloudflare) Name() string {
	return MessengerName
}

// Push sends a message via the Cloudflare Email Sending API.
func (cf *Cloudflare) Push(m models.Message) error {
	from, err := toCloudflareAddr(m.From)
	if err != nil {
		return fmt.Errorf("invalid from address %q: %v", m.From, err)
	}
	to, err := toCloudflareRecipients(m.To)
	if err != nil {
		return fmt.Errorf("invalid to address: %v", err)
	}

	req := sendRequest{
		From:    from,
		To:      to,
		Subject: m.Subject,
	}

	switch m.ContentType {
	case "plain":
		req.Text = string(m.Body)
	default:
		req.HTML = string(m.Body)
		if len(m.AltBody) > 0 {
			req.Text = string(m.AltBody)
		} else if req.HTML == "" {
			// Cloudflare requires at least one of text or html.
			req.Text = " "
		}
	}

	if len(m.Headers) > 0 {
		req.Headers = make(map[string]string, len(m.Headers))
		for k, v := range m.Headers {
			if len(v) == 0 {
				continue
			}
			req.Headers[k] = v[0]
		}
	}

	if len(m.Attachments) > 0 {
		req.Attachments = make([]attachment, 0, len(m.Attachments))
		for _, f := range m.Attachments {
			a := attachment{
				Content:     base64.StdEncoding.EncodeToString(f.Content),
				Filename:    f.Name,
				Disposition: "attachment",
			}
			if ct := f.Header.Get("Content-Type"); ct != "" {
				a.Type = strings.Split(ct, ";")[0]
			}
			req.Attachments = append(req.Attachments, a)
		}
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return cf.exec(b)
}

// Flush flushes the message queue.
func (cf *Cloudflare) Flush() error {
	return nil
}

// Close closes idle HTTP connections.
func (cf *Cloudflare) Close() error {
	cf.c.CloseIdleConnections()
	return nil
}

func (cf *Cloudflare) exec(body []byte) error {
	httpReq, err := http.NewRequest(http.MethodPost, cf.url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+cf.o.APIToken)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "listmonk")

	r, err := cf.c.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}()

	respBody, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading Cloudflare response: %v", err)
	}

	var resp apiResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("cloudflare API HTTP %d: %s", r.StatusCode, truncate(string(respBody), 200))
	}

	if r.StatusCode < 200 || r.StatusCode >= 300 || !resp.Success {
		msg := formatAPIErrors(resp, r.StatusCode)
		// Cloudflare's beta sometimes returns internal_server for bad senders;
		// hint the usual fix when that happens.
		if strings.Contains(msg, "internal_server") {
			msg += " (check that app From email uses an onboarded Cloudflare domain, e.g. you@yourdomain.com — not Name <email> placeholders)"
		}
		return fmt.Errorf("cloudflare API error: %s", msg)
	}

	if resp.Result != nil && len(resp.Result.PermanentBounces) > 0 {
		return fmt.Errorf("cloudflare permanent bounce: %s", strings.Join(resp.Result.PermanentBounces, ", "))
	}

	return nil
}

// toCloudflareAddr converts "Name <addr@host>" / bare addresses into the
// shape Cloudflare Email Sending expects.
func toCloudflareAddr(s string) (any, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty address")
	}
	a, err := mail.ParseAddress(s)
	if err != nil {
		// Fall back to treating the whole string as an address if it looks like one.
		if strings.Contains(s, "@") && !strings.ContainsAny(s, "<>") {
			return s, nil
		}
		return nil, err
	}
	if a.Name == "" {
		return a.Address, nil
	}
	return emailAddr{Address: a.Address, Name: a.Name}, nil
}

func toCloudflareRecipients(to []string) (any, error) {
	if len(to) == 0 {
		return nil, fmt.Errorf("no recipients")
	}
	out := make([]any, 0, len(to))
	for _, t := range to {
		addr, err := toCloudflareAddr(t)
		if err != nil {
			return nil, err
		}
		out = append(out, addr)
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return out, nil
}

func formatAPIErrors(resp apiResponse, status int) string {
	if len(resp.Errors) == 0 {
		return fmt.Sprintf("HTTP %d", status)
	}
	parts := make([]string, 0, len(resp.Errors))
	for _, e := range resp.Errors {
		if e.Message != "" {
			parts = append(parts, e.Message)
		} else {
			parts = append(parts, fmt.Sprintf("code %d", e.Code))
		}
	}
	return strings.Join(parts, "; ")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
