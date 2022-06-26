# my-first-golang-cli-app

A CLI app for managing a to-do list, based on sqlite and cobra.

## How to develop

- `go build main.go` builds the application and creates a new binary file.
- `./main` runs the application

## Data structure

```
type Todo struct {
	Id                       int
	Name, Description, Tasks string
}
```

Each todo entry on the database consist of `name`, `description` and `tasks` fields of string type. `tasks` field represents an array of strings, seperated with pipe operator.

## Commands & usage

Local database needs to be initialized upon first use. To this end, please run `init` command.

- `init` initializes a database in the same folder (sqlite-database.db)
- `todo list` lists todos
- `todo new` prompts an interactive user interface, for getting the details of a todo entry.
- `todo get [todo_id]` gets a todo with the given id
- `todo update [todo_id]` updates a todo with the given id
- `todo add-task [todo_id] [task]` adds a task to the todo with the given id
- `todo remove-task [todo_id] [task]` removes a task from the todo with the given id
- `todo delete [todo_id]` deletes a todo with the given id
