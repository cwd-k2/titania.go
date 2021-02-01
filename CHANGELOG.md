# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Change

- Now `tester.Detail.Error` includes `build_stderr`.

## [0.5.1] - 2021-02-02

### Removed

- Removed `jessevdk/go-flags` dependency.

## [0.5.0] - 2021-02-01

### Change

- `internal/pkg/viewer` is now a part of `pkg/tester`.
- Output while quiet testing is more informative.

### Added

- Added new option `--quiet`.

### Removed

- Experimental option `--async` was removed.

## [0.4.0] - 2021-01-31

### Change

- `paizaio` API change. Now you have to use `(*paizaio.Client).RunnersXXX` with `paizaio.Request` value, and get `paizaio.Response` value.

## [0.3.0] - 2021-01-30

### Added

- New command `pio`. Executes a single file using [paiza.io](https://paiza.io) API.

### Change

- Many of refactorings.
- Now `pretty` library is internal.
- Cut out the `tester/viewer` as `viewer` internal library.
- Now `tester` and `paizaio` libraries are public package.

### Fixed

- Fix `paizaio.RunnersGetDetailsResponse.BuildTime` was not handled properly.

## [0.2.2] - 2020-11-07

### Change

- Some refactorings.
- The output json's `test_matter` fields are now `name` fields.

## [0.2.1] - 2020-10-25

### Change

- Introduce [project layout](https://github.com/golang-standards/project-layout).
- Introduce this CHANGELOG.md.
- Now `pretty` library in this project is public.
