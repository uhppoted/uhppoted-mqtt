{{- with .set_door_delay -}}
### `{{.command}}`

Sets the control mode for a controller door.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "door": 3,
      "control": "normally closed"
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-door-control",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "control": "normally closed"
      }
    }
  }
}
```
{{end -}}


