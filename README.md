# uhppote-mqtt

Wraps the `uhppote-core` device API in an MQTT endpoint for use with access control systems based on the 
*UHPPOTE UT0311-L0x* TCP/IP Wiegand access control boards.

Supported operating systems:
- Linux
- MacOS
- Windows

## Raison d'être

`uhppoted-mqtt` implements an MQTT endpoint that facilitates integration of the access control function with other 
systems (e.g. web servers, mobile applications) where the controller boards are located behind a firewall that does
not allow ingress. 

It also facilitates integration of access control with IoT systems based on [AWS IoT](https://aws.amazon.com/iot),
[Google Cloud IoT](https://cloud.google.com/solutions/iot) or the [IBM Watson IoT Platform](https://internetofthings.ibmcloud.com).

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.6.1    | Maintenance release to update module dependencies                                         |
| v0.6.0    | Maintenance release to update module dependencies                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

## Installation

### Building from source

#### Dependencies

| *Dependency*                          | *Description*                                          |
| ------------------------------------- | ------------------------------------------------------ |
| [com.github/uhppoted/uhppote-core][1] | Device level API implementation                        |
| [com.github/uhppoted/uhppoted-api][2] | common API for external applications                   |
| golang.org/x/sys/windows              | Support for Windows services                           |
| golang.org/x/lint/golint              | Additional *lint* check for release builds             |
| github.com/eclipse/paho.mqtt.golang   | Eclipse Paho MQTT client                               |
| github.com/gorilla/websocket          | paho.mqtt.golang dependency                            |

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

[1]: https://github.com/uhppoted/uhppote-core
[2]: https://github.com/uhppoted/uhppoted-api

