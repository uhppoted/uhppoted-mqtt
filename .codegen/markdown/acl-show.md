{{- with .acl_show -}}
### `{{.command}}`

Retrieves the system-wide access control permissions for a card.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "card-number": 8165538
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "acl:show",
      "response": {
        "card-number": 8165538,
        "permissions": [
          {
            "door": "Gryffindor",
            "end-date": "2021-12-31",
            "start-date": "2021-01-01"
          },
          {
            "door": "Slytherin",
            "end-date": "2021-12-31",
            "start-date": "2021-01-01"
          }
        ]
      }
    }
  }
}
```
{{end -}}

