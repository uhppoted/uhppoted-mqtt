{{- with .get_device -}}
### `{{.command}}`

Returns the controller information for a UHPPOTE controller.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
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
      "method": "get-device",
      "response": {
        "address": {
          "IP": "192.168.1.100",
          "Port": 60000,
          "Zone": ""
        },
        "date": "2018-11-05",
        "device-id": 405419896,
        "device-type": "UTO311-L04",
        "gateway-address": "192.168.1.1",
        "ip-address": "192.168.1.100",
        "mac-address": "00:12:23:34:45:56",
        "subnet-mask": "255.255.255.0",
        "timezone": {},
        "version": "0892"
      }
    }
  }
}
```
{{end -}}


