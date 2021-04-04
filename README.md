[![Go Reference](https://pkg.go.dev/badge/github.com/yexelm/shorty.svg)](https://pkg.go.dev/github.com/yexelm/shorty)
[![Go Report Card](https://goreportcard.com/badge/github.com/yexelm/shorty)](https://goreportcard.com/report/github.com/yexelm/shorty)

# Description

Shorty is a practice project showcasing a simple link shortener built on top of Redis. 
It is intended to be run in Docker. Application data is stored in `/data` volume of Redis container.

## Dependencies

- `go 1.16`
- `docker`

## API

```
POST / -d '<original URL>'
```

Saves an original URL passed via POST request body into Redis, generates and returns a unique short alias for the original URL. If this URL is already present in Redis, simply returns an existing alias for it.

```
GET /<short_alias>
```

Retrieves original full URL saved into Redis earlier by its <short_alias>

## Example

```shell
make run
```
```shell
curl localhost:8080 -d 'google.com'
```
> localhost:8080/b
```shell
curl localhost:8080/b
```
> google.com
```shell
curl localhost:8080 -d 'golang.org'
```
> localhost:8080/c 
 ```shell
curl localhost:8080/c
```
> golang.org
```shell
curl localhost:8080 -d 'google.com'
```
> localhost:8080/b

## Environment variables

All env variables can be set in the .env file (example provided in the repository).

- `HOST_PORT` application port;
- `CONTAINER_PORT` docker container port;
- `REDIS_URL` URL used for connection to Redis;
- `DB_NUM` Redis db number where the data is stored.

## Make commands

```
make run
```

Builds shorty via docker-compose.yml and runs it in Docker.

```
make down
```

Stops shorty containers running in Docker and removes all related images.

```
make clear
```

Same as `make down` but also removes `redis-data` volume where the application data is stored.

## Metrics
Basic metrics are provided by Prometheus and available via `/metrics` handler on the `:8081` port.
