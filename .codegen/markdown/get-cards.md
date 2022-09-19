{{- with .get_cards -}}
### `{{.command}}`

Retrieves a list of the cards stored on a controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
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
      "method": "get-cards",
      "response": {
        "device-id": 405419896,
        "cards": [
          8165537,
          8165539,
          8165538
        ]
      }
    }
  }
}
```
{{end -}}


