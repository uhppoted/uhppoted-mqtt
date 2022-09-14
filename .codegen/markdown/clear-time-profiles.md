{{- with .clear_time_profiles -}}
### `{{.command}}`

Clears all time profiles on a a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
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
      "method": "clear-time-profiles",
      "response": {
        "cleared": true,
        "device-id": 405419896
      },
    }
  },
}
```
{{end -}}
