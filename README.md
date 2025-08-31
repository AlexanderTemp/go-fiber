# Fiber GO API/REST + Conexión Clickhouse 🚀

![Go](https://img.shields.io/badge/Go-1.25.0-blue?logo=go)
![Fiber](https://img.shields.io/badge/Fiber-2.52.9-purple?logo=go)
![ClickHouse](https://img.shields.io/badge/ClickHouse-2.40.1-orange?logo=clickhouse)
![Air](https://img.shields.io/badge/Air-1.62.0-lightgrey?logo=go)

## 📋 Requisitos

- Go 1.25.0
- Air v1.62.0 (desarrollo)

## 🔧 Instalación de Go

```bash

tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
```

🐟 Configuración para Fish Shell

```bash
set -gx GOROOT /usr/local/go
set -gx GOPATH "$HOME/go"
set -gx PATH $PATH /usr/local/go/bin $GOPATH/bin
```

## 🛠️ Desarrollo

```
cp .env.sample .env
go mod tidy
```

Se hace uso de air para live-loading.

```
go install github.com/air-verse/air@latest

air -c .air.toml
```

## 🚀 Producción

```
go build -o restapi .

./restapi
```

## Apuntes

- Fiber intenta no crear copias innecesarias en memoria.

- Para eso, muchos valores de \*fiber.Ctx (por ejemplo c.Params("foo"), c.Query("id"), c.Body(), etc.) son en realidad referencias a un buffer interno reutilizable.

- Ese buffer se recicla entre cada request para ahorrar memoria y velocidad.

```go
var saved string

func handler(c *fiber.Ctx) error {
    saved = c.Params("foo") // ❌ se va a sobrescribir en la próxima request
    return nil
}
```

### Patrón mental

- Dentro del handler → usa c.Params(), c.Query(), c.Body(), etc. sin miedo.

- Fuera del handler / en goroutines → haz copia (utils.CopyString, append([]byte(nil), val...), etc.).

---

### Estructura

```go
app.Method(path string, ...func(*fiber.Ctx) error)
```

- app is an instance of Fiber
- Method is an HTTP request method: GET, PUT, POST, etc.
- path is a virtual path on the server
- func(\*fiber.Ctx) error is a callback function containing the Context executed when the route is matched
