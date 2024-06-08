## Basic-CRUD

This is a simple REST API for managing employees using Go, Redis, and the Chi router. It supports basic CRUD operations.

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/08Ankit04/basic-crud.git
cd basic-crud

docker-compose up

```

## Running Tests

### Unit Tests

To run unit tests and see the coverage percentage, execute the following command:

```sh
cd basic-crud/app/server
go test -cover -v ./...
```
