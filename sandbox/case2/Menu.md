####

1. Configure local env, see [Local Env](../init/LocalEnv.md)
2. Create init migrations from [Case 1](../case1/Menu.md)
3. Create migrations for managing columns
```bash
    ./quack --models=sandbox/case2/models --path=sandbox/case2/migrations run columns_migration
```
3. Apply migration to DB
```bash
   $HOME/go/bin/goose postgres "user=auqck dbname=quack password=pass host-postgres sslmode=disable" -dir=sandbox/case2/migrations up
```
