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

## How to Use

Make sure `titania.go` is in your `$PATH`.

```bash
$ titania.go [options] [directories]
```

Here, one directory should stand for one problem or topic.

A directory contains one problem, some amount of test cases, and different solutions in different languages.

### Prerequisite

The directories should contain `titania.json`, like below.

```json
{
  "client": {
    "host": "https://api.paiza.io",
    "api_key": "guest"
  },
  "test_target": [
    { "pattern": "./source_code/main.*" },
    {
      "pattern": "./source_code/wrong_answer.rb",
      "expect": { "default": "FAIL", "test_case/01": "PASS", "test_case/03": "EXECUTION TIMEOUT" }
    },
    { "pattern": "./source_code/time_out.rb" },
    { "pattern": "./source_code/build_error.c", "expect": "BUILD FAILURE" }
  ],
  "test_method": {
    "file_name": "./special_judge.cpp",
    "on_exit": 0,
    "input_order": ["input", "source_code", "build_error", "stdout"]
  },
  "test_case": [
    {
      "directory": "./test_case",
      "input_extension": ".in",
      "output_extension": ".ans"
    }
  ]
}
```

### Fields in `titania.json`

#### `"client"`

The paiza.io api client config. Fields are `"host"` and `"api_key"`. See [here](https://paiza.io/en)

```
{
  "host":    string
  "api_key": string
}
```

#### `"test_target"`

An array of objects.

```
{
  "pattern": string
  "expect":  string | map[string]string  (Optional)
}
```

The `"pattern"` field is shell glob pattern, required to collect each single file to be tested.

The `"expect"` field is the expected test result for the specified target. You can set more detailed one using object. Without this field, the expected result is treated as `{ "default": "PASS" }` or `"PASS"`.

```
{
  "pattern": "./source_code/main.rb",
  "expect": {
    "default": "PASS",
    "test_case/01": "EXECUTION TIMEOUT",
    "test_case/02": "FAIL"
  }
}
```

#### `"test_case"`

An array of objects.

```
{
  "directory":        string
  "input_extension":  string
  "output_extension": string
}
```

The `"directory"` field indicates the directory that contains test-cases, which are composed of some input-output pairs of files with same basename.

The `"input_extension"` and `"output_extension"` field are simple. With which the program recognize which is the input and which is the expected output.

**Note**: A input-output pair should have input, but not about output.

#### `"test_method"`

An object.

```
{
  "file_name":   string
  "delimiter":   string
  "on_exit":     int
  "input_order": []string
}
```

Without `"test_method"`, results are judged by simple comparison of `"test_target"`s' output and expected output.

With this, more complicated judge will be available, if you have any idea.

The `"file_name"` and `"on_exit"` field are simple. Just specify a path to the test method file and specify on which exit code of test target the judge will be conducted. The `"delimiter"` is a special string, which can be used in `"input_order"` section.

Okay, let's see the default value for it.

```
{
  ...,
  "input_order": ["stdout", "delimiter" , "input", "delimiter", "answer"]
}
```

Then the input for this test method goes like this.

`<STDOUT><delimiter><INPUT><delimiter><ANSWER>`

Available elements for `"input_order"` is...

- `"input"`
  - The input for `"test_target"`
- `"answer"`
  - The expected output or answer you've prepared for `"test_case"`
- `"stdout"`
  - The output on STDOUT.
- `"stderr"`
  - The output on STDOUT.
- `"build_stdout"`
  - The output on STDOUT, when building/compiling.
- `"build_stderr"`
  - The output on STDOUT, when building/compiling.
- `"source_code"`
  - Then source code of `"test_target"`
- `"delimiter"`
  - You can use this special string in input_order. Default is `"\x00"` (null character).
- `"newline"`
  - A newline (`"\n"`)
- `"tab"`
  - A tab character (`"\t"`)

And the test method will be like below.

**Note**: the first line of the output of the test method will be treated as the result.

```ruby
# output, input, answer can be separated by null charactor.
# ruby's `gets` will read until the argument.
# so, gets nil means read all from stdin.
output, input, answer = gets(nil).split("\0")

# and test method should output PASS or FAIL or other RESULT on the first line.
STDOUT.puts output == answer ? "FAIL" : "PASS"
```

### Example

Examples are [here](https://github.com/cwd-k2/titania.go/blob/master/examples).

```bash
$ cd example

# if no arguments specified, titania.go will take all subdirectories.
$ titania.go > /dev/null

# tested directories' paths as arguments
$ titania.go example_01 example_02 > /dev/null

# you can specify languages to test on, with -lang flag
$ titania.go --lang=ruby,haskell > /dev/null
```

### Options

#### `--lang=LANG1[,LANG2[,..]]`

You can narrow down tests by languages.

#### `--quiet`

Suppres runtime information output.

#### `--pretty`

Pretty printing the output json.

#### `--tmpdir=DIRNAME`

You can set where to put the intermediate files when executing tests. If not specified, intermediate files will not be created.

#### `--maxjob=N`

You can set the maximum number of concurrent jobs. As default, there is no limit.

### Source codes; Languages

`titania.go` detects source codes' languages by their extensions. You can see the languages available [here](https://github.com/cwd-k2/titania.go/blob/master/pkg/runner/langtype.go).

## Results

**Test results' details** will be written to **STDOUT**, in JSON format. Like below.

```
[
  {
    "name": "example_01",
    "test_method": "default",
    "test_targets": [
      {
        "name": "source_code/main.hs",
        "language": "haskell",
        "test_cases": [
          {
            "name": "test_case/01",
            "time": "0.00",
            "expect": "PASS",
            "result": "PASS",
            "output": "50\n",
            "others": "[1 of 1] Compiling Main             ( Main.hs, Main.o )\nLinking Main ...\n",
            "errors": ""
          }, ...
        ], ...
      }, ...
    ], ...
  }, ...
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
