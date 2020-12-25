# Description

Shorty is a toy project showcasing a simple link shortener built on top of Redis. It can be run either locally or in
Docker. If shorty is running in Docker, Redis persistently saves data into volume `redis-data` described in
docker-compose.yml

## API

```
POST /
```

Save a full URL passed through POST request body in Redis and receive a unique short alias for it

```
GET /:short_alias
```

Retrieve original full URL saved into Redis earlier by its short_alias

## Env variables

- `REDIS_URL` URL used for connection to Redis
- `DB_NUM` Redis db number for storing application data
- `APP_PORT` application port
- `TEST_DB_1` Redis db number used for testing storage
- `TEST_DB_2` Redis db number used for testing handlers

## Make commands

```shell
make test
```

Builds shorty via docker-compose.test.yml, launches tests, stops and removes containers.

```shell
make run
```

Builds shorty via docker-compose.yml and runs it in Docker.

```shell
make down
```

Stops shorty running in Docker and removes all related containers.

```shell
make clear
```

Same as `make down` but also removes `redis-data` volume where application data is stored

