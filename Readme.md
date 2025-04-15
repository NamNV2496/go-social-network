# This is a basic social network demo by Golang.

# 1. Technologies

    1. grpc - "google.golang.org/grpc"
    2. goqu - "github.com/doug-martin/goqu/v9"
    3. kafka - "github.com/IBM/sarama"
    4. redis - "github.com/go-redis/redis/v8"
    5. websocket - "github.com/gorilla/websocket"
    6. gorm - "gorm.io/gorm" (post-service only)
    7. Elastic search - "github.com/elastic/go-elasticsearch/v8"
    8. prometheus - "github.com/prometheus/client_golang/prometheus"
    9. grafana for monitoring
    10. email - "github.com/jordan-wright/email"
    11. sms - "github.com/twilio/twilio-go"

# 2. Prerequisite

## How to build protobuf

```
    cd http_gateway
    make generate

then copy to target service
```



access each service
```
    make generate
```

## Grafana custom

after login, you can setup dashboard for monitoring

![grafana1](docs/grafana1.png)

![grafana2](docs/grafana2.png)

![grafana3](docs/grafana3.png)

![grafana4](docs/grafana4.png)

![grafana5](docs/grafana5.png)

access to download dashboard template

https://grafana.com/grafana/dashboards/1860-node-exporter-full/

![grafana6](docs/grafana6.png)

import template by id or json

![grafana7](docs/grafana7.png)

![grafana8](docs/grafana8.png)

![grafana9](docs/grafana9.png)




# 3. How to start

## Run by Docker

    1. run `docker compose up`
    2. Check status of services
    3. Access prometheus: [http://localhost:9090/targets](http://localhost:9090/targets)
    4. Access grafana: [http://localhost:3000/login](http://localhost:3000/login) (user: admin - password: admin)
    5. Access `localhost:3001`


## Run on Local

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

# 4. Data migration

sql files location:

    user-service\internal\database\migrations\mysql

# 5. Demo

## Install wire to generate dependency injection

    go install github.com/google/wire/cmd/wire@latest

## Login

![alt text](docs/otp.png)
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


# Monitoring

![monitoring1](monitoring1.png)