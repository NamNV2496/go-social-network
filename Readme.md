# This is a basic social network demo by Golang.

# 1. Technologies

    1. grpc - "google.golang.org/grpc"
    2. goqu - "github.com/doug-martin/goqu/v9"
    3. kafka - "github.com/IBM/sarama"
    4. redis - "github.com/go-redis/redis/v8"
    5. websocket - "github.com/gorilla/websocket"
    6. gorm - 	"gorm.io/gorm" (post-service only)

# 2. How to start

## Install

    go install github.com/google/wire/cmd/wire@latest


## Docker

    1. run `docker compose up`
    2. Check status of services

## Local

    1. go to "cd deployments" and run "docker-compose up"
    2. run http-gateway "cd http_gateway/" and "go run cmd/main.go"
    3. run user-service "cd user-service/" and "go run cmd/main.go"
    4. run post-service "cd post-service/" and "go run cmd/main.go"
    5. run newsfeed-service "cd newsfeed-service/" and "go run cmd/main.go"
    6. run message-service "cd message-service/" and "go run cmd/main.go"
    7. run FE "cd web/" and "start index.html"

![flow](docs/flow.png)

## User

    namnv - namnv
    knm - knm
    baobq - baobq

![alt text](docs/login.png)

## Follower

![alt text](docs/follower.png)

# 3. Data migration

    user-service\internal\database\migrations\mysql

# 4. Demo

## Newsfeed

![alt text](docs/newsfeed.png)

## View comment

![alt text](docs/viewComment.png)
## View followers post

if `namnv` posted a post. `knm` and `baobq` will see it in their newsfeed.

![alt text](docs/viewPost.png)

if `knm` posted a post. Only `namnv` will see it in his newsfeed `baobq` will not. Because only `namnv` follow `knm`
if `baobq` posted a post. Only `namnv` will see it in his newsfeed `knm` will not. Because only `namnv` follow `baobq`

![alt text](docs/viewPost1.png)

Here is redis
![alt text](docs/viewPost2.png)

## popup in newsfeed

Create new chat room private with friend from newsfeed. Only 2 member of chat can see and join the chat

![alt text](docs/popup.png)


## Your wall

![alt text](docs/wall.png)

![alt text](docs/wall1.png)

## Search another people

![alt text](docs/search.png)

![alt text](docs/search1.png)


