# Project-todo

For a Computer Science student who loves programming and go school for job. High self-teaching level. The school task is not the core target. And it makes sense for he to create a customized todo list for different task management about project development and learning some framework. It will focus on Tasks-Dependency Tree and energy-requirement about a task, which makes clean for him to manage a lot of tasks on school (which not suitable for project development.).

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

## USAGE

TODO. I haven't finish this cli, 😃

## Design

see the cli design on [design.md](/design.md)
