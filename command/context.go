package command

import (
	"github.com/guntutur/todo-cli/datastore"
	"log"

	"github.com/spf13/cobra"
)

var ContextCommand = &cobra.Command{}

func init() {
	ContextCommand = &cobra.Command{Use: "todo", Short: "Starting todo apps, using default DB : " + datastore.DbName, Version: "0.0.1"}
}

// Execute is the entry point
func Execute() {
	err := ContextCommand.Execute()
	if err != nil {
		log.Fatal("cannot start todo app - ", err)
	}
}
