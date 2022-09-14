{{- with .get_time_profile -}}
### `{{.command}}`

Retrieves a time profile from a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "profile-id": 29
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-time-profile",
      "response": {
        "device-id": 405419896,
        "time-profile": {
          "id": 29,
          "start-date": "2021-01-01",
          "end-date": "2021-12-31",
          "weekdays": "Monday"
          "segments": [
            {
              "start": "08:30",
              "end": "17:00"
            },
            {
              "start": "00:00",
              "end": "00:00"
            },
            {
              "start": "00:00",
              "end": "00:00"
            }
          ],
        }
      }
    }
  }
}
```
{{end -}}
