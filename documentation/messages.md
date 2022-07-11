# Command and Control Functions

This document describes the command and control functions implemented by _uhppoted_mqtt_.

## Security

### Authentication

### Authorisation

### Encryption

## Commands

1.  `get-devices`
2.  `get-device`
3.  `get-time`
4.  `set-time`
5.  `get-door-delay`
6.  `set-door-delay`
7.  `get-door-control`
8.  `set-door-control`
9.  `record-special-events`
10. `get-status`
11. `get-cards`
12. `delete-cards`
13. `get-card`
14. `put-card`
15. `delete-card`
16. `get-time-profile`
17. `set-time-profile`
18. `clear-time-profiles`
19. `get-time-profiles`
20. `set-time-profiles`
21. `set-task-list`
22. `get-events`
23. `get-event`
24. [`open-door`](#open_door)
25. `acl-show`
26. `acl-grant`
27. `acl-revoke`
28. `acl-upload-file`
29. `acl-upload-s3`
30. `acl-upload-http`
31. `acl-download-file`
32. `acl-download-s3`
33. `acl-download-http`
34. `acl-compare-file`
35. `acl-compare-s3`
36. `acl-compare-http`

### `open-door`

Remotely opens a door after verifying that the card in the _request_ has permission:
- the card number must match a line in the _uhppoted-mqtt_ `cards` file 
- the card number must be a valid card on the controller with permission to open the
  door

The `cards` file is the first level of authorisation and is intended to restrict the use
of the `open-door` command to a limited set of supervisor or override cards. It defaults 
to _mqtt/cards_ file in the _uhppoted_ _etc_ directory:
- `/usr/local/etc/com.github.uhppoted/mqtt/cards (MacOS)`
- `/etc/uhppoted/mqtt/cards (Linux)`
- `\Program Data\uhppoted\mqttcards (Windows)`
- `.\mqtt\cards (Windows)`

but can be configured with the _mqtt.cards_ value in the _uhppoted.conf_ file:
```
# MQTT
...
mqtt.cards = /usr/local/etc/com.github.uhppoted/mqtt/cards
...
```

Each line in the file should be a regular expression that matches one or more authorised cards. A catch-all `.*` regular expression will authorise all cards e.g.:
```
.*
```

In addition, the card number in the request must be valid for the controller:
- the start date must be on or before _today_
- the end date must be on or after _today_
- the card must have access to the door

e.g.
```
uhppote-cli get-card 405419896 8165538                                

405419896  8165538  2022-01-01 2022-12-31 Y N N Y
```
Request:
```
{
    "message": {
        "request": {
            "request-id": "<request-id>",
            "client-id": "<client-id>",
            "reply-to": "<topic>",
            "device-id": <controller-id>,
            "door": <door-id>,
            "card-number": <card-number>
        }
    }
}

request-id   (optional) message ID, returned in the response
client-id    (required) client ID for authentication and authorisation (if enabled)
reply-to     (optional) topic for reply message. Defaults to uhppoted/gateway/replies (or the configured reply topic) if not provided.
device-id    (required) controller serial number
door         (required) door (1..4) to open
card-number  (required) card number used to validate access
```

Response:
```
{
  "message": {
    "reply": {
      "request-id": <request-id>,
      "client-id": <client-id>,
      "method": "open-door",
      "response": {
        "device-id": 405419896,
        "door": 4,
        "opened": true
      },
      ...
    }
  },
  ...
}

request-id  message ID from the request
client-id   client ID from the request
device-id   controller serial number
door        door (1..4) from the request
opened      true/false
```

Example:
```
{
  "message": {
    "request": {
      "request-id": "AH173635G3",
      "client-id": "QWERTY",
      "reply-to": "uhppoted/reply/97531",
      "device-id": 405419896,
      "door": 4,
      "card-number": 405419896
    }
  }
}
{
  "message": {
    "reply": {
      "client-id": "QWERTY",
      "method": "open-door",
      "nonce": 185,
      "request-id": "AH173635G3",
      "response": {
        "device-id": 405419896,
        "door": 4,
        "opened": true
      },
      "server-id": "uhppoted"
    }
  },
  "hmac": "fd480076bc713b9f331c0fe33feb470fd8877002448cdcd712d83792d4cb4919"
}
```



