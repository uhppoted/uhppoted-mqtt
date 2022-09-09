{{- with .record_special_events -}}
### `{{.command}}`

Enables/disables event logging for door open, close and pushbutton events.

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
      "enabled": true
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "record-special-events",
      "request-id": "AH173635G3",
      "response": {
        "DeviceID": 405419896,
        "Enable": true,
        "Updated": true
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "c7767d44fd50c7cda2ee5e1bf1af2d86b8a2367a4a330fda0a84427cc8c1c8ef"
}
```
{{end -}}


