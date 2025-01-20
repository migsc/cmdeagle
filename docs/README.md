
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

### Install with Go's package manager

```sh
go install github.com/migsc/cmdeagle
```

Prerequisites: 
- Go 1.23.2 or later. Installation intructions are available [here](https://go.dev/doc/install).
<!-- - Node.js (v16.17.0+) -->
<!-- - Node.js Package Manager (npm) -->

### Initialize a CLI starter project named yourcli

```sh
cmdeagle init yourcli
```

### Build the CLI from the root directory of your project:

```sh
cmdeagle build
```

### Create a sub command

...

### Define build script for your command

...

### Define start script for your command

...

### Build and run your new sub command

...

### Defining build scripts for executables

...

## Prerequisites
- Install cmdeagle

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
