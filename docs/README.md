
# cmdeagle

_**WARNING:** This is very much a work in progress, but we're close to releasing a stable version 1. Check back soon. Feedback is also much appreciated._

A versatile build tool that allows you to build cross-platform CLI applications written in any programming language of your choice.

### Features:

- Define commands and subcommands from a single YAML configuration file.
- Reuse your existing external scripts and binaries to build commands.
- Define, parse, and validate your arguments and flags into environment variables your scripts can use.
- Create build steps, declare assets and data to bundle for each command, bundled assets, and more.

# Quick Start

The easiest way to install cmdeagle right now is with [Go](https://go.dev/doc/install), on a unix-like system (macOS, Linux, etc). More platforms will be supported soon.

### 1) Install with Go's package manager

```sh
go install github.com/migsc/cmdeagle
```

Go 1.23.2 or later is required. Install it from [Golang's website](https://go.dev/doc/install).
<!-- - Node.js (v16.17.0+) -->
<!-- - Node.js Package Manager (npm) -->

### 2) Initialize a CLI starter project named yourcli

```sh
cmdeagle init <YOUR_CLI_NAME>
```

Where `<YOUR_CLI_NAME>` is the name of the executable file you want to build. By default, the binary will be named after the directory you run the command from.

You can change this later in the `.cmd.yaml` file ofr your new CLI project.

### 3) Build the CLI:

You can build a binary for your CLI that will only run on your current operating system and architecture from the root directory of your project:

```sh
cmdeagle build -o .
```

Which will build the executable file to the current working directory.

Change the `-o` flag to build the executable file to a different directory, or run it without arguments, and it will build it to the default executable directory, which varies by operating system.

On macOS and on Linux, it will be built to either `./usr/local/bin` directory if you have write permissions, or `~/.local/bin` directory if you don't.

On Windows, it will be built to the `%LOCALAPPDATA%\Programs` directory.

### 4) Run your CLI

You can run the executable file from your current working directory, or from anywhere on your system if you add its directory to the `PATH` environment variable. Let's invoke its help command to see what it can do.

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

Every one of those flags and arguments are defined in the `.cmd.yaml` file. Have a look and read the comments to learn how each configuration key works.

The sample scripts are defined in several different languages to help you get started. For Python, JavaScript and other interpreted languages, the scripts are bundled together with the executable file via the `includes` key. 

```sh
includes:
# You can bundle files with your CLI by declaring the paths to them in the `includes` field.
# This is useful for things like static assets, media, configuration files, data files, etc.
- "./greet.sh"
- "./greet.js"
- "./greet.py"
```

But for Go, and other compiled languages, the executable needs to be built and installed to the system's default executable directory seperately. This is done with the `build` key which defines a build script for either your CLI or a specific subcommand in the `.cmd.yaml` file.

```sh
build: |
  go build -o $CLI_BIN_DIR/mycli-go-binary greet.go
  chmod +x $CLI_BIN_DIR/mycli-go-binary
```

Note that `$CLI_BIN_DIR` is a special environment variable that cmdeagle provides you in the build script that points to the directory where the executable file should be built.

The start script You can also read the [reference](#reference) section for more information on how to define your own commands, flags, and arguments.

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

The `completion` command will generate a script for your CLI to use in your shell. This is made possible because cmdeagle uses [cobra](https://github.com/spf13/cobra) under the hood. You can turn this off by setting the `completion` key to `false` in the `.cmd.yaml` file.



<!-- 

TODO: Move this to the reference section.

- <OPERATING_SYSTEM> defaults to your operating system. 

Where `<OPERATING_SYSTEM>` defaults to your operating  the operating system you want to build for, `<ARCHITECTURE>` is the architecture you want to build for, and `<OUTPUT_FILEPATH>` is the filepath where the executable file will be built.

You can also use the `-o` flag to specify the name and filepath of the executable file you want to build. By default, the executable file will be built in the current working directory. -->

# Reference

## Using the cmdeagle CLI

### Initializing CLI configuration with the init Command

### Build your CLI with the build Command
- Targeting a Specific Platform

## Configuring Your CLI

Your CLI's schema is defined in a [YAML](https://en.wikipedia.org/wiki/YAML#cite_note-19) file named `.cmd.yaml`. There, you define your CLI's basic information and metadata, top-level command and sub commands, validation and parsing rules for your arguments and flags, build steps, bundled assets, and more.

### Defining Your CLI's Basic Information and Metadata

...

### Top-level Command and Sub Commands

...

#### Defining Command's Basic Information 

...

#### Define Command-Level Environment Variables


##### `env` key

...

#### Defining Command-Level Validation

...

#### Command Lifecycle Scripts

...

##### The Command Scripts Lifecycle
- `requires` - Validates dependencies, once before building and once before running your command.
- `build` - Builds your command and create an exectuable file
- `include` - Defines bundled assets to include in your CLI's executable file after building your command
- `start` - Executed when you or your users run your command

##### Defining `build` script

...

##### Defining `include` script

...

##### Defining `start` script 

...

##### Defining Bundled Assets To Include

...

#### Arguments and Flags

Arguments and flags are defined with `args` and `flags` key.

<!-- ##### Definition -->

Their definitions are relatively similar.

```yaml
args:
  - name: name
    type: string
    description: "Name to greet"
    default: "World"
    required: true
flags:
  - name: uppercase
    shorthand: u
    type: boolean
    description: "Convert greeting to uppercase"
    default: false
```

###### `name`

The name key is important for the CLI to identify your argument or flag. It's also 
what you will use to reference their values within in your scripts.

See [Using Environment Variables](#using-environment-variables) for more information on how to reference their values.

In the future, we plan to relax the requirement for arguments to be named in order
to better support arbitrary number of arguments.

It's also worth noting that cmdeagle assumes that your arguments are positional,
and that the order of your arguments in the configuration file determines their
order in the command line. 

###### `description` key

The description key defines a short description of the argument or flag. This description will be used in the help command to describe the argument or flag.

<!-- Long and short descriptions are not yet supported. -->

<!-- ###### `examples` key

... -->

###### `type` key

The type key defines how the CLI will parse the value of the argument or flag. Ultimately though, your scripts will still receive the raw value as a string due to the limitations of the shell.

- `string` (default value) - Effectively a no-op.
- `int` - Parses the value as an integer. Fails if the value cannot be parsed as an integer.
- `float` / `number` - Parses the value as a floating-point number. Fails if cannot be parsed as a number.
-  `bool` /`boolean` - Parses the value as a boolean. The resulting value will either be `"true"` or `"false"`. For flags, any value that is not `"false"` will be considered `"true"`, and the absence of the flag will be considered `"false"`. For arguments, the value must be `"true"` or `"false"`. Otherwise, the value will be considered invalid and the execution will fail.
<!-- - `date` -->
<!-- - `json[]`
- `json{}` -->

###### `pattern` key

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

###### `default` key

The default value key defines the value that will be used if the argument or flag is not provided.

```yaml
default: "World"
```

Note that the default value should match the `type` of the argument or flag. If it doesn't, the CLI will fail to build your command.

###### `required` key

Will fail if the argument or flag is not provided.

```yaml
required: true # defaults to false
```

###### `depends-on` key

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


###### `conflicts-with` key

...

###### `validation` key (argument-and-flag-level)

You can define rules for a single argument or flag using the `validation` key within the argument or flag definition.


###### `eq` key
...

###### `gt` key
...

###### `gte` key
...

###### `lt` key
...

###### `lte` key
...

###### `min` key
...

###### `max` key
...

###### `range` key
...



###### `is-existing` key

- `file`
- `dir`
- `url`


###### `has-permissions` key

- `readable`
- `writable`
- `executable`



<!-- ###### `alternative-for` key

...  -->

### Validation and Parsing

#### `validation` key

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

### Linux/Unix/macOS
- Link to Windows installation

### Node.js

### Typescript

## Installation

## Install CLI 
