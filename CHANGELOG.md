# CHANGELOG

## [Unreleased]

### Changed

1. Reworked event struct in `get-status`, `get-event` and `get-events` response to include:
   - event type code and description
   - event reason code and description
   - event direction code and description

2. Added protocol version to configuration.

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