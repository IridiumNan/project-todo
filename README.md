# Project-todo

**Version 0.0.1**

For a Computer Science student who loves programming and high self-teaching level. The school task is not the core target. And it makes sense for he to create a customized todo list for different task management about project development and learning some framework. It will focus on Work-Dependency Tree and energy-requirement about a work, which makes it easier for him to manage a lot of works on school (which not suitable for project development.).

---

## Feature

- Energy Based

This project is the todo work manager with filter (or engine) based on your energy status.

When you are tied or sleepy, you can get the work with tag `EnergyLow`.
When you are fine but the environment is noisy or hard to focus, you can get the work with tag `EnergyMedium`.
When you are able 90% focus on your work, you can get the work with tag `EnergyHigh`.

- Project Based

When you run `todo init` command on your project dir (like git init)
A data dir which store your todo work for this project will be created.
Which means you can manage this by git and carry it around your different machines.

> [!NOTE]
> Also `todo init /path/to/project`

## Quick Start

- create a new directory for todo metadata and context storage.

```bash
cd /path/to/project

todo init
```

- create a new work on this directory
  this CLI will open a context configuration file, you will assign its title, dependencies, context

```bash
todo new
```

- pop a work matches your current energy level

```bash
todo next -e=0 # for energy low

# -e=1 for energy medium

# -e=2 for energy high
```

## Usage

```bash
Usage:
  todo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Open config file with env editor
  help        Help about any command
  init        init a project-todo on current dir
  log         Open the log file with specific program, default use less
  new         Create a new work then push it to todo list with statusTODO
  next        Pop a work on todo list with specific energy level

Flags:
  -h, --help   help for todo

Use "todo [command] --help" for more information about a command.
```

## Design

See the cli design on [design.md](/design.md)

## Todo

The done works analysis will be support on v0.0.2
