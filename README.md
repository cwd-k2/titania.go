# [WIP] titania.go

Just a code-testing/judging tool, using [paiza.io](https://paiza.io) API. Intended to replace other local judgment tools.

**This program is still in a WIP state, so big changes may be introduced. Be careful.**

## What is this for?

### Expected situation

This program is intended to be used for the situations like below.
- Programming contests || Judging your codes locally
- Coding tests
- Just testing your algorithm/logic is sutable for any possible situations
- When you came up with an excellent programming/algorithm problem and wanna have a small contest with your friends, and they use different languages, say, A uses `C++`, B uses `Haskell` while C uses `Scala`..., but you don't have all those runtime environment

The prerequisites will be like below.
- A topic/problem is provided.
- A few or many inputs are provided.
- There is some expected output/answer for each input.
- A submitted code will be a single file.

### What it is not for

`titania.go` is **NOT** for the developing situations: when you're working with a great and fancy project and you want to test your method or class implementation is right or not.

If you cut out the logic or implementation into a single file, this tool can be applied somehow, but uh not so suitable I guess.

## Installation

It's simple. Just run `go get`.

```bash
$ go get -u github.com/cwd-k2/titania.go/cmd/titania.go
```

## Usage

Make sure `titania.go` is in your `$PATH`.

```bash
$ titania.go [options] [directories]
```

Here, one directory should stand for one problem or topic.

A directory contains one problem, some amount of test cases, and different solutions in different languages.

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

`titania.go` detects source codes' languages by their extensions. You can see the languages available [here](https://github.com/cwd-k2/titania.go/blob/master/internal/tester/util.go).

### Input and Answer

For test cases to run, `titania.go` requires both input and expected output files.

Correspoiding input and answer files should have same **names**, except their extensions.

### Additional testing method

Sometimes you want to do a specialized judgement. When the problem goes like '...print one of the possible answers', just comparing the output and the expected one is not enough.

`titania.go` just compares the STDOUT of the submitted code and expected output/answer by default (simply check the equality), but you can set another testing method, by setting `test_method`. It means, you can also write an additional code to check if the answer was right, or partially-right, or wrong-but-the-quality-is-worth-to-be-praised.

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

For test method execution, the code's output, the test case's input and expected answer are handed through STDIN, and they are joined by null characters.

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

Here you can see some parts of this result output are awfully named (`test_matter`??? what???). Yes, it is my fault, and if you have a good idea then please tell me, I will totally appreciate it.

```json
[
  {
    "test_matter": "example_01",
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
            "is_expected": true,
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
            "is_expected": false,
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
            "is_expected": true,
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

This is because we use STDOUT for pipe, redirect, or some other possible ways.

If you want to record the results, below is the usage.

```bash
$ titania.go [directories] > filename.json
```

Or, if you just want to know runtime information, below is the way you go (Like the examples above).

```bash
$ titania.go [directories] > /dev/null
```
