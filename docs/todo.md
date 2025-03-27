# Todo

For next version v1.0

- Docs
  - [x] Fix placeholder link to reference name key line 35
- [x] favicon for docs
- [x] fix the examples to use proper env variables
- init command
  - [x] Fix template app's config (line 61 cannot unmarshal)to the former.
  - [x] Fix template app's config (line 61 cannot unmarshal)
  - [x] Make build command for template app work without the chmod command
  - [ ] remove comments before release so that they're not visible in source code
- build command
  - [x] need to provide arguments for arch and handle multiple
- CI/CD
  - [ ] version increment and release script
    - modifies version in docs and package.json
    - releases to Go registry
- [ ] test and figure out what to do with existing variable interpolation functionality.
scrap it? keep it? could be good for cross platform.
- [x] completion setting
- [ ] app level metadata settings need to be printed in the help output
- [ ] version needs to be printed in the help output
- [ ] boolean accepted values need to be accepted correctly
- [ ] add install with npm instructions
- [ ] timer needs to have a list command that is not dependent on jq

## v1.1

- [ ] Remove the need for internet access to build
- CI/CD
  - [ ] release and installation with npm
    - docs must be updated to reflect the new installation method
    - must be tested
- [ ] ability to override CLI_BIN_DIR by either...
  - defining an environment variable
  - using the root or subcommand level 'env' setting in your '.cmd.yaml'
  - then update the docs to reflect the new setting (end of Building Your CLI section)
  - will need to figure out a way to make it work across platforms

## v2

- [ ] Switch to config schema where command, arg and flag names become keys in
the config file
  - template config file needs to be updated
  - docs need to be updated
  - cli-examples need to be updated

## v2.1

- Windows support
  - [ ] On Windows, it will be built to the `%LOCALAPPDATA%\Programs` directory.

## v3

## Someday

- [ ] take full advantage of cobra features

TBD

## Done

- [x] command, arg, and flag `name` keys' values as keys in setting arrays
- [x] fix Args variable key parsing
