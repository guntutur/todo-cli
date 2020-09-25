package datastore

import (
	"crypto/tls"
	"log"
	"os"
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

// Todo holds todo information
type Todo struct {
	ID          string
	Created     time.Time
	TaskContent string
	Tags        string
}
