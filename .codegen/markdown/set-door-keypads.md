{{- with .set_door_keypads -}}
### `{{.command}}`

Activates/deactivates the reader access keypads for a controller.

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
      "keypads": {
        "1": true,
        "2": true,
        "3": false,
        "4": true
      }
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-door-keypads",
      "response": {
        "device-id": 405419896,
        "keypads": {
          "1": true,
          "2": true,
          "3": false,
          "4": true
        }
      }
    }
  }
}
```
{{end -}}


