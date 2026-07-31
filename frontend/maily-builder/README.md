# Maily builder (listmonk)

Embeds [Maily](https://maily.to/) (`@maily-to/core` + `@maily-to/render`) as a UMD bundle for the Vue admin UI.

## Status

**Scaffold / Phase 1** — builds a standalone editor bundle. Not yet wired into campaign content types (see parent `frontend/` integration plan).

## Build

```bash
cd frontend/maily-builder && yarn install && yarn build
```

Output is copied to `frontend/public/static/maily-builder/` by the root Makefile.

## Integration contract

Same pattern as `frontend/email-builder`:

- `body_source` — Maily JSON (`JSONContent`)
- `body` — HTML from `@maily-to/render` with Go template placeholders (`{{ .Subscriber.Email }}`, etc.)

Variables are **not** resolved at edit time; listmonk's Go template engine resolves them at send time.
