## Introduction

This is a simple backend API that manages a merchant account.

# Project Setup

### Docker

* Install [Docker for Mac](https://docs.docker.com/docker-for-mac/install/)

* Setup and boot the Docker containers:

```
docker compose up -d
```

## Interaction

- Seed a admin user 
- Use the api_secret of admin user to create a merchant.
- Merchants can then manage their members using their own api keys
- Refer to [OpenAPI documentation](openapi.yaml)

## Todos
- Increase test coverage
- Add database seeds