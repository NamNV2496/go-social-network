# This is a basic social network demo by Golang.

# 1. Technologies

    1. grpc - "google.golang.org/grpc"
    2. goqu - "github.com/doug-martin/goqu/v9"
    3. kafka - "github.com/IBM/sarama"
    4. redis - "github.com/go-redis/redis/v8"
    5. websocket - "github.com/gorilla/websocket"

# 2. How to start

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

![alt text](docs/image4.png)

Follower: 

![alt text](docs/image-2.png)

Check following from current_user

![alt text](docs/image7.png)
![alt text](docs/image8.png)

if `namnv` posted a post. `knm` and `baobq` will see it in their newsfeed.

![alt text](docs/image-1.png)

if `knm` posted a post. Only `namnv` will see it in his newsfeed `baobq` will not. Because only `namnv` follow `knm`
if `baobq` posted a post. Only `namnv` will see it in his newsfeed `knm` will not. Because only `namnv` follow `baobq`

![alt text](docs/image.png)

Here is redis
![alt text](docs/image3.png)

Create new chat room private with friend from newsfeed
![alt text](docs/image5.png)

Only 2 member of chat can see and join the chat
![alt text](docs/image6.png)

