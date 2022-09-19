{{- with .get_door_delay -}}
### `{{.command}}`

Retrieves the control mode for a controller door.

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
      "door": 3
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-door-control",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "control": "controlled"
      }
    }
  }
}
```
{{end -}}


