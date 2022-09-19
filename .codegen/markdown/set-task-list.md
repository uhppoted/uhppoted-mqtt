{{- with .set_task_list -}}
### `{{.command}}`

Stores a tasklist to a controller.

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
      "tasks": [
          {
            "task": "trigger once",
            "door": 3,
            "start-date": "2021-01-01",
            "end-date": "2021-12-31",
            "weekdays": "Monday,Wednesday,Friday",
            "start": "08:27"
          }
        ]
      }
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "set-time-profile",
      "response": {
        "device-id": 405419896,
        "warnings": []
      }
    }
  }
}
```
{{end -}}
