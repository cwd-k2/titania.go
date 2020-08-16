# [WIP] titania.go

Just a testing tool, for single-file-programs using paiza.io API.

## Installation

```bash
$ go get -u github.com/cwd-k2/titania.go
```

## Run

```bash
$ titania.go [options] [directories]
```

### Prerequisite

Make sure `titania.go` is in your `$PATH`.

And the target directories should contain `titania.json`, that is like below.

```
{
  "host": "http://api.paiza.io:80", // paiza.io API host
  "api_key": "guest",               // paiza.io API api_key
  "source_code_directories": ["./source_*"], // directories containing source codes
  "test_case_directories": ["./test_case"],  // directories containing input/answer for test
  "test_case_input_extension": ".in",        // input files' extension
  "test_case_output_extension": ".out",      // answer files' extension
  "max_processes": 10               // I haven't implemented yet
}
```

### Example

Examples are in [here](https://github.com/cwd-k2/titania.example).

```bash
$ cd example

# just specify one directory.
$ titania.go example_01

# two or more is also okay.
$ titania.go example_01 example_02

# or you can run with no args, then titania.go will take all subdirectories as targets.
$ titania.go

# you can specify languages to test on, with -lang flags
$ titania.go --lang=ruby,haskell
```

## **IMPORTANT**

All _INTERESTING_ (but not so important) information that pop up at runtime, will write out to **STDERR**.

And the **test results' details** will be written to **STDOUT**, in JSON format.

So, if you want to record the results, below is the usage.

```bash
$ titania.go [directories] > <filename>
```

Or, if you just want to know runtime information (includes summary), just go like below.

```bash
$ titania.go [directories] > /dev/null
```
