{{- with .put_card -}}
### `{{.command}}`

Adds or updates a card record on a controller.

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
      "card": {
        "card-number": 8165538,
        "start-date": "2021-01-01",
        "end-date": "2021-12-31",
        "doors": {
          "1": true,
          "2": false,
          "3": 55,
          "4": false
        }
      }
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "put-card",
      "response": {
        "device-id": 405419896
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": true,
            "2": false,
            "3": 55,
            "4": false
          }
        }
      }
    }
  }
}
```
{{end -}}
