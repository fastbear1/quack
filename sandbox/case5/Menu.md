####

1. Configure local env, see [Local Env](../init/LocalEnv.md)
2. Create and run all previous migrations from [Case 1](../case1/Menu.md), [Case 2](../case2/Menu.md), [Case 3](../case2/Menu.md), [Case 4](../case4/Menu.md)
3. Create migrations for managing columns
```bash
    ./quack --models=sandbox/case5/models --path=sandbox/case5/migrations run columns_altering
```
3. Apply migration to DB
```bash
   $HOME/go/bin/goose postgres "user=quack dbname=quack password=pass host-postgres sslmode=disable" -dir=sandbox/case5/migrations up
```
