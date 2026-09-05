# Design of project-todo

For a Computer Science student who loves programing and go school for job. High self-teaching level. The school task is not the core target. And it makes sense for he to create a customized todo-list for different task management about project development and learning some framework. It will focus on Tasks-Dependency Tree and energy-requirement about a task, which makes clean for him to manage a lot of tasks on school (which not suitable for project development.).

## Data FlowChart

```mermaid
graph TB
    search --> pwd{"pwd/./project-todo/"}
    pwd --> |"Not Found"|cache
    pwd --> |"Found"|read --> fzf --> less_open
```

- Where cache is a file store previous visit todo dir

```txt
/path/to/project1/.project-todo
/path/to/project2/.project-todo
...
```

---

## Energy Consumption Level

- High

Task required focusing. You should ensure you are clear mind and environment is quite & possible to deep work.

- Medium

Task required some energy but you can do it on a noise environment. You can do it on meet or class.

- Low

Just some click or copy-paste task which can be done almost everywhere.

---

## Task ID and Dependencies

Some task has the dependencies which means you should do the previous task so that this task is available.

TODO: Write about its' build and management.

---

## DDL List

For student who develop projects and program everyday. There are some school task has the ddl. so this project will provides a ddl list to remaind user that something must be done for school. This is globally stored on \$XDG_STATE_DIR/project-todo/. And run `todo ddl` command for manage them.

---

## Analysis

Consider an `atomic task` which should not be split into several pieces or more tasks.

For example, go to F2. Go to bed. Or have a drink.

It makes no sense to split these tasks.

However, for most cases, an atomic task makes barely progress.

So we should combines several atomic task by specific sequence to make a real thing done.

We call this structure as `work`.

A work contains not just atomic tasks, it also contains `resources`, `estimated time`, `workflow` etc...

This system should not manage atomic tasks but work.

Because human will use different strategy when doing a same work.

It should flexible and readable.

---

## Core Structure

### System Domain

For this system, there is no need to use field Context and links or definition of done. It makes no sense for the system (it can't parse it). All it's function as below.

- Dependency Graph: make all work provided is really executable
- Energy Routing: Filter the work by high, medium, low energy cost for different occasions
- Life Cycle State: Update `StatusTODO`, `StatusDOING`, `StatusDONE`
- Document Pointer: Store the path to the document for specific work

### Human Domain

- Context & Links: Document reference, manual links, the golang interface just care about how the file be rendered
- Atomic Step List (Workflow): The execution sequence
- Definition of Done: The criteria for human confirm if a work is done.

### Code

See [demo.go](./demo.go)
