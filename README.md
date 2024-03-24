![build](https://github.com/uhppoted/uhppoted-mqtt/workflows/build/badge.svg)
![build](https://github.com/uhppoted/uhppoted-mqtt/workflows/ghcr/badge.svg)

# uhppoted-mqtt

Wraps the `uhppote-core` device API in an MQTT endpoint for use with access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards.

Supported operating systems:
- Linux
- MacOS
- Windows
- ARM7 _(e.g. RaspberryPi)_

## Raison d'être

`uhppoted-mqtt` implements an MQTT endpoint that facilitates integration of the access control function with other systems (e.g. web servers, mobile applications) where the controllers are located behind a firewall that does not allow ingress. 

It also facilitates integration of access control with IoT systems based on e.g. [AWS IoT](https://aws.amazon.com/iot),
[Google Cloud IoT](https://cloud.google.com/solutions/iot) or the [IBM Watson IoT Platform](https://internetofthings.ibmcloud.com).

---
### Contents

- [Release Notes](#release-notes)
- [Installation](#installation)
   - [Docker](#docker)
   - [Building from source](#building-from-source)
- [Command line](#uhppoted-mqtt)
- [Notes](#notes)
---


## Release Notes

### Current Release

**[v0.8.7](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.7) - 2023-12-01**

1. `set-door-passcodes` command to set supervisor passcodes for a door.
2. Published received events to the _live_ events topic.
3. Renamed _master_ branch to _main_ in line with current development practice.
4. Replaced `nil` event pointer in `get-status` with zero value.
5. Fixed _double_ events in `events::Listen` (cf. https://github.com/uhppoted/uhppoted-mqtt/issues/15)
6. Reworked events _feed_ to poll for unretrieved events.


## Installation

Executables for all the supported operating systems are packaged in the [releases](https://github.com/uhppoted/uhppoted-rest/releases):

The release tarballs contain the executables for all the operating systems - OS specific tarballs with all the _uhppoted_ components can be found in [uhpppoted](https://github.com/uhppoted/uhppoted/releases) releases.

Installation is straightforward - download the archive and extract it to a directory of your choice. To install `uhppoted-mqttd` as a system service:
```
   cd <uhppote directory>
   sudo uhppoted-mqttd daemonize
```

`uhppoted-mqttd help` will list the available commands and associated options (documented below).

The `uhppoted-mqttd` service requires the following additional files:

- `uhppoted.conf`

### `uhppoted.conf`

`uhppoted.conf` is the communal configuration file shared by all the `uhppoted` project modules and is (or will 
eventually be) documented in [uhppoted](https://github.com/uhppoted/uhppoted). `uhppoted-mqttd` requires:
- the _MQTT_ section to define the configuration for the MQTT client connection and endpoint
- the _devices_ section to resolve non-local controller IP addresses and door to controller door identities.

A sample [uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/main/runtime/simulation/405419896.conf) file
is included in the `uhppoted` distribution. Alternatively, a starter `uhppoted.conf` file can be created by executing:
```
uhppoted-mqtt config > uhppoted.conf
```

### Docker

A public _Docker_ image is published to [ghcr.io](https://github.com/uhppoted?tab=packages&repo_name=uhppoted-mqtt). 

The image is configured to use the `/usr/local/etc/uhppoted/uhppoted.conf` file for configuration information.

#### `docker compose`

A sample Docker `compose` configuration is provided in the [`scripts/docker/compose`](scripts/docker/compose) folder. 

To run the example, download and extract the [zipped](script/docker/compose.zip) scripts and supporting files into folder
of your choice and then:
```
cd <compose folder>
docker compose up
```

The MQTT server can be tested using an MQTT client, e.g. _mqtt_:
```
mqtt publish --topic 'uhppoted/gateway/requests/devices:get' \
              --message '{ "message": { "request": { "request-id": "AH173635G3",     \
                                                     "client-id":  "QWERTY54",       \
                                                     "reply-to":   "uhppoted/reply", \
                                                      "nonce":     5 }}}'

```

The default image is configured for TCP only. To enable a TLS connection to the MQTT broker, enable it in the _uhppoted.conf_ file
on the Docker volume mapped to _/usr/local/etc/uhppoted_ and copy the certificates and keys to the _/usr/local/etc/uhppoted/mqtt_ folder on 
the Docker volume mapped to _/usr/local/etc/uhppoted_, e.g.
```
docker cp broker.pem  mqttd:/usr/local/etc/uhppoted/mqtt
docker cp client.key  mqttd:/usr/local/etc/uhppoted/mqtt
docker cp client.cert mqttd:/usr/local/etc/uhppoted/mqtt
```

#### `docker run`

To start a REST server using Docker `run`:
```
docker pull ghcr.io/uhppoted/restd:latest
docker run --publish 8080:8080 --publish 8443:8443 --name restd --mount source=uhppoted,target=/var/uhppoted --rm ghcr.io/uhppoted/restd
```

The REST server can be tested using _curl_, e.g.:
```
curl -X 'GET' 'http://127.0.0.1:8080/uhppote/device' -H 'accept: application/json' | jq .
```

#### `docker build`

For inclusion in a Dockerfile:
```
FROM ghcr.io/uhppoted/restd:latest
```


### Building from source

Assuming you have `Go` and `make` installed:

```
git clone https://github.com/uhppoted/uhppoted-mqtt.git
cd uhppoted-mqtt
make build
```

If you prefer not to use `make`:
```
git clone https://github.com/uhppoted/uhppoted-mqtt.git
cd uhppoted-mqtt
mkdir bin
go build -trimpath -o bin ./...
```

The above commands build the `'uhppoted-mqtt` executable to the `bin` directory.


### Building from source

#### Dependencies

| *Dependency*                                             | *Description*                                          |
| -------------------------------------------------------- | ------------------------------------------------------ |
| [uhppote-core](https://github.com/uhppoted/uhppote-core) | Device level API implementation                        |
| [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) | common API for external applications                   |
| github.com/eclipse/paho.mqtt.golang                      | Eclipse Paho MQTT client                               |
| golang.org/x/sys                                         | Support for Windows services                           |
| golang.org/x/net                                         | paho.mqtt.golang dependency                            |
| github.com/gorilla/websocket                             | paho.mqtt.golang dependency                            |

## uhppoted-mqtt

Usage: ```uhppoted-mqtt <command> <options>```

Supported commands:

- `help`
- `version`
- `config`
- `run`
- `daemonize`
- `undaemonize`

Defaults to `run` if the command it not provided i.e. ```uhppoted-mqtt <options>``` is equivalent to ```uhppoted-mqtt run <options>```.

### `run`

Runs the `uhppoted-mqtt` MQTT client. Intended for use as a system service that runs in the background to handle MQTT messages directed at the endpoint. 

Command line:

` uhppoted-mqtt [--debug] [--console] [--config <file>] `

```
  --config      Sets the uhppoted.conf file to use for controller configurations. 
                Defaults to the communal uhppoted.conf file shared by all the uhppoted 
                modules.
  --console     Runs the MQTT endpoint as an application, logging events to the
                console.
  --debug       Displays verbose debugging information, in particular the communications with the UHPPOTE controllers
```

### `daemonize`

Registers the `uhppoted-mqtt` MQTT endpoint as a system service that will be started on system boot. The command creates the necessary system specific service configuration files and service manager entries.

Command line:

`uhppoted-mqtt daemonize `

### `undaemonize`

Unregisters the `uhppoted-mqtt` MQTT endpoint as a system service, but does not delete any created log or configuration files. 

Command line:

`uhppoted-mqtt undaemonize `

## Notes

### Events

_uhppoted-mqtt_ publishes controller events to two topics:

- `uhppoted/gateway/events`
- `uhppoted/gateway/events/live`

The `uhppoted/gateway/events` topic is an events _feed_ that attempts to ensure that all events are published to the topic at least 
once, including events that may have occurred while the service was offline. Events are retrieved in small batches so if there is a
significant backlog the event published in the feed may lag significantly behind the current event.

The `uhppoted/gateway/events/live` topic is an events _feed_ that publishes controller events as they are received. However, events
that occur while the service is offline won't be published to this topic.

The topics are configurable in _uhppoted.conf_:
```
…
mqtt.topic.root = uhppoted/gateway
mqtt.topic.events = ./events
mqtt.topic.events.live = ./events/live
…
```




