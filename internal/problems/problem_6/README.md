# Problem 6

**Context:** You are on-call. The async scraping endpoint (`POST /scrape/async`) is slowly eating up memory on the servers until they OOM crash.

**Expected:** The endpoint should start a background job and wait for it, but gracefully handle client disconnects.

**Symptom:** The test simulates clients disconnecting early. Currently, this causes background goroutines to leak and hang forever.

---

## Explicación del leak

### El problema del canal `chan bool` sin buffer

El handler arranca el worker y espera el resultado con un canal sin buffer:

```go
done := make(chan bool)
go performHeavyScrape(req.URL, done)

select {
case <-done:
    // ...
case <-r.Context().Done():
    // el handler retorna, pero la goroutine sigue viva
}
```

Un canal sin buffer solo deja completar el envío (`done <- true`) si hay **alguien recibiendo en ese mismo instante**. El flujo problemático es:

1. El cliente se desconecta, así que el `select` gana por la rama `r.Context().Done()`.
2. El handler retorna y **ya nadie recibe** de `done`.
3. El worker despierta tras su `time.Sleep`, ejecuta `done <- true` y se queda **bloqueado para siempre** porque no hay receptor.

Cada goroutine colgada mantiene viva su stack y sus variables capturadas. Bajo tráfico con desconexiones/timeouts, las goroutines se acumulan sin límite hasta que el proceso se queda sin memoria (el OOM del enunciado).

El canal `bool` en sí no es el problema; el problema es **acoplar el envío con la existencia de un receptor**. Dos formas de romper ese acoplamiento:

- **Buffer:** `make(chan bool, 1)` hace que el envío siempre tenga hueco aunque nadie lea, así el worker termina y sale.
- **Envío no bloqueante:** `select { case done <- true: default: }` para que el worker nunca se quede esperando.

### El tema del `context`

El buffer arregla el leak, pero el worker **sigue haciendo todo el trabajo pesado** aunque el cliente ya no esté: gasta CPU, red y memoria en un resultado que nadie va a leer.

La solución "production-grade" es **propagar el `context` del request al worker** para que el trabajo se cancele en cuanto el cliente se va:

```go
func performHeavyScrape(ctx context.Context, url string, done chan<- bool) {
    select {
    case <-time.After(100 * time.Millisecond): // en real: la llamada de red usando ctx
        done <- true
    case <-ctx.Done():
        // cliente desconectado: abortamos sin escribir en el canal
        log.Printf("Scrape aborted for %s: %v", url, ctx.Err())
    }
}
```

Claves:

- `r.Context()` se cancela automáticamente cuando el cliente se desconecta o expira el timeout.
- Pasándolo al worker, este corta su trabajo en la rama `ctx.Done()` y **retorna** → la goroutine muere, sin leak.
- En un scraper real no se usa `time.After`; se pasa el `ctx` directamente a la operación bloqueante para que se aborte sola:

```go
reqHTTP, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
resp, err := http.DefaultClient.Do(reqHTTP)
```

Combinar **canal con buffer** (cinturón de seguridad ante la carrera "termina justo cuando el handler ya salió") con **cancelación por `context`** (cortar trabajo inútil) elimina el leak y deja de malgastar recursos.
