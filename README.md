# Scafold Fiber GO REST

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

## Estructura

```go
app.Method(path string, ...func(*fiber.Ctx) error)
```

- app is an instance of Fiber
- Method is an HTTP request method: GET, PUT, POST, etc.
- path is a virtual path on the server
- func(\*fiber.Ctx) error is a callback function containing the Context executed when the route is matched
