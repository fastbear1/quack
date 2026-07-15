# quack
Tool for auto creating goose migration files according to gorm models and database schema. This project was inspired by python alembic tool.

## Project status
Project status and current restirction
 - status is beta
 - only PostgreSQl database driver
 - different table definition restriction(details [Restrictions](#restrictions))

## Installing locally
```bash
    git clone github.com/fastbear1/quack
    cd ./quack
    go build cmd/main.go
```

## How to use
Getting help information
```bash
    quack --help
```
```bash
Uasge:
  quack [flags] command
Commands:
  - run - quack(run and create) goose migration file
  - help - show help information
  - version - show current version(can be used for checking config file)
  command usage example:
    - quack help - show help information
    - quack run - run creating a new migration files
flags:
  - models - directory where all gorm struct models live
  - uri - URI connection database string
  - dbname - database name
  - path - directory where a new migration files will be stored
  - exclude - exclude gorm struct models(usually it's a embeded struct or not a model struct)
  - db-exclude - exclude database tables(for example goose_migrations table)
  flag usage examples:
    - quack --models=models --path=migrations --uri=postgres://user:pass@host:port/database --exclude=Base,TestUsers run
    - quack --modesl=internal/models --path=models --uri=postgres://user:pass@host:port/database --exclude=Base --db-exclude=goose,goose_migrations run`
)
```

----------------
### Restrictions

#### Embed field
 - only anonymous declaration of embed struct

#### Relation fields
 - only explicit declaration for reference fields with tag contains foreignKey attribute

#### Alter column
 - only data, nullable and default value checking for altering column

#### parsing Index tags
 - only unique, index type, expression and column list are used for creaeting indices. Same params used to check index state(alter index)
