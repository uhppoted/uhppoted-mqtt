# CHANGELOG

## Unreleased

### Added
1. ARMv6 build target for Pi ZeroW


## [0.8.9](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.9) - 2024-09-06

### Added
1. TCP/IP support.

### Updated
1. Updated to Go 1.23.


## [0.8.8](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.8) - 2024-03-27

### Added
1. `restore-default-parameters` API function to reset a controller to the manufacturer default configuration.
2. Added public Docker image to ghcr.io.

### Updated
1. Bumped Go version to 1.22.


## [0.8.7](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.7) - 2023-12-01

### Added
1. `set-door-passcodes` command to set supervisor passcodes for a door.
2. Published received events to the _live_ events topic.

### Updated
1. Renamed _master_ branch to _main_ in line with current development practice.
2. Replaced `nil` event pointer in `get-status` with zero value.
3. Fixed _double_ events in `events::Listen` (cf. https://github.com/uhppoted/uhppoted-mqtt/issues/15)
4. Reworked events _feed_ to poll for unretrieved events.


## [0.8.6](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.6) - 2023-08-30

### Added
1. `set-door-keypads` command to activate/deactivate reader access keypads.
2. Added Dockerfile and associated Makefile targets.
3. Preliminary documentation for _uhppoted.conf_ configuration file.


## [0.8.5](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.5) - 2023-06-13

### Added
1. `set-interlock` API function to set controller door interlock mode.

### Updated
1. Reworked to use cards with pointerless 'from' and 'to' dates 
2. Added _staticcheck_ linting to build


## [0.8.4](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.4) - 2023-03-17

### Added
1. `doc.go` package overview documentation.
3. Added PIN support to get-card and put-card APIs.

### Updated
1. Replaced (most) logging with uhppoted-lib logging.
2. Fixed Windows event logging.
3. Added Cloud Discovery section to AWS Greengrass HOWTO.


## [0.8.3](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.3) - 2022-12-16

### Added
1. Added ARM64 to release build artifacts
2. Basic working version of AWS Greengrass integration HOWTO

### Changed
1. Removed the 'false' option for acl.verify (ref. https://github.com/uhppoted/uhppoted/discussions/17)
2. Replaced service lockfile implementation with _flock_ equivalent
3. Replaced MQTT client 'soft-lock' with _flock_ equivalent
4. Updated _systemd_ unit file to wait on `network-online.target`
5. Removed _zip_ files from release artifacts (no longer necessary)

### Removed
1. Removed ARM7 specific daemonize (was only required for _softlock_).


## [0.8.2](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.2) - 2022-10-14

### Added
1. HOWTO writeup for integrating with AWS Greengrass (in progress).
2. ARM7 specific `daemonize` implementation to delete the MQTT client lock file on service start.
3. User guide for MQTT requests/responses and security
4. Added `mqttd.acl.verify` to configuration

### Changed
1. Updated go.mod to Go 1.19
2. Reworked MQTT client connection lockfile to implement 'soft lock'
3. Added logic to shutdown application on too many MQTT client disconnects within the monitoring interval
4. Added mime-type to `acl:download` and `acl:compare` request messages. Zip files now expect a mime-type of
   `application/zip`
5. `acl:download` and `acl:compare` can fetch plain TSV files if the mqtt.acl.verify configuration allows
   unsigned downloads (i.e. `mqtt.acl.verify = none` or `mqtt.acl.verify = not-empty,RSA`)
6. Fixed response message capitalization for `record-special-events`
7. Reworked RecordSpecialEvents to not use wrapped requests/responses
8. Added 'swipe open' and 'swipe close' event reasons to message internationalisation.


## [0.8.1](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.1) - 2022-08-01

### Changed

1. Reworked event struct in `get-status`, `get-event`, `get-events` and `listen` to include:
   - event type code and description
   - event reason code and description
   - event direction code and description
2. Added (optional) protocol version to configuration.
3. Added (optional) translation locale to configuration.
4. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.


## [0.8.0](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.8.0) - 2022-07-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.0

## [0.7.3](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.7.3) - 2022-06-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.7.3

### [0.7.2](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.7.2)

1. Migrated to uhppoted-lib `config` command implementation
2. Reworked `get-events` to return the `first`, `last` and `current` event indices.
3. Reworked `get-event`  to handle `first`, `last`, `current`, `next` and a valid index.

### [0.7.1](https://github.com/uhppoted/uhppoted-mqtt/releases/tag/v0.7.1)

1. Task list support:
   -  `set-task-list`

2. Migrated to IUHPPOTED interface + implementation


## Older

| *Version* | *Description*                                                                             |
| --------- | ----------------------------------------------------------------------------------------- |
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

