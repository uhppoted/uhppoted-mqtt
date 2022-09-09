{{- with .get_time -}}
### `{{.command}}`

Sets the controller date and time.

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
      "date-time": "2022-09-09 11:25:04"
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "set-time",
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "date-time": "2022-09-09 11:25:04 PDT"
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "1e3a06eb2022315bc5fb6335d4ff3ef62c75c4a6c5b395fd43e5f1876703bb27"
}
```
{{end -}}


