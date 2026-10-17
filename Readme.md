# This is a basic social network demo by Golang.

# 1. Technologies

    1. grpc - "google.golang.org/grpc"
    2. goqu - "github.com/doug-martin/goqu/v9"
    3. gorm - "gorm.io/gorm" (post-service only)
    4. kafka - "github.com/IBM/sarama"
    5. redis - "github.com/go-redis/redis/v8"
    6. websocket - "github.com/gorilla/websocket"
    7. Elastic search - "github.com/elastic/go-elasticsearch/v8"
    8. prometheus - "github.com/prometheus/client_golang/prometheus"
    9. grafana for monitoring
    10. email - "github.com/jordan-wright/email"
    11. sms - "github.com/twilio/twilio-go"
    12. template - "html/template" for notification form
    13. DSA: Trie - for violent comment rule checking (post-service)
    13.1. DSA: Trie - for address auto suggestion (in user-service)
    14. Nginx: for load balancing
    15. Post validator - "github.com/go-playground/validator/v10"
    16. Regex - "regexp"
    17. Langfuse

![alt text](docs/flow.png)

# 2. How to start

## Run by Docker

    1. Copy `.env.langfuse.example` to `.env.langfuse` and replace its placeholder secrets.
    2. Run `docker compose up`
    3. Check status of services
    4. Access prometheus: [http://localhost:9090/targets](http://localhost:9090/targets) to check is your service up successfully
    5. Access grafana: [http://localhost:3000/login](http://localhost:3000/login) (user: admin - password: admin)
    6. Access `localhost:3001` for web
    7. Access `localhost:3001/comment_rule.html` for comment-rule
    8. Access `http://localhost:3000/` for langfuse

## Run on Local

    1. go to "cd deployments" and run "docker-compose up"
    2. run http-gateway "cd http_gateway/" and "go run cmd/main.go"
    3. run user-service "cd user-service/" and "go run cmd/main.go"
    4. run post-service "cd post-service/" and "go run cmd/main.go"
    5. run newsfeed-service "cd newsfeed-service/" and "go run cmd/main.go"
    6. run message-service "cd message-service/" and "go run cmd/main.go"
    7. run FE "cd web/" and "start index.html"

# 3. Information

### User

    namnv - namnv
    knm - knm
    baobq - baobq

### Follower

![alt text](docs/follower.png)

### Data migration

sql files location:

    user-service\internal\database\migrations\mysql

### Install wire to generate dependency injection

    go install github.com/google/wire/cmd/wire@latest

# 4. Features

```txt
- Login / Register
- Newsfeed
- View and add new comment on post with comment rule
- View followers post
- Follow / Unfollow
- View your wall and your friends post
- Search another people
- Chat with friends
- Email / SMS OTP
- Notification
```
## 4.1 Login / Register

- Login

![alt text](docs/login.png)

- Register

![alt text](docs/register.png)

- OTP

![alt text](docs/otp.png)

## 4.2 Newsfeed 

![alt text](docs/newsfeed.png)

## 4.3 View and add new comment on post with comment rule

- View comment on post

![alt text](docs/viewComment.png)

- Comment violent rule checking

![alt text](docs/comment_violent_rule.png)

## 4.4 View followers post

if `namnv` posted a post. `knm` and `baobq` will see it in their newsfeed.

![alt text](docs/viewPost.png)

if `knm` posted a post. Only `namnv` will see it in his newsfeed `baobq` will not. Because only `namnv` follow `knm`
if `baobq` posted a post. Only `namnv` will see it in his newsfeed `knm` will not. Because only `namnv` follow `baobq`

![alt text](docs/viewPost1.png)

## 4.5 Follow / Unfollow

![alt text](docs/popup.png)

## 4.6 View your wall and your friends post

- View your wall

![alt text](docs/wall1.png)

- View your friends post

![alt text](docs/wall.png)

## 4,.7 Search another people

- Search by name

![alt text](docs/search.png)

- Act when search

![alt text](docs/search1.png)

## 4.8 Chat with friends

- Chat room list

![alt text](docs/chat.png)

- Chat with friend

![alt text](docs/chat1.png)

## 4.8 Email / SMS OTP

## 4.9 Notification

## 4.10 Location searching

Support search by ward and district (will support city in the future)

![alt text](docs/search3.png)

![alt text](docs/search2.png)


<!-- ![alt text](docs/notification.png) -->

# Monitoring

![monitoring2](docs/monitoring2.png)
![monitoring3](docs/monitoring3.png)
![monitoring4](docs/monitoring4.png)

prometheus read this [link](./user-service/readme.md)

<details>
  <summary>Grafana custom</summary>
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
</details>

# What's next?

- add CDC for new post
- map reducing for postdata

## Langfuse API keys (Public / Secret)

If you run Langfuse locally (docker-compose) you'll need a project Public Key (PK) and Secret Key (SK) so the services can upload traces.

Quick steps to get PK/SK:

1. Open the Langfuse UI in your browser: http://localhost:3000
2. Sign in (or create the initial admin user if you haven't).
3. Go to your **Project** (top-left), then **Settings** → **API keys** (or **Project Settings** → **API Keys**).
4. Create a new API key — the UI will show the **Public Key** and **Secret Key**. Copy both immediately (the secret is shown once).

Where to set them

- In the service config file: set under `comment_rule.langfuse` and `comment_rule.langfuse_trace` in `post-service/configs/local.yaml`.

Example `local.yaml` snippet:

```yaml
comment_rule:
    langfuse:
        enabled: true
        endpoint: http://localhost:3000
        public_key: pk-xxxxxxxx
        secret_key: sk-xxxxxxxx
    langfuse_trace:
        endpoint: http://localhost:3000
        public_key: pk-xxxxxxxx
        secret_key: sk-xxxxxxxx
        trace_name: comment_rule_ai_check
        verbose_dump: true
```

- Or provide via docker-compose/service env vars (example already present for Langfuse services):

```yaml
environment:
    - LANGFUSE_INIT_PROJECT_PUBLIC_KEY=pk-xxxxxxxx
    - LANGFUSE_INIT_PROJECT_SECRET_KEY=sk-xxxxxxxx
```

Note: keep the Secret Key private; do not check it into source control.
