{{- with .get_door_delay -}}
### `{{.command}}`

Retrieves the open delay for a controller door.

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
      "door": 3
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-door-delay",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 7
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "506fe61c581fbd3cd438aa89d984b6f382bdf95a8f0170a1b3a2edd073fe4c03"
}
```
{{end -}}


