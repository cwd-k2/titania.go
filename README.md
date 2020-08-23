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

The directories should contain `titania.json`, like below.

```
{
  "client": {
    "host": "http://api.paiza.io:80",
    "api_key": "guest"
  },
  "test_target": [
    { "pattern": "./source_code/main.*" },
    { "pattern": "./source_code/wrong_answer.rb", "expect": "FAIL" },
    { "pattern": "./source_code/time_out.rb" },
    { "pattern": "./source_code/build_error.c", "expect": "BUILD FAILURE" }
  ],
  "test_case": [
    {
      "directory": "./test_case",
      "input_extension": ".in",
      "output_extension": ".ans"
    }
  ]
}
```

### Example

Examples are [here](https://github.com/cwd-k2/titania.example).

```bash
$ cd example

# if no arguments specified, titania.go will take all subdirectories.
$ titania.go > /dev/null

# tested directories' paths as arguments
$ titania.go example_01 example_02 > /dev/null

# you can specify languages to test on, with -lang flag
$ titania.go --lang=ruby,haskell > /dev/null
```

### Source codes; Languages

`titania.go` detects source codes' languages by thier extensions.

### Input and Answer

For test cases to run, `titania.go` requires both input and expected output files.

Correspoiding input and answer files should have same **names**, except thier extensions.

### Additional testing method

`titania.go` compare execution STDOUT and answer by default, but you can set another testing method, by setting `test_method`.

```json
{
  "client": {
    "host": "http://api.paiza.io:80",
    "api_key": "guest"
  },
  "test_target": [
    { "pattern": "./source_*/*.*" }
  ],
  "test_case": [
    { "directory": "./test_case", "input_extension": ".in", "output_extension": ".out" }
  ],
  "test_method": {
    "file_name": "./test_method.rb"
  }
}
```

For test method execution, a single case's execution output, test case's input and expected answer are handed through STDIN, and they are joined by null character.

Simply, goes like this.

`<STDOUT>\0<STDIN>\0<ANSWER>\0`

Okay, let's see an example.

```ruby
# output, input, answer can be separated by null charactor.
# ruby's `gets` will read until the argument.
# note that all below have \0 on their end.
output = gets "\0"
input  = gets "\0"
answer = gets "\0"

# and test method should output PASS or FAIL
STDOUT.puts output == answer ? "FAIL" : "PASS"
STDERR.puts "Just for fun, inversing the result."
```

This '\0 separator' strategy would be a kind of awful, but I didn't come up with any other solutions.


## Results

**Test results' details** will be written to **STDOUT**, in JSON format. Like below.

```json
[
  {
    "test_topic": "example_01",
    "test_method": "default",
    "fruits": [
      {
        "test_target": "source_code/main.c",
        "language": "c",
        "expect": "PASS",
        "details": [
          {
            "test_case": "test_case/01",
            "result": "PASS",
            "expected": true,
            "time": "0.00",
            "output": "50\n",
            "error": ""
          },
          (snip)
        ]
      },
      (snip)
      {
        "test_target": "source_code/time_out.rb",
        "language": "ruby",
        "expect": "PASS",
        "details": [
          {
            "test_case": "test_case/01",
            "result": "EXECUTION TIMEOUT",
            "expected": false,
            "time": "",
            "output": "",
            "error": ""
          },
          (snip)
        ]
      },
      {
        "test_target": "source_code/build_error.c",
        "language": "c",
        "expect": "BUILD FAILURE",
        "details": [
          {
            "test_case": "test_case/01",
            "result": "BUILD FAILURE",
            "expected": true,
            "time": "",
            "output": "",
            "error": "Main.c:1:1: error: unknown type name 'include'\ninclude <stdio.h>\n^\nMain.c:1:9: error: expected identifier or '('\ninclude <stdio.h>\n        ^\n2 errors generated.\n"
          },
          (snip)
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
$ titania.go [directories] > filename.json
```

Or, if you just want to know runtime information, below is the way you go (Like example).

```bash
$ titania.go [directories] > /dev/null
```
