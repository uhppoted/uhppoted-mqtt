### `get-devices`

Returns a list of all UHPPOTE controllers found via a UDP broadcast on the local LAN or specifically
configured in _uhppoted.conf_.


```
Request:

topic: uhppoted/gateway/requests/devices:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-devices",
      "response": {
            "device-id": "<controller-id>",
            "device-type": "<controller-type>",
            "ip-address": "<address>",
            "port": "<UDPv4 port>",
      },
      ...
    }
  },
  ...
}

request-id   message ID from the request
client-id    client ID from the request
device-id    controller serial number
device-type  controller type (UTO311-L0x)
ip-address   controller IPv4 address
port         UDP port for controller commands
```


Example:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "nonce": 5
    }
  }
}

{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "get-devices",
      "request-id": "AH173635G3",
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
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "77b8917b80e18837fc9c13e451194d17f2712fb95b8f5283f9c762dfe1ed4f55"
}
```
