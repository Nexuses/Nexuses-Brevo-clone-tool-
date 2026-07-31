# Subscriber import samples

Use these files with **Subscribers → Import** in the admin UI (or the [import API](https://listmonk.app/docs/apis/import.html)).

## `subscriber-import-template.csv`

Example bulk-import file. Required and optional columns:

| Column       | Required | Description |
|-------------|----------|-------------|
| `email`     | Yes      | Subscriber e-mail address (plain address only). |
| `firstname` | No       | First name. Stored for templates as `{{ .Subscriber.FirstName }}`. |
| `lastname`  | No       | Last name. Stored for templates as `{{ .Subscriber.LastName }}`. |
| `company`   | No       | Company name. Stored for templates as `{{ .Subscriber.Company }}`. |
| `name`      | No       | Legacy display name (single field). Used when `firstname` / `lastname` are empty. |
| `attributes`| No       | JSON object as a **single CSV cell**. Use doubled quotes (`""`) inside the cell for JSON string quotes. |

Rules:

- First row must be the header with these exact names (column order can vary).
- If `firstname` and/or `lastname` are set, the display name is built from them unless `name` is also set (then `name` wins).
- If no name fields are set, listmonk derives a display name from the e-mail local part.
- Default delimiter is comma (`,`); set another delimiter on the import screen if needed.
- Upload `.csv` or a `.zip` containing one CSV file.
- Extra columns are ignored.

## HTML / campaign placeholders

Map CSV data to templates with [Go template](https://listmonk.app/docs/templating.html) syntax:

| CSV / DB field | Template expression |
|----------------|---------------------|
| `email` | `{{ .Subscriber.Email }}` |
| `firstname` | `{{ .Subscriber.FirstName }}` |
| `lastname` | `{{ .Subscriber.LastName }}` |
| `company` | `{{ .Subscriber.Company }}` |
| `name` (legacy) | `{{ .Subscriber.Name }}` |
| `attributes` → `city` | `{{ .Subscriber.Attribs.city }}` |
| `attributes` → `tier` | `{{ .Subscriber.Attribs.tier }}` |

Example campaign HTML:

```html
<p>Hi {{ .Subscriber.FirstName }},</p>
<p>Thanks for being part of {{ .Subscriber.Company }}.</p>
<p>You are on the <strong>{{ .Subscriber.Attribs.tier }}</strong> plan in {{ .Subscriber.Attribs.city }}.</p>
<p>Questions? Reply to this e-mail or contact us at support@example.com.</p>
```

Attribute keys in the CSV JSON must match the names used after `.Subscriber.Attribs.` (case-sensitive).

## Spreadsheet editors

Open `subscriber-import-template.csv` in Excel, Google Sheets, or LibreOffice. When saving:

- Keep UTF-8 encoding.
- For cells with JSON in `attributes`, wrap the whole JSON in quotes and escape inner quotes as `""` (see rows in the template).
- Do not add extra columns unless you merge custom fields into the `attributes` JSON object.
