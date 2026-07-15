[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/)
[![Status](https://img.shields.io/badge/beta-0.42.1-blue)](https://github.com/fastbear1/quack/releases)
[![Tests](https://github.com/fastbear1/quack/actions/workflows/citest.yml/badge.svg)](https://github.com/fastbear1/quack/actions/workflows/citest.yml)

# quack
Tool for auto creating goose migration files according to gorm models and database schema. This project was inspired by python alembic tool.

# DISCLAIMER: 
**Always check migration file before applaying it even if everything looks fine**

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
### Config files
Config file description
```yaml
version: 1                                            //version for local commands for future use
database:
  uri: "postgres://user:passwor@postgres:port/database" // database connetion URI
  name: "database"                                      // name for database 
  exclude:                                              // exclude tables from parsing
    - "goose_migrations"
    - "other_table"
models: 
  path: "models"                                        // path(directory name) with gorm models
  exclude:                                              // exclude models from parsing
    - "Base"
    - "Model"
migrations:                                             // path where to store migration files
  path: "migrations"
```
Excluded gorm models are usually embed strcut which used in many tables.

Example:
```go
type Base struct {
	ID uuid.UUID `gorm:"index;type:uuid;primary_key;default:gen_random_uuid()"`
	// Test comment
	CreatedAt time.Time `gorm:"type:timestamp without time zone;not null;default:now();<-:create"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone;not null;default:now()"`
}
```
Excluded database tables are independent tables or misc tables used by another client. 

Example:

 - goose_migrations
 - cache (table for caching routines)

### Usage
Getting help information
```bash
    quack --help
```
```bash
Uasge:
  quack [flags] command filename
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
    - quack --models=models --path=migrations --uri=postgres://user:pass@host:port/database --exclude=Base,TestUsers run MIGRATIONFILE_NAME
    - quack --modesl=internal/models --path=models --uri=postgres://user:pass@host:port/database --exclude=Base --db-exclude=goose,goose_migrations run MIGRATIONFILE_NAME 
)
```

Command flags has priority on config file parameters. For different example of usage see [Playground](./playground/Menu.md) cases. 
After file was created and checked use goose migration tool for applying newest migrations to database.

----------------
### Restrictions

#### Embed field
 - only anonymous declaration of embed struct

    Correct:
    ```go
        type User struct {
        	Base
        	Name   string `gorm:"not null"`       // full name field
        	Email  string `gorm:"not null"`       // User email field
        }
    ```
    Incorrect
    ```go
        type Blog struct {
          ID      int
          Author  Author `gorm:"embedded;embeddedPrefix:author_"`
          Upvotes int32
        }
    ```

#### Relation fields
 - only explicit declaration for reference fields with tag contains foreignKey attribute

    Correct:
    ```go
        type User struct {
        	gorm.Model
            CreditCardID int `gorm:type:smallint`
            CreditCard CreditCard `gorm:"foreignKey:CreditCardID;referenceName:user_creditcard_credit_card_id_id"`
        }
    ```
    Incorrect
    ```go
        type User struct {
            gorm.Model
            CreditCard CreditCard
        }
    ```


#### Alter column
 - only data, nullable and default value checking for altering column

#### Parsing Index tags
 - only unique, index type, expression and column list are used for creaeting indices. Same params used to check index state(alter index)

#### Constaraint and indices names
 - preffered way is to declare names for constraints and indices

   Example:
   ```go
       type Command struct {
            Base
	        Name  string    `gorm:"type:varchar(255);not null"`
	        Cid   uuid.UUID `gorm:"type:uuid;index;indexName:commands_cars_cid__id;default:uuidv4()"`
	        OwnerId uuid.UUID `gorm:"type:uuid;not null"`
	        Owners  Owners    `gorm:"foreignKey:OwnerId;referenceName:commands_owner_owner_id_id;constraint:OnDelete:CASCADE;"`
        } 
   ```

