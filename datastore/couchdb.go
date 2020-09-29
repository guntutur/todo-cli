package datastore

import (
	"context"
	"crypto/rand"
	"fmt"
	_ "github.com/go-kivik/couchdb"
	"github.com/go-kivik/kivik"
	"log"
	"strings"
	"time"
)

var couchDBHost string
var DbName string

type TodoCDB struct {
	Key         int64  `json:"key"`
	Created     string `json:"created_date"`
	TaskContent string `json:"task_content"`
	Tags        string `json:"tags"`
}

type TodoCDBObj struct {
	ID        string    `json:"_id"`
	Rev       string    `json:"_rev,omitempty"`
	Todo      TodoCDB   `json:"todo"`
	DeletedAt string    `json:"deletedAt"`
	DeletedBy DeletedBy `json:"deletedBy"`
}

type DeletedBy struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func createId() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return strings.ToLower(uuid)
}

func init() {
	DbName = "todos_guntutur"
	couchDBHost = "http://admin:iniadmin@13.250.43.79:5984"
}

func dialCouchDBClient() *kivik.DB {
	client, err := kivik.New("couch", couchDBHost)
	if err != nil {
		log.Fatal("couchdb connect remote endpoint failed", err)
	}

	e, err := client.DBExists(context.TODO(), DbName)
	if !e {
		log.Printf("Default db %s not exists, creating..", DbName)
		err = client.CreateDB(context.TODO(), DbName)
		if err != nil {
			log.Fatalf("couchdb create db failed : %s", err)
		}
	}

	return client.DB(context.TODO(), DbName)
}

func ListTodosCDB(tags string) []TodoCDBObj {
	c := dialCouchDBClient()

	var todos []TodoCDBObj
	rows, err := c.AllDocs(context.TODO(), kivik.Options{"include_docs": "true"})
	if err != nil {
		log.Fatal("failed to fetch all docs", err)
	}

	for rows.Next() {
		var todoCdbObj TodoCDBObj
		if err := rows.ScanDoc(&todoCdbObj); err != nil {
			panic(err)
		}

		if tags == "" {
			if todoCdbObj.DeletedAt == "" {
				todos = append(todos, todoCdbObj)
			}
		} else {
			if tags == todoCdbObj.Todo.Tags {
				todos = append(todos, todoCdbObj)
			}
		}
	}
	if len(todos) == 0 {
		fmt.Println("no todos founds")
		return nil
	}
	c.Close(context.TODO())
	return todos
}

func CompleteTodoCDB(id string) {
	c := dialCouchDBClient()
	defer c.Close(context.TODO())

	row := c.Get(context.TODO(), id)
	var todoCdbObj TodoCDBObj
	if err := row.ScanDoc(&todoCdbObj); err != nil {
		log.Println("Todo not found")
	}

	todoCdbObj.Todo.Tags = "completed"
	_, err := c.Put(context.TODO(), id, todoCdbObj)
	if err != nil {
		log.Fatal("failed to complete todo id", id)
	}

	log.Printf("todo id %s completed. use './todo list' to show\n", id)
}

func CreateTodoCDB(taskContent string) {
	c := dialCouchDBClient()
	defer c.Close(context.TODO())

	t := time.Now()
	layout := "01/02/2006, 03:04:05 PM"
	baseTodo := TodoCDB{
		Key:         t.UnixNano() / 1e6,
		Created:     t.Format(layout),
		TaskContent: taskContent,
		Tags:        defaultTag,
	}

	todo := map[string]interface{}{"todo": baseTodo}
	rev, err := c.Put(context.TODO(), createId(), todo)
	if err != nil {
		panic(err)
	}

	log.Printf("todo with rev %s saved! use './todo list' to show", rev)
}

// actually update the tags to delete
func DeleteTodoCDB(id string) {
	c := dialCouchDBClient()
	defer c.Close(context.TODO())

	row := c.Get(context.TODO(), id)
	var todoCdbObj TodoCDBObj
	if err := row.ScanDoc(&todoCdbObj); err != nil {
		log.Println("Todo not found")
	}

	todoCdbObj.DeletedAt = time.Now().String()
	tmpId := strings.Split(DbName, "_")[1]
	todoCdbObj.DeletedBy = DeletedBy{
		ID:    tmpId,
		Email: tmpId + "@gmail.com",
	}
	_, err := c.Put(context.TODO(), id, todoCdbObj)
	if err != nil {
		log.Fatal("failed to delete todo id", id)
	}

	log.Printf("todo id %s deleted. use './todo list' to show\n", id)
}
