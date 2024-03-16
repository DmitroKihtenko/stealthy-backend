# Stealthy Backend
Backend service component of Stealthy application. REST API web server.
Encapsulates user's service business logic of Stealthy application:
 - provides service of sign up and sign in process for users
 - provides service of retrieving data of authenticated user
 - provides service of uploading, downloading user's files and viewing
details about uploaded files

### Technologies
Build with:
 - Golang 1.21.1
 - Gin framework 1.9.1
 - Go mongo 1.13.0

Works on HTTP web protocol.

### Requirements
Installed Docker and Docker-compose plugin

### API
You can view API details using Openapi standard mapping `/swagger/index.html`

### How to up and run
Configure application
1. Copy files: `.env.example` to `.env`, `config.yaml.example` to
`config.yaml`.
2. Make changes you need in configuration files (details about configs can
be found in `config.yaml.example` and `.env.example` files)

Build docker images
Build docker images and start service
```bash
docker compose up
```

Stop and remove containers after application use
```bash
docker compose down
```

### How to run application tests
```shell
docker compose -f docker-compose-test.yml build test && \
docker compose -f docker-compose-test.yml run --rm test
```

### How to check and format code
Check code style:
```shell
docker compose -f docker-compose-test.yml build check-code && \
docker compose -f docker-compose-test.yml run --rm check-code
```

Format code:
```shell
docker compose -f docker-compose-test.yml build check-style && \
docker compose -f docker-compose-test.yml run --rm check-style gofmt -l ./
```
