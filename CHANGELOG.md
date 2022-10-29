# CHANGELOG

## [Unreleased]

### Changed
1. Removed the 'false' option for acl.verify (ref. https://github.com/uhppoted/uhppoted/discussions/17)

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
