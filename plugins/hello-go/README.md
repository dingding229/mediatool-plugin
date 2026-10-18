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
