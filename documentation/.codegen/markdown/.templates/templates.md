{{define "request"}}
```
Request:
{
    "message": {
        "request": {
{{- range .request.fields}}
            "{{.field}}": "{{.value}}",
{{- end}}
        }
    }
}
{{range .request.fields}}
{{printf "%-12s" .field}} {{.description}}
{{- end}}
```
{{end}}


{{define "response"}}
```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "{{.command}}",
      "response": {
{{- range .response.fields}}
            "{{.field}}": "{{.value}}",
{{- end}}
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
{{- range .response.fields}}
{{printf "%-12s" .field}} {{.description}}
{{- end}}
```
{{end}}