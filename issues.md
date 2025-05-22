# Issues

These are meant to transfer into Github Issues later.


- Print app-level metadata in help output
    - [feature] [help] [usability]
    - The CLI's help output should display application-level metadata (such as author, license, etc.) to provide users with more context.

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

- `init` command: skip interactive prompts when arguments are provided
    - [feature] [init] [minor] [usability]
    - If all required arguments are passed to the `init` command, it should skip interactive prompts and use the provided values, streamlining automation and scripting.

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