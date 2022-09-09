{{- with .get_door_delay -}}
### `{{.command}}`

Retrieves the control mode for a controller door.

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
      "method": "get-door-control",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "control": "controlled"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "d4eab0760fc831bd16760db16055d72b46129c8540b432bd5967bb5f8ba6b665"
}
```
{{end -}}


