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
	tags := cmd.Flag("tags").Value.String()

	//var todos []datastore.TodoRedis
	var todos []datastore.TodoCDBObj
	if tags == "completed" || tags == "active" || tags == "" {
		//todos = datastore.ListTodos(tags)
		todos = datastore.ListTodosCDB(tags)
	} else if tags == "deleted" {
		log.Fatalf("todos marked deleted are saved in memory but won't be shown anywhere :)")
	} else {
		log.Fatalf("provide valid tags - active, completed, or deleted")
	}

	todoTable := [][]string{}

	for _, todo := range todos {
		todoTable = append(todoTable, []string{todo.ID, todo.Todo.TaskContent, todo.Todo.Tags, todo.Todo.Created})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Task Content", "Tags", "Created"})

	for _, v := range todoTable {
		table.Append(v)
	}
	table.Render()
}
