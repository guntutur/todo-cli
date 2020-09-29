## Todo Cli Apps Build using golang and redis

go version go1.12.9 darwin/amd64
Redis server v=5.0.5 sha=00000000:0 malloc=libc bits=64 build=31cd6e21ec924b46

### Run Locally
This assume you have redis installed and exposed on its default `6379` port<br>
run command `go buid -o todo` in root dir<br/>
which will build binary executable called in root `todo`.
for detail and explanation of feature, please use command below

`./todo -h`

which would output something like this

````
Starting todo apps, using default DB : <db_name>

Usage:
  todo [command]

Available Commands:
  complete    complete todo
  create      create todo with task_content
  delete      delete todo, will not be shown afterwards, but are kept in memory datastore
  help        Help about any command
  list        list all todos

Flags:
  -h, --help      help for todo
  -v, --version   version for todo

Use "todo [command] --help" for more information about a command.
````

to change <db_name>, go to command/datastore/couchdb.go at line 52
the built in functionality of cli apps are made possible with help of [github.com/spf13/cobra](https://github.com/spf13/cobra)

### create task
`./todo create --task_content "your task"`

### list all task
`./todo list`

### list all task with input tags
currently there is only `active`, `completed`, and `deleted` tags known
`./todo list --tags <inpt_tags>`

### complete task
`./todo complete --id "task_id_integer"`

### delete task
deleted tags won't shown in list, but are kept in memory datastore
`./todo delete --id "task_id_integer"`

## RUn Cli Apps with docker
````
docker build -t todos_guntur_cli:dev .
````

then run
````
docker run -it todos_guntur_cli:dev sh
````
to sh into docker image, we will automatically entering /app where our todo binary resided, from there, basic curds operation can be done


## Run Cli Apps and redis with docker

To integrate our apps with redis container, run below command
````
docker-compose up -d --build
````

this will bring up both our cli apps and redis inside their own container<br/>
after the operation succeeded, run below command to sh into individual services

# cli apps
````
docker-compose exec app sh 
````
after successfully executing above command, we are in directory /app where our source code and binary executable of `todo` resided<br/>
try to launch the binary by executing `./todo list`<br/>
if no error in console, we are successfully integrating our cli apps inside container with redis on the another container, try a few cruds operation then moving on to next section

# redis datastore
````
docker-compose exec redis sh
````

## verify redis installation
using exec redis above, enter command `redis-cli`
after prompt open, type in `ping` and then enter to see `pong` in the output

verify that all data exists after basic cruds operation with `SCAN 0 COUNT 1000 MATCH todo*` in redis-cli console
show individual record with `HGETALL todo:1`

use below command to stop all running contaier
````
docker-compose down
````
