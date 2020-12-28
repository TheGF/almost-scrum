# Almost Scrum
Almost Scrum is a simple application to support agile project
management. It aims to be agile and technically simple.

## Basic Principles
- simple storage based on file system
- so simple to be manually hackable (if required)
- integrated with Git
- built for local usage
- comes with web UI
- distributed as one single executable

## Basic concepts
Project management is built on **tasks** and **boards**. 
A task is an activity assigned to one person 
A board is a collection of tasks.

Each new project has a *Backlog* board.



## Command Line Command
Almost Scrum comes with a single executable *ash*
All commands accept two optional parameter 
- -p path: the location of the project
- -d level: debug level

In case no path is provided, ash looks for a project in
the current folder and all its parents.

### Command init
    ash [-p path] init

Create a new project in the current directory or in
the provided path. It fails if a project already exists.

The command first looks for a .git folder in the path.
If the folder is found, it creates a new folder *.ash* and
there it creates the project folders and configuration

If no .git folder is in the path, it creates the project
folders and configuration in the current path

### Command board
    ash [-p path] board

Show the boards and highlight the current one

### Command board new
    ash [-p path] board new [name]

Create a new board

### Command board current
    ash [-p path] board default [filter]

Change the current board

### Command new
    ash [-p path] new [title]

Create a new task.
Optionally the title of the task can be provided. 
When not provided, the user is required to enter it


### Command owner
    ash [-p path] owner [filter]

Change a task owner. Only the current owner is expected to
change reassign a task


### Command edit
    ash [-p path] edit [filter]

Edit a task. A user can edit only tasks he owns. 

### Command mv
    ash [-p path] mv [filter]

Move a task to a different board.


## Command Line Command

