#### Prepare local environment

##### Install goose
Check installation steps [Goose install](https://github.com/pressly/goose#install)

##### Build quack locally
```bash
    make build
```

##### Run local db
Default connection URI - postgres://auqck:pass@postgres:5432/quack

```bash
    docker-compose -f misc/dokcer-compose.yaml up -d
```

##### Create config file
Create quack_config.yaml file in project directory
```yaml
    version: 0.1
    database:
      uri: "postgres://quack:pass@postgres:5432/quack"
      name: "quack"
    exclude:
      - "goose_migrations"
    models:
      path: "misc/models"
    exclude:
    migrations:
      path: "misc/migrations"
```

##### Create you first migration file
```bash
    quack -models=sandbox/case1/models -path=sandbox/case1/migrations run init_migration
```

##### Apply migration files
```bash
    $HOME/go/bin/goose postgres "user=quack dbname=quack password=pass host=postgres sslmode=disable" -dir=misc/migrations up
```
