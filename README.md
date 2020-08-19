# [WIP] titania.go

Just a testing tool, using paiza.io API.

## Installation

It's simple. Just run `go get`.

```bash
$ go get -u github.com/cwd-k2/titania.go
```

## Usage

Make sure `titania.go` is in your `$PATH`.

```bash
$ titania.go [options] [directories]
```

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

Examples are [here](https://github.com/cwd-k2/titania.example).

```bash
$ cd example

# if no arguments specified, titania.go will take all subdirectories as targets.
$ titania.go > /dev/null

# target directories' paths as arguments
$ titania.go example_01 example_02 > /dev/null

# you can specify languages to test on, with -lang flag
$ titania.go --lang=ruby,haskell > /dev/null
```

### Output

**Test results' details** will be written to **STDOUT**, in JSON format. Like below.

```json
[
  {
    "target": "example_02",
    "fruits": [
      {
        "source_code": "source_haskell/main.hs",
        "language": "haskell",
        "details": [
          {
            "test_case": "test_case/01",
            "result": "PASS",
            "time": "0.00",
            "output": "50\n",
            "error": ""
          },
          {
            "test_case": "test_case/02",
            "result": "PASS",
            "time": "0.00",
            "output": "20\n",
            "error": ""
          }
        ]
      },
      {
        "source_code": "source_ruby/main.rb",
        "language": "ruby",
        "details": [
          {
            "test_case": "test_case/01",
            "result": "PASS",
            "time": "0.09",
            "output": "50\n",
            "error": ""
          },
          {
            "test_case": "test_case/02",
            "result": "PASS",
            "time": "0.10",
            "output": "20\n",
            "error": ""
          }
        ]
      }
    ]
  }
]
```


## Philosophy

All fancy (but not so important) information that pop up at runtime, will write out to **STDERR**.

This is because we use STDOUT for pipe, redirect, or so.

If you want to record the results, below is the usage.

```bash
$ titania.go [targets] > filename.json
```

Or, if you just want to know runtime information, below is the way you go (Like example).

```bash
$ titania.go [targets] > /dev/null
```
