# [WIP] titania.go

Just a testing tool, for single-file-programs using paiza.io API.

## Installation

```sh
$ go get -u github.com/cwd-k2/titania.go
```

## Run

```sh
$ titania.go [directories]
```

### Prerequisite

Make sure `$GOBIN` is included in your `$PATH`.

And the target directories should contain `titania.json`.

```
{
  "host": "http://api.paiza.io:80", // paiza.io host
  "api_key": "guest",               // paiza.io api_key
  "source_code_directories": ["./source_*"], // directories containing source codes
  "test_case_directories": ["./test_case"],  // directories containing input/answer for test
  "test_case_input_extension": ".in",        // input files' extension
  "test_case_output_extension": ".out",      // answer files' extension
  "max_processes": 10               // meh I haven't implemented yet
}
```

`titania.json`'s example is [here](./example/example_01/titania.json) or [there](./example/example_02/titania.json)

### Example

Examples are in [here](./example).

```sh
$ cd example
# just specify one directory.
$ titania.go example_01
# two or more is also okay.
$ titania.go example_01 example_02
# or you can run with no args, then titania.go will take all subdirectories as targets.
$ titania.go
```
