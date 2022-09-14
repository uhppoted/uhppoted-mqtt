{{- with .get_devices -}}
### `{{.command}}`

Returns a list of all UHPPOTE controllers found via a UDP broadcast on the local LAN or specifically
configured in _uhppoted.conf_.

{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
    }
  }
}

{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "get-devices",
      "response": {
        "devices": {
          "201020304": {
            "device-type": "UTO311-L02",
            "ip-address": "192.168.1.101",
            "port": 60000
          },
          "303986753": {
            "device-type": "UTO311-L03",
            "ip-address": "192.168.1.100",
            "port": 60000
          },
          "405419896": {
            "device-type": "UTO311-L04",
            "ip-address": "192.168.1.100",
            "port": 60000
          }
        }
      }
    }
  }
}
```
{{end -}}


