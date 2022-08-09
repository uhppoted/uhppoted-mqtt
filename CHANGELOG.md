# CHANGELOG

## [Unreleased]

### Added
1. HOWTO writeup for integrating with AWS Greengrass.

### Changed
1. Updated go.mod to Go 1.19


## [v0.8.1] - 2022-08-01

### Changed

1. Reworked event struct in `get-status`, `get-event`, `get-events` and `listen` to include:
   - event type code and description
   - event reason code and description
   - event direction code and description
2. Added (optional) protocol version to configuration.
3. Added (optional) translation locale to configuration.
4. Resolved INADDR_ANY to interface IPv4 address for controller listener address health check.


## [v0.8.0] - 2022-07-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.8.0

## [v0.7.3] - 2022-06-01

### Changed
1. Updated for compatibility with [uhppoted-lib](https://github.com/uhppoted/uhppoted-lib) v0.7.3

### v0.7.2

1. Migrated to uhppoted-lib `config` command implementation
2. Reworked `get-events` to return the `first`, `last` and `current` event indices.
3. Reworked `get-event`  to handle `first`, `last`, `current`, `next` and a valid index.

### v0.7.1

1. Task list support:
   -  `set-task-list`

2. Migrated to IUHPPOTED interface + implementation
