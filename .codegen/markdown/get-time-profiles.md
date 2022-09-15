{{- with .get_time_profiles -}}
### `{{.command}}`

Retrieves a range of time profiles from a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```

  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "from": 2,
      "to": 254
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-time-profile",
      "response": {
      "method": "get-time-profiles",
      "response": {
        "device-id": 405419896,
        "profiles": [
          {
            "end-date": "2021-12-31",
            "id": 29,
            "segments": [
              {
                "end": "17:00",
                "start": "08:30"
              },
              {
                "end": "00:00",
                "start": "00:00"
              },
              {
                "end": "00:00",
                "start": "00:00"
              }
            ],
            "start-date": "2021-01-01",
            "weekdays": "Monday"
          }
        ]
      }
    }
  }
}
```
{{end -}}
