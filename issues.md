# Issues

These are meant to transfer into Github Issues later.


- Print app-level metadata in help output
    - [feature] [help] [usability]
    - The CLI's help output should display application-level metadata (such as author, license, etc.) to provide users with more context.

- Don't print Version, Author or License when they are empty
    - [bug]

- Print version in help output
    - [feature] [help] [usability]
    - The CLI's help output should include the current version number, making it easy for users to verify which version they are running.

- Accept all boolean values correctly
    - [bug] [flags] [usability]
    - Boolean flags should accept all common representations (true/false, t/f, 1/0, yes/no, y/n) for a smoother user experience.

- Add install with npm instructions to documentation
    - [docs] [npm] [installation]
    - The documentation should include clear instructions for installing the CLI using npm.

- Timer: add a list command that does not depend on jq
    - [feature] [timer] [usability]
    - The timer feature should include a list command that works without requiring jq, to improve accessibility and reduce dependencies.

- Release and installation with npm (docs and testing)
    - [release] [npm] [docs] [testing]
    - Update documentation to reflect the new npm installation method and ensure the release and installation process is thoroughly tested.

- Allow overriding CLI_BIN_DIR
    - [feature] [env] [cross-platform]
    - Users should be able to override CLI_BIN_DIR by:
        - Defining an environment variable
        - Using the root or subcommand level 'env' setting in .cmd.yaml
      Update the documentation accordingly and ensure this works across platforms.

- Add "What's Next" section to Quick Start in docs
    - [docs] [onboarding] [reference]
    - At the end of the "Quick Start" section, add a "What's Next" section that links to key reference sections (such as configuration and CLI usage) to help new users discover next steps.

- `init` command: allow creating a bare template
    - [feature] [init] [major]
    - The `init` command should support creating a minimal/bare template project, without sample commands or scripts, for advanced users who want to start from scratch.

- `init` command: properly skip interactive prompts when arguments are provided
    - [feature] [init] [minor] [usability]
    - If all required arguments are passed to the `init` command, it should skip interactive prompts and use the provided values, streamlining automation and scripting.
    - But if only some of the arguments are passed, we need to still launch the interactive initialization.

- Add support for  updating shell files on users' behalf for autocompletion of commands
    [feature] [minor]

- Add a new `debug` command
    - [feature] [debug] [major]
    - Introduce a `debug` command to help users troubleshoot their CLI applications, inspect configuration, and diagnose issues.

- Migrate to new config format (v2)
    - [breaking] [config] [v2]
    - Redesign the configuration schema so that command, arg, and flag names become keys in the config file. Update templates, docs, and CLI examples accordingly.

- Rich TUI for built-in `completion` command (future)
    - [feature] [completion] [TUI] [future]
    - Enhance the built-in CLI's `completion` command to provide a rich terminal UI (TUI) using react-tea (React renderer for Bubble Tea), making it more interactive and user-friendly.

- Terminal-based browser for `help` command (future)
    - [feature] [help] [TUI] [future]
    - Upgrade the `help` command to launch a terminal-based browser for documentation, improving discoverability and navigation for users.

- Update `init` command help output to accurately describe argument handling
    - [bug] [init] [help] [usability]
    - The help output for the `init` command should clearly and accurately explain how arguments are handled, including when interactive prompts are skipped, how defaults are applied, and what happens when only some arguments are provided. The current help text is misleading or incomplete.

- Provide additional starter templates for initialization (Node.js/TypeScript, Python, Go, etc.)
    - [feature] [init] [templates] [enhancement] [low-priority]
    - Offer users a choice of starter templates when initializing a new CLI project, such as Node.js with TypeScript/JavaScript, Python, Go, and others. This will make it easier for users to get started in their preferred language and ecosystem.

- Create a new repo or directory in this repo to example CLIs built with cmdeagle
    - [documentation] [usability]
    - Should add the timer CLI that I created

- Remove the need for internet access to build
    - [feature] [build] [offline] [v1] [minor]
    - Refactor the build process so that, once dependencies are cached, users can build their CLI projects without requiring an internet connection.

- Windows support: build to %LOCALAPPDATA%\Programs directory
    - [feature] [windows] [cross-platform] [v2] [minor]
    - Ensure that on Windows, the CLI is built and installed to the %LOCALAPPDATA%\Programs directory, matching platform conventions.

- Take full advantage of Cobra features
    - [enhancement] [cobra] [future]
    - Refactor and extend the CLI to leverage more of Cobra's advanced features, such as richer argument/flag validation, command aliases, and improved help output.