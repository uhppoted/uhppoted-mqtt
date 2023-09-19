{{- with .set_door_passcodes -}}
### `{{.command}}`

Sets up to four supervisor passcodes for a controller door, with valid passcodes being in the range [1..999999]. 
Invalid passcodes are replaces by 0 (no code).

{{template "request"  . -}}
{{template "response" . }}

Example:
```
topic: uhppoted/gateway/requests/device/door/passcodes:set

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "door": 3,
      "passcodes": [12345,999999,54321]
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "method": "set-door-control",
      "response": {
        "device-id": 405419896,
        "door": 3,
        "passcodes": [12345,999999,54321]
      }
    }
  }
}
```
{{end -}}

