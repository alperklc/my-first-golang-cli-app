package main

import (
	"github.com/alperklc/my-first-golang-cli-app/cmd"
	"github.com/alperklc/my-first-golang-cli-app/db"
)

func main() {
	database := db.NewDatabase()
	database.OpenDatabase()

	commander := cmd.NewCommander(database)
	commander.RegisterCommands()
	commander.Execute()
}
