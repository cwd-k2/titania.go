# [WIP] titania.go

Just a testing tool, for single-file-programs using paiza.io API.

## Installation

It's simple. Just run `go get`

```bash
$ go get -u github.com/cwd-k2/titania.go
```

## Usage

Make sure `titania.go` is in your `$PATH`.

```bash
$ titania.go [options] [directories]
```

The options should be before other arguments.

### Prerequisite

The target directories should contain `titania.json`, like below.

```
{
  "host": "http://api.paiza.io:80", // paiza.io API host
  "api_key": "guest",               // paiza.io API api_key
  "source_code_directories": ["./source_*"], // directories containing source codes (relative paths)
  "test_case_directories": ["./test_case"],  // directories containing input/answer for test (relative paths)
  "test_case_input_extension": ".in",        // input files' extension
  "test_case_answer_extension": ".out"       // answer files' extension
}
```

### Example

Examples are in [here](https://github.com/cwd-k2/titania.example).

```bash
$ cd example

# if no arguments specified, titania.go will take all subdirectories as targets.
$ titania.go

# target directories' paths as arguments
$ titania.go example_01 example_02

# you can specify languages to test on, with -lang flag
$ titania.go --lang=ruby,haskell
```

## Philosophy

All _interesting_ (but not so important) information that pop up at runtime, will write out to **STDERR**.

And the **test results' details** will be written to **STDOUT**, in JSON format.

So, if you want to record the results, below is the usage.

```bash
$ titania.go [targets] > filename.json
```

Or, if you just want to know runtime (fancy) information, below is the way you go.

```bash
$ titania.go [targets] > /dev/null
```
