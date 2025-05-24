# Issues

These are meant to transfer into Github Issues later.

- Terminal-based browser for `help` command (future)
    - [enhancement] [help] [TUI] [future]
    - Upgrade the `help` command to launch a terminal-based browser for documentation, improving discoverability and navigation for users.

- Update `init` command help output to accurately describe argument handling
    - [bug] [init] [help] [usability]
    - The help output for the `init` command should clearly and accurately explain how arguments are handled, including when interactive prompts are skipped, how defaults are applied, and what happens when only some arguments are provided. The current help text is misleading or incomplete.

- Provide additional starter templates for initialization (Node.js/TypeScript, Python, Go, etc.)
    - [enhancement] [init] [templates] [enhancement] [low-priority]
    - Offer users a choice of starter templates when initializing a new CLI project, such as Node.js with TypeScript/JavaScript, Python, Go, and others. This will make it easier for users to get started in their preferred language and ecosystem.

- Create a new repo or directory in this repo to example CLIs built with cmdeagle
    - [documentation] [usability]
    - Should add the timer CLI that I created

- Remove the need for internet access to build
    - [enhancement] [build] [offline] [v1] [minor]
    - Refactor the build process so that, once dependencies are cached, users can build their CLI projects without requiring an internet connection.

- Windows support: build to %LOCALAPPDATA%\Programs directory
    - [enhancement] [windows] [cross-platform] [v2] [minor]
    - Ensure that on Windows, the CLI is built and installed to the %LOCALAPPDATA%\Programs directory, matching platform conventions.

- Take full advantage of Cobra features
    - [enhancement] [cobra] [future]
    - Refactor and extend the CLI to leverage more of Cobra's advanced features, such as richer argument/flag validation, command aliases, and improved help output.