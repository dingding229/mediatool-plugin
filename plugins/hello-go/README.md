# hello-go

A built-in lifecycle example, not an installable executable. The registry exposes
only metadata/factory references. After the Host authorizes initialization, the
factory receives this instance's typed Host client. `InitPlugin` reads and
CAS-updates the `starts` key; form/page output is synthetic and carries no secrets.

The schema/golden and real stdio test child use only mock state, not a database,
Redis, production account or media. `StopService` is idempotent. The application
serializes plugin callbacks. A failed CAS is not blindly retried.
