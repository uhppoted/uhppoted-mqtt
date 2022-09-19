{{- with .acl_grant -}}
### `{{.command}}`

Grants system-wide access control permissions for a card. The permissions are
added to any existing permissions.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/{{ .request.topic }}

{
  "message": {
    "request": {
{{- template "request-preamble"}}
    "request": {
      "card-number": 8165538,
      "start-date": "2022-01-01",
      "end-date": "2022-12-31",
      "doors": [
        "Gryffindor",
        "Slytherin"
      ]
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "acl:grant",
      "response": {
        "granted": true
      }
    }
  }
}
```
{{end -}}
