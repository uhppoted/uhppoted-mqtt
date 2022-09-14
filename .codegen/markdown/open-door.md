{{- with .open_door -}}
### `{{.command}}`

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
{{template "request"  . -}}
{{template "response" . }}

Example:
```
{
  "message": {
    "request": {
{{- template "request-preamble"}}
      "device-id": 405419896,
      "door": 4,
      "card-number": 405419896
    }
  }
}
{
  "message": {
    "reply": {
{{- template "response-preamble"}}
      "method": "open-door",
      "response": {
        "device-id": 405419896,
        "door": 4,
        "opened": true
      }
    }
  }
}
```
{{end -}}


