# Media Tool built-in Go plugin library

Single Go module: `github.com/dingding229/mediatool-plugin` (Go 1.26.9).

This is **not** a runtime marketplace, an executable download service or an upload
format. Media Tool pins a commit of this module and compiles its registry into the
app image. Merely listing a definition never constructs a plugin. Only an
administrator-selected and enabled instance is run by the application's
`plugin-worker` subprocess (application supervisor integration is G2).

## Layout

- `sdk/`: mandatory/optional interfaces, contracts, bounded stdio JSON-RPC,
  instance-local typed Host helpers, typed event registration, mock Host and golden tools.
- `plugins/<id>/`: reviewed implementation, manifest, icon and documentation;
  optional prebuilt `ui/dist` is consumed at application image build time.
- `registry/`: pure metadata and factory references, no initialization callbacks.

No package `init`, IO, process execution or goroutine creation at package scope.
The SDK and plugins must not import the main application's internals or persistence
clients. Public DTO/enum/schema definitions have one source here. `KnownScope`
means a known contract, **not** a grant or evidence that an application implements
it. Every Host call must be authorized by the Host against its connection's
instance, scope/selectors, caller and current grant/generation revision.

## Development and release

1. Add `plugins/<id>` and its metadata; register its factory. Raise the plugin
   version and record history when changing behavior. No secrets in metadata/UI.
2. Test with `sdk/plugintest`: deny-by-default mock Host, revision checks, error/delay
   injection, reverse calls, lifecycle golden. The hello example only updates its
   own KV. JSON component schemas/extension dispatch and business clients are
   extended in F2/F3; declared unimplemented behavior is not a successful no-op.
3. Run, serially:
   ```sh
   go vet ./...
   go test -p 1 -count=1 -v ./...
   go test -race -p 1 -count=3 ./...
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./...
   ```
   Use the fixed Go toolchain with `GOTOOLCHAIN=local`. CI executes a real Go test
   child over stdio; it is not a distributable plugin binary or proof of the app's
   Linux process isolation. No database/Redis/115/production account is used.
4. Push this repository using the project's approved Git credentials and wait for
   its CI. Then bump the module pin in `dingding229/mediatool` and pass that
   application's CI, Actions -> immutable Docker Hub image -> backed-up server
   release flow. Pushing this repo alone **does not update production**.

Runtime does not download/compile code. Do not introduce external repositories,
source/commit/version selectors, ZIP/binary upload or independent plugin rollback.
The application preserves configuration/KV/data when a definition disappears;
expanded permissions require reapproval. This repository carries no credentials
or other projects' authentication. Missing push access is a blocker, not a reason
to change the distribution model.

## Trust and current scope

Reviewed code is trusted, with defects isolated by workers rather than a malicious
code sandbox. Application G2 uses UID 65532, explicit secret-free environment/FDs,
parent dumpable protection, process-group/reaper and resource limits. Landlock is
best effort. Same-UID readable files remain a review-controlled risk; this SDK
implements **none** of those OS controls. No setuid/UID pool/seccomp or production
security changes are required. API/service/command/module/event Host integration
is F2/F4, typed business clients F3. No relay, transcoding or HLS is added.
