# Book API — Go and Redis

A CRUD learning project from a **QLUB knowledge-transfer session**, built with Go, Gorilla Mux, and Redis. It explores HTTP handlers, JSON models, Redis access, background cleanup, and an external ISBN lookup.

## Run locally

Use Go 1.23.5 or a compatible newer toolchain, plus Redis at `localhost:6379` with no password and database 0. The Redis address is configured in [pkg/config/config.go](pkg/config/config.go).

```bash
git clone https://github.com/paranietharan/bookapi.git
cd bookapi
go mod download
go run ./cmd
```

Start Redis before the API. The HTTP server listens on `http://localhost:8080`; keep its terminal open, since the program waits for Enter to exit.

## API routes

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/books` | List books |
| POST | `/books` | Create a book |
| GET | `/books/{id}` | Retrieve a book |
| PUT | `/books/{id}` | Update a book |
| DELETE | `/books/{id}` | Delete a book |
| GET | `/book/isbn/{id}` | Fetch details from Open Library and add a book |

The ISBN route has a write side effect despite using GET and needs internet access. Request fields are defined in [bookDto.go](pkg/models/bookDto.go) and [book.go](pkg/models/book.go).

```bash
curl http://localhost:8080/books
curl -X POST http://localhost:8080/books   -H 'Content-Type: application/json'   -d '{"name":"Demo Book","author":"Demo Author","category":"Learning","description":"Example record"}'
```

## Implementation notes

- [router.go](pkg/router/router.go) registers endpoints and starts the server.
- [handler.go](pkg/handler/handler.go) handles requests.
- [store.go](pkg/store/store.go) combines Redis records with a process-local list; list state is not restored from Redis on restart.
- The cleanup loop is scheduled every 40 seconds. The current removal condition uses a timestamp age greater than 20 seconds, despite comments describing a longer period. Use disposable data when exploring this behavior.

## Development checks

```bash
go build ./...
go test ./...
```

No `*_test.go` files are currently included. Persistence consistency, timestamp updates, and concurrent access would need further work before using this as a durable application.
