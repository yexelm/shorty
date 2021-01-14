[![Go Report Card](https://goreportcard.com/badge/github.com/yexelm/shorty)](https://goreportcard.com/report/github.com/yexelm/shorty)

# Description

Shorty is a practice project showcasing a simple link shortener built on top of Redis. It can be run either locally or
in Docker. If shorty is running in Docker, Redis persistently saves data into volume `redis-data` described in
docker-compose.yml

## Dependencies

- `go 1.15`
- `docker`

## API

```
POST /
```

Saves a full URL passed through POST request body into Redis, generates and returns a unique short alias for this URL

```
GET /<short_alias>
```

Retrieves original full URL saved into Redis earlier by its <short_alias>

## Environment variables

All env variables are set through .env file (example provided in the repository).

- `HOST_PORT` application port
- `CONTAINER_PORT` docker container port  
- `REDIS_URL` URL used for connection to Redis
- `DB_NUM` Redis db number for storing application data (original URL to short alias relation)

- `TEST_HANDLERS_DB` Redis db number used for testing handlers
- `TEST_STORAGE_DB` Redis db number used for testing storage

## Make commands

```
make test
```

Builds shorty via docker-compose.test.yml, launches tests, stops and removes containers.

```
make run
```

Builds shorty via docker-compose.yml and runs it in Docker.

```
make down
```

Stops shorty running in Docker and removes all related containers.

```
make clear
```

Same as `make down` but also removes `redis-data` volume where application data is stored
