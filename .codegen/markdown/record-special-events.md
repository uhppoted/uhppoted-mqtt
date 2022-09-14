{{- with .record_special_events -}}
### `{{.command}}`

Enables/disables event logging for door open, close and pushbutton events.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "enabled": true
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "record-special-events",
      "response": {
        "DeviceID": 405419896,
        "Enable": true,
        "Updated": true
      }
    }
  }
}
```
{{end -}}


