{{- with .get_time -}}
### `{{.command}}`

Returns the controller date and time.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-time",
      "response": {
        "date-time": "2022-09-08 11:01:04 PDT",
        "device-id": 405419896
      }
    }
  }
}
```
{{end -}}


