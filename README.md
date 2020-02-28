# uhppote-mqtt

MQTT endpoint for access control systems based on the *UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards. 

Supported operating systems:
- Linux
- MacOS
- Windows

## Raison d'être

The manufacturer supplied application is 'Windows-only' and provides limited support for integration with other
systems.

## Releases

## Installation

### Building from source

#### Dependencies

| *Dependency*                        | *Description*                                          |
| ----------------------------------- | ------------------------------------------------------ |
| com.github/uhppoted/uhppote-core    | Device level API implementation                        |
| com.github/uhppoted/uhppoted-api    | External API implementation                            |
| golang.org/x/sys/windows            | Support for Windows services                           |
| golang.org/x/lint/golint            | Additional *lint* check for release builds             |
| github.com/eclipse/paho.mqtt.golang | Eclipse Paho MQTT client                               |
| github.com/gorilla/websocket        | paho.mqtt.golang dependency                            |

### Binaries

## uhppoted-mqtt

Usage: *uhppoted-mqtt \<command\> \<options\>*

Defaults to 'run' unless one of the commands below is specified: 

- daemonize
- undaemonize
- config
- help
- version

Supported 'run' options:
- --console
- --debug

