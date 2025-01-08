# cmdeagle

**WARNING: This is very much a work in progress, but we're close to releasing a stable version 1. Check back soon. Feedback is also much appreciated.**

A versatile framework that allows users to build cross-platform command-line tools using the language of their choice. It's designed to simplify the process of creating complex CLI applications by referencing external scripts and  binaries, using a YAML configuration file to define commands, arguments, flags, build commands and more.

Internally, it heavily leverages existing CLI libraries and tools such as:

  - Cobra [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
  - Viper [https://github.com/spf13/viper](https://github.com/spf13/viper)
  - pflag [https://github.com/spf13/pflag](https://github.com/spf13/pflag)
  - huh [https://github.com/charmbracelet/huh](https://github.com/charmbracelet/huh)
  - log [https://github.com/charmbracelet/log](https://github.com/charmbracelet/log)

Big thanks to [@spf13](https://github.com/spf13) and [@charmbracelet](https://github.com/charmbracelet) for their work on these libraries. Building cmdeagle was made possible with their contributions to the opensource community. 

## Prerequisites

- Go 1.22 or higher: You need to have Go installed to install cmdeagle and build CLIs. You can install Go by following the instructions at https://go.dev/doc/install.
- Node.js 14.0 or higher


## Installation

```
go install github.com/migsc/cmdeagle
```

### Install from source

```
git clone https://github.com/migsc/cmdeagle.git
cd cmdeagle
make build
```

Then add the following line to your `.zshrc` or `.bashrc` to make it easier to invoke `cmdeagle` directly. Then run `source ~/.zshrc` or `source ~/.bashrc` to apply the changes.

```
PATH="$PATH:PATH_TO_YOUR_CLONE_DIR/cmdeagle/bin"
```

Where `PATH_TO_YOUR_CLONE_DIR` is the path to the directory you cloned cmdeagle into.

## Usage

1. Initialize a new CLI project by running in a directory of your choice:

```
mdkir ./mycli && cd mycli
cmdeagle init
```

This will create a `..cmd.yaml` file in your project root to define your CLI structure. The schema is not documented yet, but you can read the comments in the file to understand what each field does.

2. Build your CLI:

```
cmdeagle build
```

3. Run your CLI.

For Linux and macOS, if you have `./.local/bin` in your `PATH` you can simply run:

```
mycli
```

Otherwise you can run it like so:

```
cd
./.local/bin/mycli
```

If you need help setting up your `PATH` variable, add the following line in your `.zshrc` or `.bashrc`. This will make  it easier to run your CLI by invoking it directly. Then run `source ~/.zshrc` or `source ~/.bashrc` to apply the  changes.

```
PATH="$PATH:~/.local/bin"
```

### Configuration Guide

The `.cmd.yaml` file is the heart of your CLI application. It defines your commands, arguments, flags, and their behaviors.

#### Basic Structure

```yaml
name: "mycli"                    # Name of your CLI binary
description: "My CLI tool"       # Description shown in help
version: "1.0.0"                # Version of your CLI
author: "Your Name"             # Author information
license: "MIT"                  # License information

# Dependencies required to run your CLI
requires:
  node: ">=16.0.0"              # Specify version constraints
  python3: "*"                  # Any version is acceptable
  go: "^1.22.0"                # Major version must match

# Files to bundle with your CLI
includes:
- "./scripts/helper.sh"
- "./config/default.json"

# Command definitions
commands:
- name: greet                   # Command name (invoked as: mycli greet)
  description: "Greet user"     # Command description
```

#### Defining Arguments

Arguments can be defined with various types and validations:

```yaml
args:
  vars:
  - name: username              # Argument name
    type: string               # Type: string, number, boolean
    description: "Username"    # Description for help text
    required: true            # Is this argument required?
    default: "guest"          # Default value if not provided

  - name: age
    type: number
    description: "User age"
    constraints:              # Validation constraints
      min: 18
      max: 100
    depends-on:              # Dependency relationships
    - name: username         # Only valid if username is provided

  rules:
    minimum-n-args: 1        # Minimum number of arguments
    maximum-n-args: 2        # Maximum number of arguments
```

#### Defining Flags

Flags provide optional modifications to command behavior:

```yaml
flags:
- name: verbose               # Flag name (--verbose)
  shorthand: v               # Short form (-v)
  type: boolean             # Type: boolean, string, number
  description: "Enable verbose output"
  default: false            # Default value

- name: format
  shorthand: f
  type: string
  description: "Output format"
  default: "text"
  conflicts-with:           # Mutually exclusive flags
  - json
  - yaml

- name: count
  shorthand: c
  type: number
  description: "Number of iterations"
  default: 1
```

#### Command Execution

You can specify how commands are executed using the `start` field:

```yaml
commands:
- name: process
  description: "Process files"
  start: |
    if [ "${flags.verbose}" = "true" ]; then
      echo "Processing ${args.filename}..."
    fi
    python3 ./scripts/process.py
```

#### Version Constraints

The `requires` field supports various version constraint operators:

- `*`: Any version
- `^`: Major version must match (^1.2.3 allows 1.x.x)
- `~`: Minor version must match (~1.2.3 allows 1.2.x)
- `>`: Greater than
- `<`: Less than
- `>=`: Greater than or equal
- `<=`: Less than or equal
- No operator: Exact match

Example:
```yaml
requires:
  node: ">=14.0.0"
  python3: "^3.8.0"
  go: "~1.22.0"
```

#### Nested Commands

Commands can be nested to create command hierarchies:

```yaml
commands:
- name: user
  description: "User management"
  commands:
  - name: create
    description: "Create user"
    args:
      vars:
      - name: username
        type: string
        required: true
    start: |
      ./scripts/create-user.sh

  - name: delete
    description: "Delete user"
    flags:
    - name: force
      shorthand: f
      type: boolean
      description: "Force deletion"
```

### Environment Variables

When your commands execute, cmdeagle automatically provides arguments and flags as environment variables:

- Arguments: `ARGS_<NAME>` (e.g., `ARGS_USERNAME`)
- Flags: `FLAGS_<NAME>` (e.g., `FLAGS_VERBOSE`)

These can be accessed in your scripts:

```python
# Python example
import os
username = os.environ.get('ARGS_USERNAME')
verbose = os.environ.get('FLAGS_VERBOSE') == 'true'
```

```javascript
// JavaScript example
const username = process.env.ARGS_USERNAME;
const verbose = process.env.FLAGS_VERBOSE === 'true';
```

```bash
# Shell example
username="${ARGS_USERNAME}"
if [ "${FLAGS_VERBOSE}" = "true" ]; then
  echo "Verbose mode enabled"
fi
```

For more examples, check out the sample CLI in the `examples/mycli` directory.

### For Windows

Docs coming soon.


## Distribution

This tool **does not** help you with releasing or codesigning your resulting binary. It's **highly recommendeded** for you and your users' security that you do some sort of codesigning before distributing your CLI. You can use something like [goreleaser](https://goreleaser.com/) to build and distribute your CLI for multiple platforms and handle codesigning. Worth noiting that cmdeagle itself uses this.

## Contributing

This project is under active development. Contributions, ideas, and feedback are welcome! Please open an issue or submit a pull request on the GitHub repository.

## License

[MIT License](LICENSE)
