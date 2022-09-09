{{- with .set_door_delay -}}
### `{{.command}}`

Sets the open delay for a door on a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY54",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "door": 3,
      "delay": 8
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY54",
      "method": "set-door-delay",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 8
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "51f91c1aa04879c3e713e409e4ed67224a0256a059afeec4897dcc340e84488e"
}
```
{{end -}}


