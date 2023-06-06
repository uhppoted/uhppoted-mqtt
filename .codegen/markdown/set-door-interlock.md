{{- with .set_door_interlock -}}
### `{{.command}}`

Sets the door interlock mode for a controller.

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
      "interlock": 4
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-door-interlokc",
      "response": {
        "device-id": 405419896,
        "interlock": 4
      }
    }
  }
}
```
{{end -}}


