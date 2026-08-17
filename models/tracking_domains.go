package models

import null "gopkg.in/volatiletech/null.v6"

// Custom tracking domain statuses.
const (
	TrackingDomainStatusPending  = "pending"
	TrackingDomainStatusVerified = "verified"
	TrackingDomainStatusFailed   = "failed"
)

// TrackingDomain is a user-owned custom click/open tracking hostname.
type TrackingDomain struct {
	Base

	UserID         int       `db:"user_id" json:"user_id"`
	Domain         string    `db:"domain" json:"domain"`
	Status         string    `db:"status" json:"status"`
	DNSRecordType  string    `db:"dns_record_type" json:"dns_record_type"`
	DNSRecordName  string    `db:"dns_record_name" json:"dns_record_name"`
	DNSRecordValue string    `db:"dns_record_value" json:"dns_record_value"`
	VerifiedAt     null.Time `db:"verified_at" json:"verified_at"`
	LastError      string    `db:"last_error" json:"last_error"`
}
