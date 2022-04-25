# basic-authentication
## To run The code follow the below step

```shell
docker-compose up -d
EXPORT CONNECTION_STRING=postgresql://local:local@localhost:5432/postgres?sslmode\=disable;
EXPORT MAX_CONNECTION=10;
EXPORT MAX_IDLE_CONNECTION=10;
EXPORT JWT_SECRET_KEY=mysecretkey
go run main.go
```

## Deployment strategy
We do have a docker file which we can deploy on any container registry and use that in any orchestration.
