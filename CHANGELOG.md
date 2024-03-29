# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.10.0] - 2022-06-26

### Added

- Option `--no-json`: do not print output json.

### Change

- Perf: Now titania.go doesn't hold input/output/code/... data.
- Hence, intermediate files should be always created now, and the default directory is `./titania-out`.

## [0.9.3] - 2021-06-24

### Added

- Added `"language"` for `"input_order"`.

## [0.9.2] - 2021-06-22

### Change

- Update Go to 1.16

## [0.9.1] - 2021-06-22

### Change

- Fix null pointer reference crash & command line parsing issues.

## [0.9.0] - 2021-06-21

### Added

- New fields `"input_prefix"` and `"answer_prefix"`. You can now use those prefixes to identify each test case.

### Change

- `"input_extension"` and `"output_extension"` are deprecated. They are now `"input_suffix"` and `"answer_suffix"` respectively.

## [0.8.2] - 2021-06-18

### Change

- Output Json's format changed. Now `"fruits"` & `"details"` are named as `"test_targets"` & `"test_cases"`, `"is_expected"` is replaced by `"expect"`, etc.

## [0.8.1] - 2021-06-17

### Change

- Doesn't create temporary files by default.

## [0.8.0] - 2021-06-17

### Added

- You can now set the maximum jobs to run concurrently by `--maxjob=N`.

## [0.7.1] - 2021-06-17

### Change

- `piorun` command now takes `os.Stdin` as default input.
- Now you can set `"delimiter"` and more detailed order for `"test_method"`.

## [0.7.0] - 2021-06-03

### Added

- Added `simplejson` subpackage.
- Added a functionality to store intermediate output files to a temporary directory.
- Added `--temdir=` option to change the temporary directory.

### Change

- `paizaio` package was renamed to `runner`.
- `pio` command was renamed to `piorun`.
- Now output json is not pretty-printed by default. You can use `--pretty` flag to pretty-print.
- `internal/pkg/pretty` package are now public.
- `internal/pkg/langtype` were merged into `runner` package.
- You can now set `config -> "test_target" -> "expected"` an object that explains what results is expected for each test case. Still, a plain string is OK.
- Now `config -> "test_method" -> "on"` is `config -> "test_method" -> "on_exit"`, and you can set an integer on which the method will be executed.

## [0.6.2] - 2021-02-17

### Added

- You can now specify input order for TestMethod execution.

## [0.6.1] - 2021-02-05

### Added

- Now TypeScript and なでしこ are available (Following the change of [paiza.io](https://paiza.io) API).

## [0.6.0] - 2021-02-02

### Change

- Now `tester.Detail.Error` includes `build_stderr`.
- Now `tester.NewTestUnit` takes `*tester.Config` as argument. You have to create a `*tester.Config` instance before `tester.NewTestUnit`

### Added

- Now `paizaio.RunnersGetDetailsResponse` has `Error` field.
- Better error handling.

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
