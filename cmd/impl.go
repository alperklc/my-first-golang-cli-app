package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alperklc/my-first-golang-cli-app/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CommandImpl interface {
	Init(cmd *cobra.Command, args []string)
	ListAllTodos(cmd *cobra.Command, args []string)
	NewTodo(cmd *cobra.Command, args []string)
	GetTodo(cmd *cobra.Command, args []string)
	UpdateTodo(cmd *cobra.Command, args []string)
	AddTask(cmd *cobra.Command, args []string)
	RemoveTask(cmd *cobra.Command, args []string)
	DeleteTodo(cmd *cobra.Command, args []string)
}

type data struct {
	logger *log.Entry
	db     db.Database
}

type promptContent struct {
	errorMsg string
	label    string
}

func NewCommandImpl(db db.Database) CommandImpl {
	return &data{
		db:     db,
		logger: log.WithFields(log.Fields{"package": "cmd::impl"}),
	}
}

func (d *data) Init(cmd *cobra.Command, args []string) {
	d.db.CreateTable()
}

func (d *data) ListAllTodos(cmd *cobra.Command, args []string) {
	d.db.DisplayAllTodos()
}

func (d *data) NewTodo(cmd *cobra.Command, args []string) {
	namePromptContent := promptContent{
		"Could not process the entry. Please enter a name for todo.",
		"Please enter a name for todo.",
	}
	name := promptGetInput(namePromptContent)

	descriptionPromptContent := promptContent{
		"Could not process the entry. Please enter a description for the todo entry.",
		"Please enter a description for the todo entry",
	}
	description := promptGetInput(descriptionPromptContent)

	tasksPromptContent := promptContent{
		"Could not process the entry. Please add tasks.",
		"Please add tasks",
	}
	tasks := promptAddTasks(tasksPromptContent)

	d.db.InsertTodo(name, description, tasks)
}

func (d *data) GetTodo(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		d.logger.Fatal("Arguments are missing.")
		os.Exit(1)
	}
	d.db.GetTodo(args[0])
}

func (d *data) UpdateTodo(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		d.logger.Fatal("Arguments are missing.")
	}

	todoFound, errGetTodo := d.db.GetTodo(args[0])
	if errGetTodo != nil {
		d.logger.Fatalf("Could not find todo: %v\n", errGetTodo)
		os.Exit(1)
	}
	newVersionOfTodo := todoFound

	namePromptContent := promptContent{
		"Could not process the entry. Please enter a name for todo.",
		fmt.Sprintf("%s is the name of the todo. Would you like to change it? Please enter a new name for todo. Hit enter to leave it", todoFound.Name),
	}
	newName := promptGetInputWithDefault(namePromptContent, todoFound.Name)
	newVersionOfTodo.Name = newName

	descriptionPromptContent := promptContent{
		"Could not process the entry. Please enter a description for the todo entry.",
		fmt.Sprintf("%s is the description of the todo. Please enter a description for the todo entry. Hit enter to leave it", todoFound.Description),
	}
	newVersionOfTodo.Description = promptGetInputWithDefault(descriptionPromptContent, todoFound.Description)

	d.db.UpdateTodo(todoFound.Id, newVersionOfTodo)
}

func (d *data) AddTask(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		d.logger.Fatal("Arguments are missing.")
		os.Exit(1)
	}

	todoFound, errGetTodo := d.db.GetTodo(args[0])
	if errGetTodo != nil {
		d.logger.Fatalf("Could not find todo: %v\n", errGetTodo)
		os.Exit(1)
	}

	tasksArray := strings.Split(todoFound.Tasks, "|")
	tasksArray = append(tasksArray, args[1])
	newTasks := strings.Join(tasksArray[:], "|")

	d.db.UpdateTodo(todoFound.Id, db.Todo{Id: todoFound.Id, Name: todoFound.Name, Description: todoFound.Description, Tasks: newTasks})
}

func (d *data) RemoveTask(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		d.logger.Fatal("Arguments are missing.")
		os.Exit(1)
	}

	todoFound, errGetTodo := d.db.GetTodo(args[0])
	if errGetTodo != nil {
		d.logger.Fatalf("Could not find todo: %v\n", errGetTodo)
		os.Exit(1)
	}

	tasksArray := strings.Split(todoFound.Tasks, "|")
	newTasksArray := []string{}
	for i := range tasksArray {
		if tasksArray[i] != args[1] {
			newTasksArray = append(newTasksArray, tasksArray[i])
		}
	}
	newTasks := strings.Join(newTasksArray[:], "|")

	d.db.UpdateTodo(todoFound.Id, db.Todo{Id: todoFound.Id, Name: todoFound.Name, Description: todoFound.Description, Tasks: newTasks})
}

func (d *data) DeleteTodo(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		d.logger.Fatal("Arguments are missing.")
		os.Exit(1)
	}

	d.db.DeleteTodo(args[0])
}
