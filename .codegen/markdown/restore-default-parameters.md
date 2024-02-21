{{- with .restore_default_parameters -}}
### `{{.command}}`

Resets a controller to the manufacturer default configuration.

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
      "method": "restore-default-parameters",
      "response": {
        "device-id": 405419896,
        "reset": true
      }
    }
  }
}
```
{{end -}}
