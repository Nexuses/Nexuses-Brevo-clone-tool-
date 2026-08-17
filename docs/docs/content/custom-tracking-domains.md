# Custom Tracking Domains

Custom Tracking Domains (CTDs) let campaigns use a customer hostname for click and open tracking (`/link/…`, `/campaign/…`) instead of the platform admin URL. Existing campaigns without a verified domain keep using the platform tracking base (or `app.root_url`).

## Platform tracking URL

Set the platform (default) tracking base URL:

| Setting | Example |
| ------- | ------- |
| `app.tracking_url` (Settings / DB) | `https://t.platform.example` |
| `LISTMONK_app__tracking_url` (env) | `https://t.platform.example` |
| `PLATFORM_TRACKING_DOMAIN` (env) | `t.platform.example` |

Rules:

- Use an absolute HTTPS URL with no trailing path, query, or fragment.
- `PLATFORM_TRACKING_DOMAIN` may be a bare hostname; listmonk promotes it to HTTPS.
- Environment values override the persisted `app.tracking_url`.
- If empty, listmonk falls back to `app.root_url` for tracking links and for the DNS CNAME target used when verifying customer domains.
- Keep `app.root_url` as the admin / brand site. Prefer a dedicated host for `app.tracking_url` when customers will CNAME to you.

Docker Compose publishes listmonk on port **9000** (`9000:9000`). Terminate TLS at a reverse proxy in front of that port; the stock compose file does not provide HTTPS.

## Customer DNS (CNAME)

For each customer tracking host (for example `track.customer.com`):

1. Create a **CNAME** whose target is the hostname from `app.tracking_url` (for example `t.platform.example`).
2. Confirm the record in Admin → Settings → Tracking domains (verify).

!!! note

    A CNAME only answers DNS. It does **not** issue or terminate TLS. HTTPS must still be configured on the platform edge for that customer hostname.

### Verification behaviour

Verification performs a real CNAME lookup against the expected host from `app.tracking_url`:

| Result | When |
| ------ | ---- |
| `pending` | Record missing, NXDOMAIN, timeout, or other temporary / propagation issues |
| `failed` | A CNAME exists but points at the wrong target |
| `verified` | Canonical CNAME matches the expected platform tracking host |

DNS can take minutes to hours to propagate. Retry verify after the record is visible with `dig` / `nslookup` from a public resolver.

## TLS for arbitrary customer domains

Browsers and many mail clients expect HTTPS on the hostname that appears in tracking URLs. The certificate **must match the customer hostname** (SAN or dedicated cert).

Recommended approach for unrelated customer zones:

- **Per-domain ACME / on-demand TLS** (Caddy `on_demand_tls`, Traefik, cert-manager, or similar) that obtains a certificate when a verified customer host first hits the edge.

Do **not** rely on a single wildcard like `*.platform.example` to cover `track.customer.com` in another DNS zone. Wildcards only help hosts under *your* zone (for example `*.t.platform.example`).

## Reverse proxy (tracking hosts only)

On customer tracking vhosts, expose only public tracking paths and block admin/API. Forward the original host and proxy headers to listmonk on port 9000.

Example Nginx server for a customer tracking host:

```nginx
# Customer tracking host → listmonk (Docker :9000)
# TLS termination assumed (ACME / on-demand) on this server block.

upstream listmonk_app {
    server 127.0.0.1:9000;
}

server {
    listen 443 ssl http2;
    server_name track.customer.com;

    # ssl_certificate / ssl_certificate_key (or ACME automation) go here.

    # Allow click + open tracking only.
    location /link/ {
        proxy_pass http://listmonk_app;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /campaign/ {
        proxy_pass http://listmonk_app;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Do not serve admin UI or private APIs on customer tracking hosts.
    location /admin {
        return 404;
    }

    location /api {
        return 404;
    }

    location / {
        return 404;
    }
}
```

Serve the admin UI and `/api` only on the platform `app.root_url` host, not on customer CNAMEs. See [HTTP routes](configuration.md#http-routes).

## Local testing

Useful for proxy and path checks without public DNS:

- **Hosts file**: map `track.customer.test` to `127.0.0.1`, then open or curl `http://track.customer.test:9000/…`.
- **curl Host header** (without editing hosts):

```shell
curl -sI -H "Host: track.customer.test" http://127.0.0.1:9000/link/
curl -sI -H "Host: track.customer.test" http://127.0.0.1:9000/admin/
```

!!! warning

    Hosts-file overrides and `curl -H "Host: …"` do **not** satisfy product DNS verification. Verification uses a real CNAME lookup (`LookupCNAME`). For local verify flows you need a resolvable CNAME (public DNS, a lab zone, or a fake resolver in automated tests).

## Production smoke test

1. Set `app.tracking_url` (or `LISTMONK_app__tracking_url`) to the platform tracking base; restart/reload if required.
2. Confirm Docker/app listens on **9000** and the edge proxies to it.
3. Add a customer domain in Admin → Tracking domains; create the CNAME to the platform tracking host.
4. Wait for propagation; run **Verify** until status is `verified` (not stuck on `pending` / `failed`).
5. Confirm TLS works for `https://<customer-host>/` with a cert that matches that hostname.
6. On the customer host: `GET /link/…` and `GET /campaign/…/px.png` succeed; `GET /admin` and `GET /api/…` are blocked.
7. Create a new campaign as a user who owns a verified domain; confirm rendered click/open URLs use `https://<customer-host>/…`.
8. Click a tracked link and load the open pixel; confirm analytics increment.
9. Confirm an older campaign without a domain still tracks via `app.tracking_url` (or `app.root_url`).

HTTP and HTTPS links in rendered HTML anchors are rewritten through the existing
listmonk link and click analytics tables. `mailto:`, `tel:`, fragment links,
subscription URLs, existing tracking URLs, images, and non-link attributes are
not rewritten.
