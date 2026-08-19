-- custom tracking domains
-- name: get-tracking-domains
SELECT id, user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value,
       verified_at, last_error, created_at, updated_at
FROM custom_tracking_domains
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: get-tracking-domain
SELECT id, user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value,
       verified_at, last_error, created_at, updated_at
FROM custom_tracking_domains
WHERE id = $1 AND user_id = $2;

-- name: create-tracking-domain
INSERT INTO custom_tracking_domains
    (user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value)
VALUES ($1, $2, $3, 'pending', 'CNAME', $4, $5)
RETURNING id, user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value,
          verified_at, last_error, created_at, updated_at;

-- name: update-tracking-domain-host
UPDATE custom_tracking_domains SET
    domain = $3,
    dns_record_name = $4,
    status = 'pending',
    verified_at = NULL,
    last_error = '',
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND status IN ('pending', 'failed')
RETURNING id, user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value,
          verified_at, last_error, created_at, updated_at;

-- name: update-tracking-domain-verification
UPDATE custom_tracking_domains SET
    status = $3,
    verified_at = (CASE WHEN $3 = 'verified' THEN NOW() ELSE NULL END),
    last_error = $4,
    dns_record_value = COALESCE(NULLIF($5, ''), dns_record_value),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, domain, base_domain, status, dns_record_type, dns_record_name, dns_record_value,
          verified_at, last_error, created_at, updated_at;

-- name: delete-tracking-domain
DELETE FROM custom_tracking_domains WHERE id = $1 AND user_id = $2;

-- name: get-tracking-domain-blocking-campaign-statuses
SELECT DISTINCT status FROM campaigns
WHERE tracking_domain_id = $1
  AND EXISTS (
      SELECT 1 FROM custom_tracking_domains
      WHERE id = $1 AND user_id = $2
  )
  AND status = ANY('{running,scheduled,paused}'::campaign_status[]);

-- name: get-latest-verified-tracking-domain-id
SELECT id FROM custom_tracking_domains
WHERE user_id = $1 AND status = 'verified'
ORDER BY verified_at DESC NULLS LAST, id DESC
LIMIT 1;
