{{- with .acl_revoke -}}
### `{{.command}}`

Revokes system-wide access control permissions for a card. The permissions are
removed from any existing permissions.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
    "request": {
      "card-number": 8165538,
      "start-date": "2022-01-01",
      "end-date": "2022-12-31",
      "doors": [
        "Dungeon"
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
        "revoked": true
      }
    }
  }
}
```
{{end}}