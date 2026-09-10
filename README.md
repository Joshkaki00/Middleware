# Middleware

Minimal Echo (v5) server demonstrating custom middleware: request tracing, header injection, and a spoof-resistant "inside the building" IP check.

## Install

```bash
go get github.com/labstack/echo/v5@v5.3.1
```

## Usage

```bash
go run .
```

The server starts on `:1323`.

```bash
curl -i http://127.0.0.1:1323/
curl -s http://127.0.0.1:1323/whoami
curl -i http://127.0.0.1:1323/internal/admin
```

Expected output for `/whoami`:

```json
{"ip":"127.0.0.1","inside":true,"via":"realip+private/loopback"}
```

## How it works

Middleware runs in this order: `Pre` (before route matching) → `Use` chain → handler.

| Middleware | Lane | Purpose |
| --- | --- | --- |
| `PreTrace` | `Pre` | Tags every request, including unmatched routes, before the router decides on a match |
| `middleware.Recover()` | `Use` | Converts panics into `500` responses instead of crashing the process |
| `middleware.RequestID()` | `Use` | Adds an `X-Request-Id` header for correlating logs |
| `ServerHeader` | `Use` | Sets a custom `Server` response header |
| `InsideTheBuilding` | `Use` | Labels the request `inside`/`outside` using `c.RealIP()` and private/loopback ranges |
| `EnforceInside` | `Use` (scoped to `/internal`) | Returns `403` for any request not labeled `inside` |

`e.IPExtractor = echo.ExtractIPDirect()` means `c.RealIP()` trusts only the raw TCP peer address, not client-supplied headers like `X-Forwarded-For` or `X-Real-IP`. This is why forged headers do not change the `inside` result — the extractor never reads them.

## Routes

| Route | Description |
| --- | --- |
| `GET /` | Returns `OK` |
| `GET /whoami` | Returns `{"ip", "inside", "via"}` for the caller |
| `GET /internal/admin` | Returns `200` only if the caller is labeled `inside`; otherwise `403` |

## Requirements

- Go 1.25 or newer (Echo v5 requirement)

## License

MIT
