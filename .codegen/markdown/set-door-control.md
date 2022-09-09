{{- with .set_door_delay -}}
### `{{.command}}`

Sets the control mode for a controller door.

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
      "door": 3,
      "control": "normally closed"
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "set-door-control",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "control": "normally closed"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "5d2550d6d5a8c11c79fc531a4391fe80f6a66f241792b8432ee102ff0c20c833"
}
```
{{end -}}


