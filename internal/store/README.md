# Project-todo

This directory is created by [Project-todo](https://github.com/IridiumNan/project-todo) which is designed to manage the todo works for project.

## Installation

Visit [repo](https://github.com/IridiumNan/project-todo) for details

## Basic USAGE

- init a new todo list on current dir

```bash
todo init

# todo init <path> also support
```

- pop next work

```bash
todo next --energy=low
# available options --energy=low medium high
# run todo next --help for details
```

- push a new work

```bash
todo new
```

- config the global file

```bash
todo config
```

- check logs

```bash
todo log vim # open log file with vim
# You can use anything you like
# less, tail etc...
```

> [!NOTE]
> This program will not make any change to your system, it just create plain files and parse them.
> So you don't need to uninstall it when something go wrong.
