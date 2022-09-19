{{- with .delete_cards -}}
### `{{.command}}`

Deletes all cards from a controller.

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
      "method": "delete-cards",
      "response": {
        "device-id": 405419896,
        "deleted": true
      }
    }
  }
}
```
{{end -}}
