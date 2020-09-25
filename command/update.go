package command

import (
	"github.com/guntutur/todo-cli/datastore"
	"github.com/spf13/cobra"
)

var completeCommand = &cobra.Command{Use: "complete", Short: "complete todo", Run: CompleteTodo}

func init() {
	completeCommand.Flags().String("id", "", "id of the todo you want to complete")
	completeCommand.MarkFlagRequired("id")
	ContextCommand.AddCommand(completeCommand)
}

func CompleteTodo(cmd *cobra.Command, args []string) {
	id := cmd.Flag("id").Value.String()
	datastore.CompleteTodo(id)
}
