### `get-device`

Returns the controller information for a UHPPOTE controller.


```
Request:

topic: <root>/<requests>/device:get

message:
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": "<controller-id>",
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the
                        configured reply topic) if not provided.
device-id    (required) controller serial number
```

```
Response:
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "get-device",
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
topic: uhppoted/gateway/requests/device:get

topic: uhppoted/gateway/requests/device:get

{
  "message": {
    "request": {
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896
    }
  }
}

{
  "message": {
    "reply": {
      "server-id": "uhppoted"
      "client-id": "QWERTY",
      "request-id": "AH173635G3",
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
