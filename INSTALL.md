# Para desarrollo

```
cp .env.sample .env
go mod tidy

air -c .air.toml
```

# Para producción

```
go build -o restapi .

./restapi
```
