# cage

A _(slightly opinionated, but configurable)_ developer tool for safely working with untrusted, unknown, or any other kind of projects on macOS. The core philosophy is simple: **don't run foreign code on your host machine**.

When you `cd` into a caged project, `cage` spins up an isolated Lima VM with only that project's directory mounted. Your host stays clean, and any code (or its dependency) that you run inside VM is isolated.

The name is inspired by FreeBSD's `jail`.

## How It Works

- Each project gets its own Lima VM, provisioned on first use
- Only the project directory is mounted into the VM — nothing else from your host is visible
- `mise` is pre-installed inside the VM; you install runtimes there, not on your host
- A shell hook detects when you enter a caged directory and reminds you to stay inside the VM
- Editors connect via SSH remote (Zed, VS Code) — language servers run inside the VM, not on your host

## Prerequisites

```sh
brew install lima
```

## Installation

```sh
go install github.com/stackfusion/cage@latest
cage install
```

`cage install` writes Lima template and patches your shell rc. It auto-detects your shell (bash, zsh, fish), but you can be explicit:

```sh
cage install --shell fish
```

## Usage

### The Short Version

```sh
cd ~/Workspace/some-project
cage
```

That's it. `cage` checks setup, initializes the project if needed, starts the VM, and drops you into a shell inside it.

### Commands

| Command               | Description                                             |
|-----------------------|---------------------------------------------------------|
| `cage`                | Guided entry point — init, start, and shell in one step |
| `cage init`           | Write a `.cage` config in the current directory         |
| `cage start`          | Create (if needed) and start the Lima VM                |
| `cage stop`           | Stop the VM                                             |
| `cage delete`         | Stop and permanently delete the VM                      |
| `cage shell`          | Open a shell inside the VM                              |
| `cage shell -- <cmd>` | Run a command inside the VM                             |
| `cage zed`            | Open the project in Zed via SSH remote                  |
| `cage code`           | Open the project in VS Code via SSH remote              |
| `cage acknowledge`    | Suppress the entry banner until `.cage` changes         |
| `cage install`        | Write Lima template and patch shell rc                  |
| `cage prune`          | Remove cage config and optionally delete the VM         |

## Shell Hook

The hook fires automatically when you `cd` into a caged directory:

```
~/Workspace $ cd some-project

⚠ cage: caged directory — VM some-project-cage is not running
cage: run cage to start and enter the VM, or cage acknowledge to suppress this
```

Once you run `cage acknowledge`, the loud banner is suppressed and replaced with a subtle one-line indicator on subsequent visits.

To set up the hook manually:

```sh
# zsh
echo 'eval "$(cage hook zsh)"' >> ~/.zshrc

# bash
echo 'eval "$(cage hook bash)"' >> ~/.bashrc

# fish
echo 'cage hook fish | source' >> ~/.config/fish/config.fish
```

## The `.cage` File

A minimal config committed alongside your project:

```yaml
vm_name: my-project-cage
```

The VM template (CPU, memory, disk, base image) lives in `~/.config/cage/lima-template.yaml` and can be edited to suit your needs.

## Editors

Both Zed and VS Code connect to the VM over SSH. Language servers, terminals, and build tools all run inside the VM.

```sh
cage zed  # requires zed on PATH
cage code # requires code on PATH and the Remote-SSH extension
```

## Safety Model

`cage` is not a security sandbox — it is a workflow guardrail. The VM is a real boundary; the shell hook and banner are reminders. The goal is to make it easy to do the right thing by default.
