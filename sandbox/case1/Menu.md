####

1. Configure local env, see [Local Env](../init/LocalEnv.md)
2. Create init migrations
```bash
    ./quack --models=sandbox/case1/models --path=sandbox/case1/migrations run init_migration
```
3. Apply migration to DB
```bash
   $HOME/go/bin/goose postgres "user=auqck dbname=quack password=pass host-postgres sslmode=disable" -dir=sandbox/case1/migrations up
```
