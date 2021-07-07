# netdata digest service

## Requirements
### Dockerized
- [Docker](https://docs.docker.com/get-docker/) - 20.10.x
- [Docker-compose](https://docs.docker.com/compose/install/) - 1.24.x
### Local
- [Go](https://golang.org/doc/install) - 1.16.x
- [Redis](https://redis.io/download) - 6.2.x
- [Nats](https://docs.nats.io/nats-server/installation)

## Basic Usage (dockerized)

```sh
# start services
make
```
```sh
# run unit tests
make test
```
```sh
# run simple verification
# terminal 1 - start services
make
# terminal 2 - run verification script
make verify
```
```sh
# run end to end tests
# terminal 1 - start services scaling netdata-digest service to 3 nodes
make e2e.serve
# terminal 2 - run end to end tests
make e2e.test
```

## Diagram
![diagram](./docs/diagram.svg)

## Architectural decisions

### Storage agnostic


### Message broken agnostic

### No in-order delivery support

### Horizontal scale support

### Alarm events persistence

## Refrences
- [golang project layout standards](https://github.com/golang-standards/project-layout)
