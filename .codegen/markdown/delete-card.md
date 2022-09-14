{{- with .delete_card -}}
### `{{.command}}`

Deletes a card from a controller.

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
      "method": "delete-card",
      "response": {
        "device-id": 405419896,
        "card-number": 8165538,
        "deleted": true
      }
    }
  }
}
```
{{end -}}
