/*
Copyright © 2022 Alper Kilci <alperkilci@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/alperklc/my-first-golang-cli-app/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Commander interface {
	Execute()
	RegisterCommands()
}

type d struct {
	logger        *log.Entry
	database      db.Database
	rootCmd       *cobra.Command
	todoCmd       *cobra.Command
	initCmd       *cobra.Command
	listCmd       *cobra.Command
	newCmd        *cobra.Command
	getCmd        *cobra.Command
	updateCmd     *cobra.Command
	addTaskCmd    *cobra.Command
	removeTaskCmd *cobra.Command
	deleteCmd     *cobra.Command
}

func NewCommander(database db.Database) Commander {
	return &d{
		database: database,
		logger:   log.WithFields(log.Fields{"package": "cmd"}),
	}
}

func (d *d) Execute() {
	err := d.rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (d *d) RegisterCommands() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-first-golang-cli-app.yaml)")

	cmdImpl := NewCommandImpl(d.database)

	d.rootCmd = &cobra.Command{
		Use:   "my-first-golang-cli-app",
		Short: "Yet another to-do list application",
		Long:  `Write to-do entries with names, descriptions and tasks.`,
	}
	d.rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	d.todoCmd = &cobra.Command{
		Use:   "todo",
		Short: "A todo is something to do",
	}
	d.initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialise a new todo database and table",
		Run:   cmdImpl.Init,
	}
	d.listCmd = &cobra.Command{
		Use:   "list",
		Short: "See a list of all notes you've added",
		Run:   cmdImpl.ListAllTodos,
	}
	d.newCmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new todo",
		Run:   cmdImpl.NewTodo,
	}
	d.getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a todo",
		Run:   cmdImpl.GetTodo,
	}
	d.updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updates a todo",
		Run:   cmdImpl.UpdateTodo,
	}
	d.addTaskCmd = &cobra.Command{
		Use:   "add-task",
		Short: "Adds a task to the todo",
		Run:   cmdImpl.AddTask,
	}
	d.removeTaskCmd = &cobra.Command{
		Use:   "remove-task",
		Short: "Removes a task from the todo",
		Run:   cmdImpl.RemoveTask,
	}
	d.deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Deletes a todo",
		Run:   cmdImpl.DeleteTodo,
	}

	d.rootCmd.AddCommand(d.initCmd)
	d.rootCmd.AddCommand(d.todoCmd)
	d.todoCmd.AddCommand(d.listCmd)
	d.todoCmd.AddCommand(d.newCmd)
	d.todoCmd.AddCommand(d.deleteCmd)
	d.todoCmd.AddCommand(d.updateCmd)
	d.todoCmd.AddCommand(d.addTaskCmd)
	d.todoCmd.AddCommand(d.removeTaskCmd)
	d.todoCmd.AddCommand(d.getCmd)
}
