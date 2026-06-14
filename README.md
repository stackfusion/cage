# Cage

Cage is a developer tool for safely working with untrusted projects on macOS.

Given the recent rise in supply-chain attacks, it’s crucial to implement proper safeguards when executing code, including dependencies, which are primary targets for attackers. The fundamental principle should be as simple as **never run code you didn’t write on your machine.**

When you `cd` into a caged project, Cage reminds and helps you to boot an isolated VM with only that project’s directory mounted. This keeps your host clean and ensures your code inside the VM is sandboxed. _Your code goes in the cage. Everything else stays free._

_The name and principle are inspired by FreeBSD's `jail`._

## Features

[![asciicast](https://asciinema.org/a/R2XWr5UcOgXulprq.svg)](https://asciinema.org/a/R2XWr5UcOgXulprq)

- Each project gets its own VM, provisioned on first boot (customizable script).
- Only project directory mounted; no other host files available.
- [`mise`](https://github.com/jdx/mise) preinstalled in VM; install runtimes there, not on host.
- Shell hook detects caged directory entry and reminds to stay inside VM.
- Zed and VS Code editors connect via SSH remote; language servers, deps, terminals, and tasks run inside VM.

## Installation

```shell
brew install lima
brew tap stackfusion/tap
brew install cage
```

You may also want to trust the `stackfusion/tap`: just follow the output from the `brew` command for more information.

## Initialization

```shell
cage install
```

This creates a VM template and patches your shell rc. It auto-detects your shell (bash, zsh, fish), but you can specify it explicitly.

```shell
cage install --shell fish
```

## Usage

```shell
cd ~/Workspace/some-project
cage
```

**That's it.**

Cage initializes the project, starts the VM, and opens a shell inside it. Use `cage help` for more information.

### Shell Hook

The hook fires automatically when you cd into a caged directory or its children.

```
~ $ cd ~/Workspace/sample-project
cage: caged directory — VM is not running
cage: run `cage` to start and enter the VM, or `cage ack` to suppress this banner
```

After `cage ack`, the loud banner is replaced with a subtle one-liner on subsequent visits.

To manually set up the hook:

```shell
# zsh
echo 'eval "$(cage hook zsh)"' >> ~/.zshrc

# bash
echo 'eval "$(cage hook bash)"' >> ~/.bashrc

# fish
echo 'cage hook fish | source' >> ~/.config/fish/config.fish
```

### The `.cage` File

A minimal config with your project.

```yaml
vm_name: my-project-cage
```

The VM template (CPU, memory, disk, base image, and provisioning script) lives in `~/.config/cage/` and can be customized.

### Editors

Zed and VS Code connect to the VM via SSH, where language servers, terminals, and build tools reside.

```shell
cage zed  # requires zed on PATH
cage code # requires code on PATH and the Remote-SSH extension
```

## Safety Model

The Cage is a workflow, not a security sandbox. The VM is the real boundary; the shell hook and banner are reminders that won’t stop you from running potentially dangerous code locally.

_The goal is to make the right thing easy to do._
