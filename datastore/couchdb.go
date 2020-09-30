package datastore

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	_ "github.com/go-kivik/couchdb"
	"github.com/go-kivik/kivik"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var couchDbUser string
var couchDbPwd string
var couchDbHost string
var couchDbPort string
var couchDbName string

const CouchDbUser = "COUCHDB_USER"
const CouchDbPwd = "COUCHDB_PWD"
const CouchDbHost = "COUCHDB_HOST"
const CouchDbPort = "COUCHDB_PORT"
const CouchDbName = "COUCHDB_NAME"

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

func constructCouchDbHost(user string, pwd string, host string, port string) string {
	return fmt.Sprintf("http://%s:%s@%s:%s", user, pwd, host, port)
}

// need a better way to set env variable from this point
// as this process is likely setting everything up every time each cruds command commence
// but at least it is working
func init() {

	_, userPresent := os.LookupEnv(CouchDbUser)
	_, pwdPresent := os.LookupEnv(CouchDbPwd)
	_, hostPresent := os.LookupEnv(CouchDbHost)
	_, portPresent := os.LookupEnv(CouchDbPort)
	_, dbNamePresent := os.LookupEnv(CouchDbName)

	if !userPresent && !pwdPresent && !hostPresent && !portPresent && !dbNamePresent {
		_, b, _, _ := runtime.Caller(0)
		d := path.Join(path.Dir(b))
		file := filepath.Dir(d) + "/setenv.sh"
		cmd := exec.Command("/bin/sh", "-c", "source "+file+" ; echo '<<<ENVIRONMENT>>>' ; env")
		bs, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}
		s := bufio.NewScanner(bytes.NewReader(bs))
		start := false
		for s.Scan() {
			if s.Text() == "<<<ENVIRONMENT>>>" {
				start = true
			} else if start {
				kv := strings.SplitN(s.Text(), "=", 2)
				if len(kv) == 2 {
					os.Setenv(kv[0], kv[1])
				}
			}
		}
		initEnvVar()
	}
}

func initEnvVar() {
	couchDbUser = os.Getenv(CouchDbUser)
	if couchDbUser == "" {
		log.Fatalf("%s env var is not set", CouchDbUser)
	}

	couchDbPwd = os.Getenv(CouchDbPwd)
	if couchDbPwd == "" {
		log.Fatalf("%s env var is not set", CouchDbPort)
	}

	couchDbHost = os.Getenv(CouchDbHost)
	if couchDbHost == "" {
		log.Fatalf("%s env var is not set", CouchDbHost)
	}

	couchDbPort = os.Getenv(CouchDbPort)
	if couchDbPort == "" {
		log.Fatalf("%s env var is not set", CouchDbPort)
	}

	couchDbName = os.Getenv(CouchDbName)
	if couchDbName == "" {
		log.Fatalf("%s env var is not set", CouchDbName)
	}

	couchDbHost = constructCouchDbHost(couchDbUser, couchDbPwd, couchDbHost, couchDbPort)
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

func dialCouchDBClient() *kivik.DB {
	client, err := kivik.New("couch", couchDbHost)
	if err != nil {
		log.Fatal("couchdb connect remote endpoint failed", err)
	}

	e, err := client.DBExists(context.TODO(), couchDbName)
	if !e {
		log.Printf("Default db %s not exists, creating..", couchDbName)
		err = client.CreateDB(context.TODO(), couchDbName)
		if err != nil {
			log.Fatalf("couchdb create db failed : %s", err)
		}
	}

	return client.DB(context.TODO(), couchDbName)
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
				if todoCdbObj.DeletedAt == "" {
					todos = append(todos, todoCdbObj)
				}
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
	tmpId := strings.Split(couchDbName, "_")[1]
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
