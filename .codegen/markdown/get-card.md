{{- with .get_card -}}
### `{{.command}}`

Retrieves a card record from a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "card-number": 8165538
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-card",
      "response": {
        "device-id": 405419896,
        "card": {
          "card-number": 8165538,
          "start-date": "2021-01-01"
          "end-date": "2021-12-31",
          "doors": {
            "1": 1,
            "2": 0,
            "3": 0,
            "4": 1
          }
        }
      }
    }
  }
}
```
{{end -}}
