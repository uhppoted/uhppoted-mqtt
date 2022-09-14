{{- with .get_door_delay -}}
### `{{.command}}`

Retrieves the open delay for a controller door.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
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
      "method": "get-door-delay",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "delay": 7
      }
    }
  }
}
```
{{end -}}


