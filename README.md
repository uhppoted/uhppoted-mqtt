![build](https://github.com/uhppoted/uhppoted-mqtt/workflows/build/badge.svg)

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

## Releases

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
| v0.8.2    | Added softlock for client connection and first draft of command documentation             |
| v0.8.1    | Added human readable event fields and fixed health check for INADDR_ANY addresses         |
| v0.8.0    | Maintenance release for version compatibility with `uhppoted-core` `v0.8.0`               |
| v0.7.3    | Maintenance release for version compatibility with `uhppoted-core` `v0.7.3`               |
| v0.7.2    | Reworked event handling (including removal of rollover)                                   |
| v0.7.1    | Added `set-task-list` command to manage the controller task list                          |
| v0.7.0    | Added support time profiles from the extended API                                         |
| v0.6.12   | Added support for `nil` events in response to `get-status`                                |
| v0.6.10   | Maintenance release for version compatibility with `uhppoted-app-wild-apricot`            |
| v0.6.8    | Maintenance release for version compatibility with `uhppote-core` `v0.6.8`                |
| v0.6.7    | Implements `special-events` message to enable/disable door events                         |
| v0.6.5    | Maintenance release for version compatibility with `node-red-contrib-uhppoted`            |
| v0.6.4    | Maintenance release for version compatibility with `uhppoted-app-sheets`                  |
| v0.6.3    | Implements ACL commands                                                                   |
| v0.6.2    | Maintenance release to update module dependencies                                         |
| v0.6.1    | Maintenance release to update module dependencies                                         |
| v0.6.0    | Maintenance release to update module dependencies                                         |
| v0.5.1    | Initial release following restructuring into standalone Go *modules* and *git submodules* |

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

A sample [uhppoted.conf](https://github.com/uhppoted/uhppoted/blob/master/runtime/simulation/405419896.conf) file
is included in the `uhppoted` distribution. Alternatively, a starter `uhppoted.conf` file can be created by executing:
```
uhppoted-mqtt config > uhppoted.conf
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




