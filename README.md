# Intellias Golang Bootcamp [2025]


## Local development infrastructure

You can run the development infrastructure by using Docker Compose tool 
(make sure that you have installed the latest version).

Run the command
```shell
  make infra-up
```
which will create the PostgreSQL database.

Use username `pguser` and password `pgpassword` to connect to `pgdb` database. 

All database schema migrations should be in-sync after the dev infra is up
(check the output from `intellias-golang-bootcamp-db-migrations` container)

We can run 
```shell
  make db-schema-migrations
```
during development to apply new migrations (which are stored in `./internal/database/migrations` director)

## Commands

### Build as CLI

```shell
  make cli-build
```

### Run a CLI

```shell
  make cli-run
```


**HINT**: Try to run
```shell
    make
```
command itself to see the lint of all supported commands.
