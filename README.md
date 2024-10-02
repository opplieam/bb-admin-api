# Buy-Better Admin API
![workflow](https://github.com/opplieam/bb-admin-api/actions/workflows/unittest.yml/badge.svg)
## Table of contents
- [Overview](#overview)
- [Project structure](#project-structure)
- [Dependencies](#dependencies)
- [Developer Setup](#developer-setup)
- [Running Test](#running-test)
- [Starting server and swagger openapi for Frontend development](#starting-server-and-swagger-openapi-for-frontend-development)
- [Running in local cluster minikube](#running-in-local-cluster-minikube)
- [Useful Command/Makefile](#useful-commandmakefile)
- [Database Schema](#database-schema)
- [Design choice](#design-choice)

## Overview
Buy Better Admin API is an admin backend to handle the admin task like `category matching`, `product matching`, 
`user managment` etc. including a helper for web scraping part.

Currently, only `category matching` is available.

`NOTE: This project is for learning purpose and not fully complete yet`

## Project structure
```
├── .gen                # auto generated by jet-db
│   ├── buy-better-admin
├── .github
│   ├── workflows
├── bin                 # go binary
├── cmd
│   ├── api             # main package for api server
│   ├── dbhelper        # a dev tool for database helper 
│   └── tokengen        # a dev tool for paseto generator
├── data                # the sql file using for seed or test db
├── internal
│   ├── middleware      # api middleware
│   ├── store           # database logic
│   ├── utils           # global utilities
│   └── v1
│       ├── category    # category handler
│       ├── probe       # liveness & readiness handler
│       └── user        # user handler
├── k8s
│   ├── base            # kustomize base
│   │   └── admin-api
│   ├── dev             # patch kustomize
│   │   └── admin-api
│   └── secret          # bitnami sealed secrets
│       └── dev
├── migrations          # generated by migrate tool
└── spec                # openapi v3 spec
```

## Dependencies
#### Infrastructure
- docker / docker-compose
- minikube
- kubectl / kustomize
#### Database tools
- CLI go [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- CLI [jet-db](https://github.com/go-jet/jet?tab=readme-ov-file#installation)
#### Testing tools
- CLI [mockery](https://vektra.github.io/mockery/latest/installation/)

## Developer Setup
1. Create environment variable store in `.env` file at root directory
```
WEB_SERVICE_ENV="dev"
WEB_ADDR=":3000"
WEB_READ_TIMEOUT=5
WEB_WRITE_TIMEOUT=40
WEB_IDLE_TIMEOUT=120
WEB_SHUTDOWN_TIMEOUT=20

DB_DRIVER="postgres"
DB_DSN="postgresql://postgres:admin1234@localhost:5432/buy-better-admin?sslmode=disable"
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME="15m"

TOKEN_ENCODED="1c0021bc344fa16c72fc522c53bfe9f77a2a597507374e56e3a275759c4c1562"
```
> For `TOKEN_ENCODED`, you can random generate using [this](https://www.browserling.com/tools/random-hex) and use 64 digits

2. Create postgres environment variable in `.postgres.env` file at root directory. This will be used by `docker-compose`. 

```
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="admin1234"
POSTGRES_DB="buy-better-admin"
```
> postgres environment variables must be match with Makefile

3. Visit `Makefile` There are 4 important variables for local development. Feel free to edit. 
```
DB_DSN
DB_NAME
DB_USERNAME
CONTAINER_NAME		
```

4. Start the development postgres db `make dev-db-up` this command does follow 
   * docker-compose with postgres image
   * sleep for 3 seconds
   * migrate up
   * seed the fake data with `dbhelper` tool

5. `go run cmd/api` start the server with port `:3000`

## Running Test
- Run only unit test `make test-unit-v`
- Run only integration test `make test-integr-v`
- Run all test `make test-all-v`

## Starting server and swagger openapi for Frontend development
![swagger](https://github.com/opplieam/buy-better/blob/main/swagger.png?raw=true)
- `make server up` starting postgres db, swagger and buy better admin server
  * `localhost:3000` - buy better admin server
  * `localhost:8081` - swagger openapi
- `make server down` shutting postgres db and swagger down

## Running in local cluster minikube
1.  create `encoded-secret.yaml` under k8s/secret/dev
```
apiVersion: v1
kind: Secret
metadata:
  name: admin-api-secret
  namespace: buy-better

type: Opaque
data:
  TOKEN_ENCODED: "your base64 encode to TOKEN_ENCODED"
```
2. `make dev-up-all`
    * starting postgres db -> migration -> seed fake data
    * starting minikube 
    * apply `bitnami-sealed-secrets` controller

3. `minikbue tunnel` to expose load balancer
4. `make dev-apply`
    * go mod tidy
    * building an image with docker
    * kustomize apply resources
    * generate and apply bitnami secret
    * restart deployment (due to bitnami seal secret controller changing certificate everytime when starting a new cluster)
     
## Useful Command/Makefile

Please visit `Makefile` for the full command.
- `make token-gen-build` build a binary of paseto token generator for testing
   * `make token-gen-valid` generate a valid token with 1 hour expiration and user_id 1
   * `make token-gen-expire` generate an expired token with user_id 1
- `make jet-gen` generate a type safe from database. run this command everytime there is a change in database schema.
- `mockery` generate a mock file. please visit `.mockery.yaml` for the setting
- `make dev-db-reset` restart the postgres container. run when you want to reset the database

## Database Schema

![db](https://github.com/opplieam/bb-admin-api/blob/main/Buy-Better-Admin.png?raw=true)

## Design choice

Most of the tools and 3rd party libraries are for learning purpose and convenience. I will try to explain some libraries.

- `go migrate` a go native migration tool with go SDK So it can run programmatic migration. 
It's useful when run integration test with `dockertest` 
- `gin` it has many useful features like `validator` and middleware. Yes, I can use `chi` router but eventually I will
use `validator` package. So I picked `gin` which already had built-in `validator`.
- `testify` a test suite feature is the reason. I planned to separate a unit and integration test with test suite.
- `jet-db` a type safe sql builder. very good dx for dynamic queries from my research.
which is the best match for this project. I also prefer SQL style rather than ORM
- `paseto` instead of jwt. This is related to the frontend development. I planned to store token for both localstorage
and cookies (token and refresh token). both storage have pros and cons. So I spread the risk into 2 storage 
(cookies and localstorage). Since paseto is an encrypted token. It made it very difficult to encrypted 
if the token is leaked.