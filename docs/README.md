
# cmdeagle

_**WARNING:** This is very much a work in progress, but we're close to releasing a stable version 1. Check back soon. Feedback is also much appreciated._

A language-agnostic CLI application build tool that allows you to create cross-platform applications written in any programming language of your choice.

## Features

- Define commands and subcommands from a single YAML configuration file.
- Reuse your existing external scripts and binaries to build commands.
- Define, parse, and validate your arguments and flags into environment variables your scripts can use.
- Create build steps and declare data/assets to install your commands.

## Quick Start

The easiest way to work with cmdeagle right now is by installing it with [Go](https://go.dev/doc/install) on a unix-like system (macOS, Linux, etc). More platforms will be supported soon.

### 1) Install with Go's package manager

```sh
go install github.com/migsc/cmdeagle@latest
```

Go 1.23.2 or later is required. Install it from [Golang's website](https://go.dev/doc/install).
<!-- - Node.js (v16.17.0+) -->
<!-- - Node.js Package Manager (npm) -->

### 2) Initializing a starter project

```sh
cmdeagle init <YOUR_CLI_NAME>
```

Where `<YOUR_CLI_NAME>` is the name of the executable file you want to build. By default, the binary will be named after the directory you run the command from.

You can change this later in the `.cmd.yaml` file for your new CLI project.

### 3) Building the CLI:

```sh
cmdeagle build -o .
```

Run this in the root directory of your project  to build a binary for your CLI to the current working directory. The binary will only run on your current operating system and architecture but you can target other platforms with flags. See the [reference]().

You can change the `-o` flag to build the executable file to a different directory, or run it without arguments, and it will build it to the default executable directory, which varies by operating system.

On macOS and on Linux, it will be built to either `./usr/local/bin` directory if you have the necessary write permissions, or `~/.local/bin` directory if you don't.

On Windows, it will be built to the `%LOCALAPPDATA%\Programs` directory.

The sample scripts are defined in several different languages to help you get started. For Python, JavaScript and other interpreted languages, the scripts are bundled together with the executable file thanks to the file paths defined in the `includes` setting.

```sh
includes:
- "./greet.sh"
- "./greet.js"
- "./greet.py"
```
This is useful for things like scripts, static assets, media, configuration files, data files, etc.

But if you're using Go, or some other compiled language, bundling executable binaries into your cmdeagle-built binary is not allowed for security reasons. So instead you must use the 'CLI_BIN_DIR' variable within a build script on the root or subcommand level `build` setting of your `.cmd.yaml`.

```sh
build: |
  go build -o $CLI_BIN_DIR/mycli-go-binary greet.go
  chmod +x $CLI_BIN_DIR/mycli-go-binary
```

Note that `$CLI_BIN_DIR` is a special environment variable that cmdeagle provides you within the shell that your build script will run in. And it points to the directory where the executable file should be built.

You can override this to point to a different directory by either defining an environment variable, or using the root or subcommand level 'env' setting in your '.cmd.yaml'. You can learn more about that here in the (#reference)[reference].

### 4) Running your CLI

You can run the executable file from your current working directory, or from anywhere on your system if you add its directory to the your system's PATH variable.

The `greet` command is a sample subcommand that you can use as a starting point for your own commands. It's defined in the `.cmd.yaml` file and it is configured to run the sample scripts in the project.

Let's test it out:

```sh
./<YOUR_CLI_NAME> greet cmdeagle 2 --uppercase --repeat 3
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
```

You can get more information about the `greet` command by running:

```sh
> ./<YOUR_CLI_NAME> help greet
```

Every one of those flags and arguments are defined in the `.cmd.yaml` file. Have a look and read the comments to learn how each configuration key works. You can also read the [reference](#reference) for more information on how to define your own commands, flags, and arguments.

For now, let's focus on the 'start' script defined for the 'greet' subcommand.

'''
# paste start script here.
'''

Note that it runs the 'my-cli-gobinary' mentioned before, if the '--use-go' flag is passed.

Let's invoke the built-in 'help' command now

```sh
> ./<YOUR_CLI_NAME> help
Usage:
  <YOUR_CLI_NAME> [flags]
  <YOUR_CLI_NAME> [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  greet       Greet the user.
  help        Help about any command

Flags:
  -h, --help   help for yourcli

Use "yourcli [command] --help" for more information about a command.
```

Where `<YOUR_CLI_NAME>` was the name you created in [step 2](#_2-initialize-a-cli-starter-project-named-yourcli).

The `completion` command will generate a script for your CLI to use in your shell. This is made possible because cmdeagle uses [cobra](https://github.com/spf13/cobra) under the hood. You can turn this off by setting the `completion` key to `false` at the root level of the `.cmd.yaml` file.

<!-- 

TODO: Move this to the reference section.

- <OPERATING_SYSTEM> defaults to your operating system. 

Where `<OPERATING_SYSTEM>` defaults to your operating  the operating system you want to build for, `<ARCHITECTURE>` is the architecture you want to build for, and `<OUTPUT_FILEPATH>` is the filepath where the executable file will be built.

You can also use the `-o` flag to specify the name and filepath of the executable file you want to build. By default, the executable file will be built in the current working directory. -->

## Reference

<!-- ## Using the cmdeagle CLI

### Initializing CLI configuration with the init Command

### Build your CLI with the build Command
- Targeting a Specific Platform -->

<!-- ## Installation -->


### Configuring Your CLI

Your CLI's schema is defined in a [YAML](https://en.wikipedia.org/wiki/YAML#cite_note-19) file named `.cmd.yaml`. There, you define your CLI's basic information and metadata, top-level command and sub commands, validation and parsing rules for your arguments and flags, build steps, bundled assets, and more.

<!-- ### Defining Your CLI's Basic Information and Metadata

...

### Top-level Command and Sub Commands

...

#### Defining Command's Basic Information 

...

#### Define Command-Level Environment Variables


##### `env` key

...

#### Defining Command-Level Validation

... -->

#### Command Lifecycle Scripts

<!-- ... -->

<!-- ##### The Command Scripts Lifecycle -->
- `requires` - Validates dependencies, once before building and once before running your command.
- `build` - Builds your command and create an exectuable file.
- `include` - Defines bundled assets to include in your CLI's executable file after building your command.
- `validate` - An optional script you can provide to validate your command's arguments and flags before running it.
- `start` - Executed when your users run your command.

<!-- ##### Defining `build` script

...

##### Defining `include` script

...

##### Defining `start` script 

...

##### Defining Bundled Assets To Include

... -->

#### Arguments and Flags

It's worth noting that cmdeagle assumes your arguments are positional, and that the order of your arguments in the configuration file determines their order in the command line. The flags you define do not have this behavior and can be defined in any order.

##### `name` key

Both arguments and flags must have a `name` key. The name key is important for the CLI to identify your argument or flag and it's also what you will use to reference their values within in your scripts.

See [Using Environment Variables](#using-environment-variables) for more information on how to reference their values.

In the future, we plan to relax the requirement for arguments to be named in order to better support arbitrary number of arguments.


##### `description` key

The description key defines a short description of the argument or flag. This description will be used in the help command to describe the argument or flag.

<!-- Long and short descriptions are not yet supported. -->

<!-- ###### `examples` key

... -->

##### `type` key

The type key defines how the CLI will parse the value of the argument or flag. Ultimately though, your scripts will still receive the raw value as a string due to the limitations of the shell.

- `string` (default value) - Effectively a no-op.
- `int` - Parses the value as an integer. Fails if the value cannot be parsed as an integer.
- `float` / `number` - Parses the value as a floating-point number. Fails if cannot be parsed as a number.
-  `bool` /`boolean` - Parses the value as a boolean. The resulting value will either be `"true"` or `"false"`. For flags, any value that is not `"false"` will be considered `"true"`, and the absence of the flag will be considered `"false"`. For arguments, the value must be `"true"` or `"false"`. Otherwise, the value will be considered invalid and the execution will fail.
<!-- - `date` -->
<!-- - `json[]`
- `json{}` -->

##### `pattern` key

Validates the value against a regular expression. It uses the [RE2 syntax](https://github.com/google/re2/wiki/Syntax) most commonly used by Perl, Python, and [Go](https://pkg.go.dev/regexp) programming languages.

```yaml
pattern: "^[a-zA-Z0-9]+$" #validates that the value is a string of alphanumeric characters
```
  <!-- # regex: "^[a-zA-Z0-9]+$" #validates that the value is a string of alphanumeric characters -->

<!-- - `uuid` - Validates the value against a UUID format.
- `email` - Validates the value against an email format.
- `url` - Validates the value against a URL format.
- `ip` - Validates the value against an IP address format.
- `ipv4` - Validates the value against an IPv4 address format.
- `ipv6` - Validates the value against an IPv6 address format. -->

##### `default` key

The default value key defines the value that will be used if the argument or flag is not provided.

```yaml
default: "World"
```

Note that the default value should match the `type` of the argument or flag. If it doesn't, the CLI will fail to build your command.

##### `required` key

Will fail if the argument or flag is not provided.

```yaml
required: true # defaults to false
```

##### `depends-on` key

Will fail if the argument or flag is not provided. In this example, the `last-name` flag will fail if the `first-name` flag is not provided.

```yaml
flags:
  - name: first-name
    type: string
    description: "Name to greet"

  - name: last-name
    type: string
    depends-on:
      - first-name
```
This results in the following environment behavior:

```sh
yourcli --first-name John 
# will succeed.

yourcli --first-name John --last-name Simpson 
# will succeed.

yourcli --last-name Simpson 
# will fail because the `first-name` flag is not provided. 

```


##### `conflicts-with` key

...

##### `validation` key (argument-and-flag-level)

You can define rules for a single argument or flag using the `validation` key within the argument or flag definition.


##### `eq` key
...

##### `gt` key
...

##### `gte` key
...

##### `lt` key
...

##### `lte` key
...

##### `min` key
...

##### `max` key
...

##### `range` key
...



##### `is-existing` key

- `file`
- `dir`
- `url`


##### `has-permissions` key

- `readable`
- `writable`
- `executable`

##### `check-stderr-on` key

This key can be used to define a custom validation by providing a script to be executed in the value of the key.

```yaml
check-stderr-on: |
  echo "Hello, World!"
```

If the script returns a non-zero exit code, or prints anything to the standard error stream, the validation will fail.

All environment variables normally available to your `build` and `start` scripts are also available to your custom validation script, including the arguments and flags passed to your command.


##### `has-stdout` key

...

<!-- ###### `alternative-for` key

...  -->

#### Validating Multiple Arguments and Flags

##### `validation` key

You can define validation rules for your arguments and flags as a whole with the `validation` key.


###### Argument-Only Validators:
...

##### `no-args` key

...

##### `arbitrary-args` key

...

##### `min-args` key

...

##### `max-args` key

...

##### `exact-args` key

...

##### `range-args` key

...


###### Conditional Validations with `and`, `or`, `not` keys

...

- Flag-Only          

##### Using Environment Variables
- Interpolated in Build Scripts
- Directly In Scripts and Source Code
- Directly In Build Scripts 
<!-- 
### Linux/Unix/macOS
- Link to Windows installation

### Node.js

### Typescript

 -->
