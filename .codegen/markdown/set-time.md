{{- with .get_time -}}
### `{{.command}}`

Sets the controller date and time.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "date-time": "2022-09-09 11:25:04"
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-time",
      "response": {
        "device-id": 405419896,
        "date-time": "2022-09-09 11:25:04 PDT"
      }
    }
  }
}
```
{{end -}}


