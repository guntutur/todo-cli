package command

import (
	"fmt"
	"github.com/guntutur/todo-cli/datastore"
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{Use: "create", Short: "create todo with task_content", Run: CreateTodo}

func init() {
	createCommand.Flags().String("task_content", "", "create todo with task_content")
	createCommand.MarkFlagRequired("task_content")
	ContextCommand.AddCommand(createCommand)
}

func CreateTodo(cmd *cobra.Command, args []string) {
	taskContent := cmd.Flag("task_content").Value.String()
	fmt.Println("created todo : " + taskContent)
	//datastore.CreateTodo(taskContent)
	datastore.CreateTodoCDB(taskContent)
}
