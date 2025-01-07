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

You need to have Go installed to install cmdeagle and build CLIs. You can install Go by following the instructions at https://go.dev/doc/install.

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
PATH="$PATH:PATH_TO_YOUR_CLONE_DIR/cmdeagle/sdk/bin"
```

Where `PATH_TO_YOUR_CLONE_DIR` is the path to the directory you cloned cmdeagle into.

## Usage

1. Initialize a new CLI project by running in a directory of your choice:

```
mdkir ./mycli && cd mycli
cmdeagle init
```

This will create a `.cmd.yaml` file in your project root to define your CLI structure. The schema is not documented yet, but you can read the comments in the file to understand what each field does.

2. Build your CLI:

```
cmdeagle build
```

3. Run your CLI.

### For Linux and macOS:


If you have `./.local/bin` in your `PATH` you can simply run:

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


### For Windows

Docs coming soon.


## Distribution

This tool **does not** help you with releasing or codesigning your resulting binary. It's **highly recommendeded** for you and your users' security that you do some sort of codesigning before distributing your CLI. Some recommendations:

- [goreleaser](https://goreleaser.com/) - Builds and distributes your CLI for multiple platforms and handles codesigning. cmdeagle itself uses this.
- [sigstore](https://sigstore.dev/) - Codesigning and verification.
- [cosign](https://github.com/sigstore/cosign) - Codesigning and verification.

## Contributing

This project is under active development. Contributions, ideas, and feedback are welcome! Please open an issue or submit a pull request on the GitHub repository.

## License

[MIT License](LICENSE)
