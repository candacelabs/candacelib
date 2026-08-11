# candacelib

`candacelib` is Candace's extended standard library: small, focused Go
packages with stable semantics and no dependency on the Candace server
monorepo. It requires Go 1.26.

## Cron

[`cron`](./cron) is a durable in-process scheduler with functional options, a
human-readable schedule DSL, Gin status-route integration, and an optional
relational PostgreSQL store generated with SQLC. Its versioned Liquid Proto
messages are boundary contracts only; the database remains ordinary typed
relational state.

```bash
go get github.com/candacelabs/candacelib/cron
```

```go
import (
	cron "github.com/candacelabs/candacelib/cron"
	cronpostgres "github.com/candacelabs/candacelib/cron/postgres"
)
```

See the [cron package README](./cron/README.md) for lifecycle, persistence,
catch-up, overlap, and idempotency guidance.

## Liquid Proto

Liquid Proto's [`liquidproto` runtime](./liquidproto),
[`liquidproto/v1` schema](./liquidproto/v1), and
[`protoc-gen-liquidproto`](./cmd/protoc-gen-liquidproto) generator are all part
of this module. Together they provide deterministic protobuf encoding,
validation at serialization boundaries, inspectable refinement errors, and
redacted error formatting.

```bash
go get github.com/candacelabs/candacelib/liquidproto
```

```go
import "github.com/candacelabs/candacelib/liquidproto"
```

Install the generator with:

```bash
go install github.com/candacelabs/candacelib/cmd/protoc-gen-liquidproto@latest
```

Schemas can import `liquidproto/v1/refinement.proto` and annotate singular
scalar fields with `(candace.liquid.v1.field)`. The generator emits a
`Validate<Message>` function for every message with annotated fields.

## License

Licensed under the [Apache License, Version 2.0](./LICENSE).
