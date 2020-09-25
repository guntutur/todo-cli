package datastore

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const redisHostEnvVar = "REDIS_HOST"
const redisPasswordEnvVar = "REDIS_PASSWORD"
const redisPasswordRequiredEnvVar = "REDIS_PASSWORD_REQUIRED"
const sslRequiredEnvVar = "REDIS_SSL_REQUIRED"
const defaultRedisHost = "localhost:6379"

var redisHost string
var redisPassword string

// true or false
var redisPasswordRequired = "true"

// true or false
var sslRequired = "true"

const todoIDCounter = "todoid"
const todoIDsSet = "todos_gentur-ids"
const defaultTag = "active"

func init() {
	redisHost = os.Getenv(redisHostEnvVar)
	if redisHost == "" {
		redisHost = defaultRedisHost
	}

	redisPasswordRequired = os.Getenv(redisPasswordRequiredEnvVar)
	if redisPasswordRequired == "" {
		redisPasswordRequired = "false"
	}

	if redisPasswordRequired == "true" {
		redisPassword = os.Getenv(redisPasswordEnvVar)
		if redisPassword == "" {
			log.Fatal("redis instance requires a password. please set in environment variable ", redisPasswordEnvVar)
		}
	}

	sslRequired = os.Getenv(sslRequiredEnvVar)
	if sslRequired == "" {
		sslRequired = "false"
	}
}

func getClient() *redis.Client {
	opts := &redis.Options{Addr: redisHost}

	if redisPasswordRequired == "true" {
		opts.Password = redisPassword
	}
	if sslRequired == "true" {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	c := redis.NewClient(opts)
	err := c.Ping().Err()
	if err != nil {
		log.Fatal("redis connect failed", err)
	}
	return c
}

func ListTodos(tags string) []Todo {
	c := getClient()
	defer c.Close()

	todoHashNames, err := c.SMembers(todoIDsSet).Result()
	if err != nil {
		log.Fatal("failed to get todo IDs", err)
	}

	todos := []Todo{}
	for _, todoHashName := range todoHashNames {
		id := strings.Split(todoHashName, ":")[1]

		todoMap, err := c.HGetAll(todoHashName).Result()
		if err != nil {
			log.Fatalf("failed to get todo from %s - %v\n", todoHashName, err)
		}

		var todo Todo
		if tags == "" {
			todo = Todo{id, todoMap["created"], todoMap["task_content"], todoMap["tags"]}
			todos = append(todos, todo)
		} else {
			if tags == todoMap["tags"] {
				todo = Todo{id, todoMap["created"], todoMap["task_content"], todoMap["tags"]}
				todos = append(todos, todo)
			}
		}
	}
	if len(todos) == 0 {
		fmt.Println("no todos founds")
		return nil
	}
	return todos
}

func CreateTodo(taskContent string) {
	c := getClient()
	defer c.Close()

	// set increment
	id, err := c.Incr(todoIDCounter).Result()
	if err != nil {
		log.Fatal("failed to increment id!", err)
	}
	todoid := "todo:" + strconv.Itoa(int(id))

	//store ID in a SET for other operations
	err = c.SAdd(todoIDsSet, todoid).Err()
	if err != nil {
		log.Fatal("failed to add todo id to SET", err)
	}

	//save todo in a HASH
	todo := map[string]interface{}{
		"created":      time.Now().String(),
		"task_content": taskContent,
		"tags":         defaultTag,
	}
	err = c.HMSet(todoid, todo).Err()
	if err != nil {
		log.Fatal("failed to save todo")
	}
	fmt.Println("todo saved! use './todo list' to show")
}

// Todo holds todo information
type Todo struct {
	ID          string
	Created     string
	TaskContent string
	Tags        string
}
