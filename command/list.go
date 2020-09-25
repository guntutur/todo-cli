package command

import (
	"github.com/guntutur/todo-cli/datastore"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var listCommand = &cobra.Command{Use: "list", Short: "list all todos", Run: ListTodo}

func init() {
	listCommand.Flags().String("tags", "", "active, completed or deleted")
	ContextCommand.AddCommand(listCommand)
}

func ListTodo(cmd *cobra.Command, args []string) {
	status := cmd.Flag("tags").Value.String()

	var todos []datastore.Todo
	if status == "completed" || status == "pending" || status == "in-progress" || status == "" {
		todos = datastore.ListTodos(status)
	} else {
		log.Fatalf("provide valid status - completed, pending or in-progress")
	}

	todoTable := [][]string{}

	for _, todo := range todos {
		todoTable = append(todoTable, []string{todo.ID, todo.TaskContent, todo.Tags, todo.Created})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Task Content", "Tags", "Created"})

	for _, v := range todoTable {
		table.Append(v)
	}
	table.Render()
}
