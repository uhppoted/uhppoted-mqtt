{{- with .delete_card -}}
### `{{.command}}`

Deletes a card from a controller.

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
      "device-id": 405419896,
      "card-number": 8165538
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "delete-card",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "card-number": 8165538,
        "deleted": true
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "4c7ca3ad707943c6c6853082d42906985e0aa669eaddecfcabe9f9a3b60fc34d"
}
```
{{end -}}
