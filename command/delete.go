package command

import (
	"github.com/guntutur/todo-cli/datastore"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{Use: "delete", Short: "delete todo, will not be shown afterwards, but are kept in memory datastore", Run: DeleteTodo}

func init() {
	deleteCommand.Flags().String("id", "", "id of the todo you want to complete")
	deleteCommand.MarkFlagRequired("id")
	ContextCommand.AddCommand(deleteCommand)
}

func DeleteTodo(cmd *cobra.Command, args []string) {
	id := cmd.Flag("id").Value.String()
	//datastore.DeleteTodo(id)
	datastore.DeleteTodoCDB(id)
}
