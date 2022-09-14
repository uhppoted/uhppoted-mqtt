{{- with .set_door_delay -}}
### `{{.command}}`

Sets the open delay for a door on a controller.

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
      "delay": 8
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-door-delay",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 8
      }
    }
  }
}
```
{{end -}}


