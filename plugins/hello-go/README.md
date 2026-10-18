# hello-go

A built-in lifecycle example, not an installable executable. The registry exposes
only metadata/factory references. After the Host authorizes initialization, the
factory receives this instance's typed Host client. `InitPlugin` reads and
CAS-updates the `starts` key; form/page output is synthetic and carries no secrets.

The schema/golden and real stdio test child use only mock state, not a database,
Redis, production account or media. `StopService` is idempotent. The application
serializes plugin callbacks. A failed CAS is not blindly retried.

### F4-A acquisition examples (1.4.0)

Subscribes to six committed account/metrics/subscription facts; deduplicates by
host event UUID in this instance's KV (`acquisition-<eventId>`). Payloads contain
only opaque IDs, revision, status and media identity, never credentials or URLs.
`resource.selection` returns the same IDs by default; `resource.download` allows
original options unchanged. Optional `resourceReverse`, `resourceReject` and
`resourceNoCleanup` demonstrate narrowing only. `resourceDelay` exercises bounded
host timeout/fallback. No hook invokes a downloader or fetches media. Installing
this version expands scopes and requires renewed administrator confirmation.

### F4-B organizing policies (1.5.0)

Subscribes to six terminal/verified organizing facts using the same idempotent
KV receipts. Five typed transfer/storage policies return original paths/handle
and `pass` by default. `organizeFolder`, `organizeReject`, `organizeDelay` are
optional test/demo settings: the host revalidates every result. These handlers
cannot access STRM bodies, resolver URLs, credentials or absolute roots.
