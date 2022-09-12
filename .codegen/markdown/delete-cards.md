{{- with .delete_cards -}}
### `{{.command}}`

Deletes all cards from a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "delete-cards",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "deleted": true
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "63b6f0198f7d243c3e2c3ccf90866507fc8544818dbea757b4723952919c11a1"
}
```
{{end -}}
