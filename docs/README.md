# Quick Start

### Install with npm

```sh
npm install -g cmdeagle
```

Prerequisites: 
<!-- - Node.js (v16.17.0+) -->
- Node.js Package Manager (npm)

### Initialize a CLI starter project named yourcli

```bash
cmdeagle init yourcli
```

### Build the CLI from the root directory of your project:

```bash
cmdeagle build
```




### Create a sub command

### Define build script for your command

### Define start script for your command

### Build and run your new sub command

### Defining build scripts for executables

## Installation by Platform

### Linux 

### macOS

### Windows

## Prerequisites
- Install cmdeagle

# Reference

## Using the cmdeagle CLI

### Initializing CLI configuration with the init Command

### Build your CLI with the build Command
- Targeting a Specific Platform

## Configuring Your CLI

### Defining Your CLI's Basic Information and Metadata

### Top-level Command and Sub Commands

#### Defining Command's Basic Information 

#### Defining Command-Level Validation

#### Command Lifecycle Scripts

##### The Command Scripts Lifecycle
- `requires` - Validates dependencies, once before building and once before running your command.
- `build` - Builds your command and create an exectuable file
- `include` - Defines bundled assets to include in your CLI's executable file after building your command
- `start` - Executed when you or your users run your command

##### Defining `build` script

##### Defining `include` script

##### Defining `start` script 

##### Defining Bundled Assets To Include

#### Arguments and Flags

##### Definition

##### Validation and Parsing
- Argument-Only
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
