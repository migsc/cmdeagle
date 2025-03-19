# cmdeagle

_**WARNING:** This is very much a work in progress, but we're close to releasing a stable version 1. Check back soon. Feedback is also much appreciated._

A language-agnostic CLI application build tool that allows you to create cross-platform applications written in any programming language of your choice.

## Features

- Define commands and subcommands from a single YAML configuration file.
- Reuse your existing external scripts and binaries to build commands.
- Define, parse, and validate your arguments and flags into environment variables your scripts can use.
- Create build steps and declare data/assets to install your commands.

## Quick Start

The easiest way to work with cmdeagle right now is by installing it with [Golang's package manager](https://go.dev/doc/install) on a unix-like system (macOS, Linux, etc). More platforms will be supported soon.

### 1) Install it with Go's package manager

```sh
go install github.com/migsc/cmdeagle@latest
```

Go 1.23.2 or later is required. Install it from [Golang's website](https://go.dev/doc/install).
<!-- - Node.js (v16.17.0+) -->
<!-- - Node.js Package Manager (npm) -->

### 2) Initializing a starter project

```sh
cmdeagle init mycli
```

Where `mycli` is the name of the executable file you want to build. By default, the binary will be named after the directory you run the command from. You can [change this later](#cli-level-name-key).

### 3) Building the CLI:

To build a binary for your CLI, run the following command in the root directory of your project:

```sh
cmdeagle build
```

On macOS and Linux, the binary will be built to either the `./usr/local/bin` directory if you have the necessary write permissions, or the `~/.local/bin` directory if you don't.

The binary will only run on your current operating system and architecture. However, you can target other platforms using specific flags. See the [reference](#reference) for more details.

Currently, internet connectivity is required for the Go build process to resolve and download modules for the wrapper application that `cmdeagle` generates to bundle your scripts and assets. If all dependencies are cached, you can build the project offline. We plan to improve this process in future releases to further reduce the need for internet connectivity during builds.

The sample scripts are defined in several different languages to help you get started. For Python, JavaScript, and other interpreted languages, the scripts are bundled together with the executable file, thanks to the file paths defined in the `includes` setting:

```sh
includes:
- "./greet.sh"
- "./greet.js"
- "./greet.py"
```

This is useful for bundling scripts, static assets, media, configuration files, data files, etc.

However, if you're using Go or another compiled language, bundling executable binaries into your `cmdeagle`-built binary is not allowed for security reasons. Instead, you must use the `CLI_BIN_DIR` variable within the `build` script setting at a subcommand or root level of your `.cmd.yaml`:

```yaml
build: |
  go build -o $CLI_BIN_DIR/mycli-go-binary greet.go
```

Here, Golang's `go build` command is used to compile the `greet.go` file and write the resulting executable to the directory specified by the `$CLI_BIN_DIR` environment variable.

Note that `$CLI_BIN_DIR` is an environment variable defined by `cmdeagle` within the shell where the `build` script runs. It points to the directory where the executable file should be built, matching the directory where the CLI's executable file will be written. Unfortunately, it's not yet possible to override this to point to a different directory, but we plan to add this feature in a future release.

### 4) Running your CLI

You can run the executable file from your current working directory or from anywhere on your system if you add its directory to your system's `PATH` variable.

The `greet` command is a sample subcommand that you can use as a starting point. It's defined in the `.cmd.yaml` file and is configured to run the sample scripts in the project.

Let's test it out:

```sh
mycli greet cmdeagle 2 --uppercase --repeat 3
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
HELLO CMDEAGLE! YOU ARE 2 YEARS OLD.
```

You can get more information about the `greet` command by running:

```sh
mycli help greet
```

The name and age arguments, along with the `--uppercase` and `--repeat` flags, are defined in the `.cmd.yaml` file. Have a look and read the comments to learn how each configuration key works. You can also read the [reference](#reference) for more details on how to define your own commands, flags, and arguments.

For now, let's focus on the `start` script defined for the `greet` subcommand:

```yaml
  start: |
    if [ "${FLAGS_USE_PYTHON}" = "true" ]; then
      python3 greet.py
    elif [ "${FLAGS_USE_JS}" = "true" ]; then
      node greet.js
    elif [ "${FLAGS_USE_GO}" = "true" ]; then
      ./{{name}}-go-binary
    else
      sh greet.sh
    fi
```

Note that it runs the `mycli-go-binary` mentioned before if the `--use-go` flag is passed. The flag could have easily been handled within the code of the `greet.py` and `greet.js` files, but we're doing it this way here to demonstrate a self-contained example of what a `start` script is capable of doing.

Let's invoke the built-in `help` command now:

```sh
mycli help

Usage:
  mycli [flags]
  mycli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  greet       Greet the user.
  help        Help about any command

Flags:
  -h, --help   help for mycli

Use "mycli [command] --help" for more information about a command.
```

Where `mycli` was the name you created in [step 2](#_2-initialize-a-cli-starter-project-named-mycli).

The `completion` command will generate a script for your CLI to use in your shell. This is made possible because `cmdeagle` uses [Cobra](https://github.com/spf13/cobra) under the hood. <!--IMPLEMENT--> You can turn this off by setting the `completion` setting to `false` at the root level of the `.cmd.yaml` file.

Currently, `cmdeagle` primarily uses Cobra for parsing arguments, flags, and subcommands. While we don't yet take full advantage of all the rich features Cobra provides, we plan to integrate more of these capabilities in future releases to enhance the functionality and flexibility of your CLI applications.


## Reference

This section provides detailed documentation for all configuration options, commands, and features available in cmdeagle. Use this reference to understand how to configure your CLI application, define commands and subcommands, set up arguments and flags, and implement the various lifecycle scripts that power your CLI's functionality.

The reference is organized by topic, starting with the configuration structure and moving through each aspect of CLI development with cmdeagle. Each section includes examples and explanations to help you implement the features in your own projects.

### Configuring your CLI 

Your CLI's schema is defined in a [YAML](https://en.wikipedia.org/wiki/YAML#cite_note-19) file named `.cmd.yaml`. There, you define your CLI's basic information and metadata, top-level command and sub commands, validation and parsing rules for your arguments and flags, build steps, bundled assets, and more.

#### Configuration structure

The `.cmd.yaml` file has a tree-like structure that mirrors the command hierarchy of your CLI application:

- The **root level** contains both application-wide configuration (like `completion`) and configuration for the root command itself (the command invoked when users run your CLI without any subcommands)
- **Subcommands** are defined under the `commands` setting and can have their own subcommands, forming a hierarchical tree
- Each command node (root or subcommand) can have its own configuration for arguments, flags, lifecycle scripts, etc.

This hierarchical structure allows you to organize complex CLIs with multiple levels of commands while maintaining a clear and logical configuration.

#### Application and root-level configuration

The root level of your `.cmd.yaml` file contains two types of configuration:

1. **Application-level configuration**: Settings that apply to the entire CLI application
2. **Root command configuration**: Settings that define the behavior of your CLI when invoked without subcommands. These same settings are available to all subcommands.

These following settings only apply at the application level:

##### Basic information and metadata settings

These settings define the core identity and metadata of your CLI application. They're primarily used to populate help text of your CLI for your users.

```yaml
# Basic CLI identity and metadata
name: mycli                                # Name of your CLI executable
description: "A tool for managing widgets" # Short description shown in help text
version: "1.0.0"                           # Version number shown with --version flag
author: "Jane Doe <jane@example.com>"      # Creator information
license: "MIT"                             # License information
```

All these settings are optional except for `name`, which defaults to your project directory name if not specified. They help users understand what your CLI does, who created it, and under what terms it can be used.

It's worth noting that the number defined in the `version` setting is displayed when users run your CLI with the `--version` flag.

##### Other application-level settings

###### `completion` setting

Controls whether the built-in command completion functionality is enabled. Defaults to `true`.

```yaml
completion: true  # Enable shell completion support
```

#### Configuring commands

Commands are the core building blocks of your CLI application. Each command (whether the root command or a subcommand) can be configured with various options that define its behavior, arguments, flags, and execution logic. This section covers all the configuration options available for commands at any level in your command hierarchy.

The following configuration keys can be used within any command definition, including the root command and all subcommands:

##### `commands` setting

At the root level, the `commands` setting defines the top-level subcommands of your CLI application. These are the commands 
that users can run directly after your CLI name.

Each command in defined in the `commands` (or the single one at the root level) defines subcommands for either the root 
command or another subcommand. It contains an array of command objects, each representing a subcommand with its own configuration.

```yaml
commands:
  - name: subcommand1
    description: "First subcommand"
    # other subcommand configuration...
    
  - name: subcommand2
    description: "Second subcommand"
    # other subcommand configuration...
    
    # Nested subcommands
    commands:
      - name: nested
        description: "A nested subcommand"
        # nested subcommand configuration...
```

With this configuration, users can run commands like:
```sh
mycli subcommand1
mycli subcommand2
mycli subcommand2 nested
```

Each command defined here can have its own configuration including [arguments](#arguments-and-flags), [flags](#arguments-and-flags), [lifecycle scripts](#command-lifecycle-configuration), and even [nested subcommands](#top-level-command-and-sub-commands) (using their own `commands` setting).

The `commands` setting is how you build the command tree structure of your CLI application, starting from these top-level commands. In the next section, we'll look at how to define subcommands and their configuration. For the sake of brevity, we'll use the term "command" to refer to both top-level command and subcommands.

#### Command lifecycle configuration

Commands in cmdeagle have a well-defined lifecycle with specific phases that you can hook into to customize behavior. These lifecycle hooks allow you to execute code at different stages of command execution, from validation and preprocessing to the main execution and cleanup. By configuring these lifecycle scripts, you can create sophisticated command behaviors while maintaining a clean separation of concerns.

##### Build Time Lifecycle

During the build phase (`cmdeagle build`), the following steps occur:

```mermaid
flowchart TD

    dev[DEVELOPER]
    user[USER]

    dev --> |Runs **cmdeagle build** command in developer environment| requires1

    requires1(Requirements Step)-->|Described dependencies and their versions checked for existence on the system| build(build)
    build(Build Step) -->|Build script runs and does any necessary compilation| include

    include(Bundling Step) --> |Assets and data bundled together and embedded into a single binary executable| exe
    
    user -->|Executes CLI application in runtime environment| exe
    exe[EXECUTABLE]

```

The settings you define to control this lifecycle are:

- [`requires`](#requires-setting) - Defines dependencies that must be present on the system for your command to run.
- [`build`](#build-setting) - Defines a script that runs during the build phase to compile or prepare your command.
- [`include`](#include-setting) - Defines files that should be bundled with your CLI application.

##### Runtime Lifecycle

When a user runs your CLI application, the following steps occur:

```mermaid
flowchart TD

    user[USER]

    user -->|Executes CLI application in runtime environment| exe
    exe[EXECUTABLE]

        
    exe --> fork_install{First time running?}
    fork_install -->|Yes| install
    install(Install Step) -->|Assets and data extracted from executable| requires2
    fork_install -->|No| requires2 
    requires2[Requirements Step] -->|Dependencies checked again in the runtime environment| validate
    validate(Validation Step) -->|Validation script is run and success is determined by lack of non-zero exit code| fork_start{Success?}
    fork_start --> |Yes| start(Application runs the **start** script)
    fork_start --> |No| exit(Application exits)


```

The settings you define to control this lifecycle are:

- [`requires`](#requires-setting) - Defines dependencies that must be present on the system for your command to run.
- [`validate`](#validate-setting) - Defines a script that runs at runtime before the main command execution.
- [`start`](#start-setting) - Defines the main script that runs when your command is executed.


Let's look at each of these settings in more detail in the following sections.

###### `requires` setting

The `requires` setting specifies dependencies that must be present on the system for your command to run. These dependencies are checked twice: once during build time and again at runtime.

```yaml
requires:
  node: ">=14.0.0"
  python3: "^3.8.0"
  go: "~1.22.0"
```

Each dependency is specified as a key-value pair where:
- The key is the name of the dependency (executable that should be available on the PATH)
- The value is a version constraint string that specifies what versions are acceptable

cmdeagle will check if these dependencies exist and meet the version requirements before proceeding with the build or execution.

**How version checking works:**
cmdeagle attempts to check the version by executing the dependency command with common version flags (like `--version`, `-version`, `--v`, `-v`, etc.) and extracting the version number from the output. This approach has limitations:
- The executable must be in the system's PATH
- The command must support one of the common version flags
- The version format must be recognizable

If cmdeagle cannot determine the version, it will only verify that the command exists but won't validate the version constraint.

If you encounter issues with version checking for a specific dependency, please [create an issue on GitHub](https://github.com/migsc/cmdeagle/issues) so we can improve support for that dependency.

**Available comparison operators:**
- `*` - Any version is acceptable (just checks if the command exists)
- `^` - Compatible with major version (e.g., `^1.2.3` matches any `1.x.y` version)
- `~` - Compatible with major and minor version (e.g., `~1.2.3` matches any `1.2.x` version)
- `>` - Greater than the specified version
- `<` - Less than the specified version
- `>=` - Greater than or equal to the specified version
- `<=` - Less than or equal to the specified version
- No operator - Exact version match (e.g., `1.2.3` only matches version `1.2.3`)

Example:
```yaml
requires:
  node: ">=16.0.0"  # Node.js 16.0.0 or higher
  python3: "*"      # Any version of Python 3
  go: "~1.22.0"     # Any 1.22.x version of Go
  ruby: "^3.0.0"    # Any 3.x.y version of Ruby
  docker: "24.0.5"  # Exactly version 24.0.5 of Docker
```

###### `build` setting

The `build` setting defines a script that runs during the build phase to compile or prepare your command. This script is executed when you run `cmdeagle build`.

```yaml
build: |
  go build -o $CLI_BIN_DIR/mycli-go-binary greet.go
```

The build script is useful for:
- Compiling source code into executables
- Generating assets or configuration files
- Preparing resources needed by your command

Environment variables available during the build script include:
- `$CLI_BIN_DIR`: The directory where binaries should be placed
- `$CLI_NAME`: The name of your CLI application
- `$CLI_DATA_DIR`: The directory where data files will be installed

###### `include` setting

The `include` setting defines files that should be bundled with your CLI application. These files are embedded into the executable during the build phase and extracted when the user first runs your CLI.

```yaml
includes:
- "./greet.sh"
- "./greet.js"
- "./greet.py"
```

This is useful for bundling:
- Scripts in interpreted languages
- Configuration files
- Static assets
- Data files

> **Note:** For security reasons, executable binaries cannot be bundled. Instead, use the `build` script to compile and place executables in the appropriate location.

###### `validate` setting

The `validate` setting defines a script that runs at runtime before the main command execution. It's used to validate arguments, flags, and other conditions before proceeding with the command.

```yaml
validate: |
  if [ -z "$ARGS_FILENAME" ]; then
    echo "Error: Filename cannot be empty" >&2
    exit 1
  fi
  
  if [ ! -f "$ARGS_FILENAME" ]; then
    echo "Error: File does not exist: ${ARGS_FILENAME}" >&2
    exit 1
  fi
```

The validate script should:
- Return a non-zero exit code if validation fails (required for the command to fail)
- Output error messages to stderr (optional but recommended)
- Return 0 (success) if validation passes

If the validation script fails, the command execution will be aborted, and a generic error message will be displayed to the user after whatever output the validation script may have produced.

All [arguments](#arguments-and-flags) and [flags](#arguments-and-flags) defined for the command are available as environment variables within the validation script. See [Using Environment Variables](#using-environment-variables) for details on how to reference these values.

###### `start` setting

The `start` setting defines the main script that runs when your command is executed. This is the core functionality of your command.

```yaml
start: |
  if [ "$FLAGS_USE_PYTHON" = "true" ]; then
    python3 greet.py
  elif [ "$FLAGS_USE_JS" = "true" ]; then
    node greet.js
  elif [ "$FLAGS_USE_GO" = "true" ]; then
    ./{{name}}-go-binary
  else
    sh greet.sh
  fi
```

The start script:
- Has access to all [arguments](#arguments-and-flags) and [flags](#arguments-and-flags) as environment variables
- Can use any files that were included with your command
- Can invoke other executables or scripts
- Is responsible for the main functionality of your command

See [Using Environment Variables](#using-environment-variables) for details on how to reference argument and flag values within your scripts.

#### Arguments and flags

Arguments and flags are the primary ways users interact with your CLI application. cmdeagle provides a robust system for defining, validating, and accessing these inputs in your command scripts.

##### Understanding the difference

- **Arguments** are positional inputs that come after the command name. They are ordered and their position matters.
- **Flags** are named inputs that can appear in any order, typically prefixed with `--` (or `-` for shorthand versions).

cmdeagle assumes your arguments are positional, and the order of your arguments in the configuration file determines their order in the command line. Flags do not have this positional behavior and can be provided in any order.

##### Defining arguments

Arguments are defined in the `args` array of a command. Each argument has several properties that control its behavior:

```yaml
args:
- name: duration
  type: string
  pattern: ^((\d+h)?(\d+m)?(\d+s)?)$|^(\d+)$
  description: "The duration of the timer in the format of #h#m#s"
  required: false
```

##### Defining flags

Flags are defined in the `flags` array of a command. They have similar properties to arguments but with some additional options:

```yaml
flags:
- name: name
  shorthand:
  - n
  type: string
  description: "A name to save the timer under to be able to recall it later."
  pattern: ^[a-zA-Z0-9_-]+$
  required: false
```

##### Common properties for arguments and flags

###### `name` setting

Both arguments and flags must have a `name` setting. This name is used to identify the input and is also how you'll reference its value in your scripts.

```yaml
name: duration
```

###### `description` setting

The description provides a short explanation of the argument or flag that will be displayed in help text.

```yaml
description: "The duration of the timer in the format of #h#m#s"
```

###### `type` setting

The type defines how cmdeagle will parse and validate the input value. Currently supported types include:

- `string`: Text input (default)
- `number`: Numeric input (integers and decimals)
- `boolean`: True/false values (for flags only)

cmdeagle will attempt to parse the input value according to the specified type. If the input value cannot be parsed into the specified type, the argument or flag will be considered invalid and the command will fail, similar to how the `validate` script works.

The default behavior for `number` arguments is to only allow decimal numbers. For `boolean` flags, the input value is case-insensitive and can be either `true`, `false`, `1`, `0`, `yes`, `no`, `y`, `n`, `on`, or `off`.

Note that in your script, these will still be available via environment variables as strings. All cmdeagle is doing for you is validating the input value according to the specified type.

###### `required` setting

Specifies whether the argument or flag must be provided. Defaults to `false` for flags and `true` for arguments.

```yaml
required: true
```

###### `pattern` setting

A regular expression pattern that the input value must match to be considered valid. If the regular expression fails to match, the argument or flag will be considered invalid and the command will fail, similar to how the `validate` script works.

```yaml
pattern: ^((\d+h)?(\d+m)?(\d+s)?)$|^(\d+)$
```

###### `default` setting

The default value to use if the argument or flag is not provided.

```yaml
default: "World"
```

##### Flag-specific properties

###### `shorthand` setting

Defines one or more single-character aliases for the flag. Users can use these with a single dash (e.g., `-n` instead of `--name`).

```yaml
shorthand:
- n
- d
```

###### `conflicts-with` setting

Specifies other flags that cannot be used together with this flag.

```yaml
conflicts-with:
- uppercase
```

###### `depends-on` setting

Specifies other arguments or flags that must be provided when this one is used.

```yaml
depends-on:
- name: name
```

##### Using environment variables

When your command runs, all arguments and flags are made available as environment variables that your scripts can access. This makes it easy to use input values in any programming language.

For arguments, the environment variable format is:
```
ARGS_NAME
```

For flags, the environment variable format is:
```
FLAGS_NAME
```

For example, in a shell script:

```bash
echo "Hello, ${ARGS_NAME}!"
if [ "${FLAGS_UPPERCASE}" = "true" ]; then
  echo "UPPERCASE MODE ENABLED"
fi
```

In a Node.js script:

```javascript
const name = process.env['ARGS_NAME'];
const uppercase = process.env['FLAGS_UPPERCASE'] === 'true';
```

In a Python script:

```python
import os
name = os.environ.get('ARGS_NAME')
uppercase = os.environ.get('FLAGS_UPPERCASE') == 'true'
```

##### Basic built-in validations

In addition to the (command-level `validate`)[#validate-setting] script, cmdeagle performs automatic validation based on the properties you define:

1. Type checking (string, number, boolean)
2. Required field validation
3. Pattern matching (if a pattern is provided)
4. Conflict and dependency validation

It's recommended to make the most of these built-in validations and piggyback off them with your `validate` script for more complex requirements. It's worth mentioning that the built-in validations are checked first, so if they fail, the `validate` script will not be run.

##### Example of complete argument and flag configuration

Here's a comprehensive example showing various argument and flag configurations:

```yaml
args:
- name: name
  type: string
  description: "Name to greet"
  default: "World"
  required: true

- name: age
  type: number
  description: "Age of the user"
  depends-on:
  - name: name

flags:
- name: uppercase
  shorthand: u
  type: boolean
  description: "Convert greeting to uppercase"
  default: false
  conflicts-with:
  - lowercase

- name: lowercase
  shorthand: l
  type: boolean
  description: "Convert greeting to lowercase"
  default: false
  conflicts-with:
  - uppercase

- name: repeat
  shorthand: r
  type: number
  description: "Repeat the greeting n times"
  default: 1
```

This configuration would allow commands like:
```
mycli greet John 25 --uppercase --repeat 3
```

### Using the `cmdeagle` CLI

The `cmdeagle` CLI is used to initialize, build, and manage your CLI application. It's fairly simple and has only a few commands.

#### `init` command

The `init` command creates a new CLI project with a basic structure and example commands.

```sh
cmdeagle init [name]
```

**Parameters:**
- `name` - The name of your CLI application (optional, defaults to the current directory name)

**Examples:**

Create a new CLI named "mycli":
```sh
cmdeagle init mycli
```

Create a CLI in the current directory:
```sh
mkdir my-awesome-cli && cd my-awesome-cli
cmdeagle init
```

This command creates several files:
- `.cmd.yaml` - The main configuration file for your CLI
- Sample greeting scripts in multiple languages (Shell, JavaScript, Python, Go)

#### `build` command

The `build` command compiles your CLI application based on the configuration in your `.cmd.yaml` file.

```sh
cmdeagle build [flags]
```

**Flags:**
- Currently, the build command supports building for your current platform. Cross-platform support is coming soon.

**Examples:**

Build your CLI:
```sh
cmdeagle build
```

After building, your CLI will be available in:
- On macOS/Linux: `/usr/local/bin` or `~/.local/bin`
- On Windows: `%LocalAppData%\Programs\mycli\bin`

#### `completion` command

The `completion` command generates shell completion scripts to enable tab completion for the `cmdeagle` commands.

```sh
cmdeagle completion [shell]
```

**Parameters:**
- `shell` - The shell to generate completion for (bash, zsh, fish, powershell)

**Examples:**


Generate bash completions:
```sh
# Create completion directory if it doesn't exist
mkdir -p ~/.bash_completion.d
cmdeagle completion bash > ~/.bash_completion.d/cmdeagle

# Add to your ~/.bashrc to load completions
echo 'source ~/.bash_completion.d/cmdeagle' >> ~/.bashrc
source ~/.bashrc
```

Generate zsh completions:
```sh
# Create completion directory if it doesn't exist
mkdir -p ~/.zsh/completion
cmdeagle completion zsh > ~/.zsh/completion/_cmdeagle

# Add to your ~/.zshrc to load completions
echo 'fpath=(~/.zsh/completion $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc
source ~/.zshrc
```

Generate fish completions:
```sh
# Fish automatically loads completions from this directory
mkdir -p ~/.config/fish/completions
cmdeagle completion fish > ~/.config/fish/completions/cmdeagle.fish
```

Generate PowerShell completions:
```powershell
# Create a directory for the completion script
mkdir -p ~/Documents/PowerShell/
cmdeagle completion powershell > ~/Documents/PowerShell/cmdeagle.ps1

# Add to your PowerShell profile to load completions
echo '. ~/Documents/PowerShell/cmdeagle.ps1' >> $PROFILE
```

Note: While PowerShell completion is listed as an option (because it's built into Cobra), Windows support for cmdeagle is not yet fully implemented. We plan to add comprehensive Windows support in a future release.

#### `help` command

The `help` command displays help information about available commands and their usage.

```sh
cmdeagle help [command]
```

**Parameters:**
- `command` - The command to get help for (optional)

**Examples:**

Get general help:
```sh
cmdeagle help
```

Get help for the init command:
```sh
cmdeagle help init
```

Get help for the build command:
```sh
cmdeagle help build
```

You can also use the `-h` or `--help` flag with any command to see its help information:
```sh
cmdeagle init --help
cmdeagle build -h
```

The help command provides detailed information about command usage, available flags, and examples to guide you through using cmdeagle effectively.