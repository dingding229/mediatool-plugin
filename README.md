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

## F2 typed extensions

Optional GetAPI/GetService/GetCommand/GetActions/GetModule now return Go handler
bindings. Runtime freezes them after InitPlugin and publishes descriptorVersion=2;
`plugin.describe` is an immutable snapshot, `plugin.invoke` uses deterministic
kind+ID handler identifiers (never function addresses or author-managed strings).
API paths are exact, no wildcards; HTTP credentials are checked by the host and
never forwarded. Services use cron/interval/date/once and optional scheduleSetting;
`hello-go` uses `scheduled:true` to enable its one-second demo timer. It also
shows API auth metadata, an internal command/action and typed TransferComplete
subscription. Module execution remains unsupported until F2-beta.

A timeout/crash can mean an unknown side effect. Do not blindly retry it. Event
identity is stable; hello uses it as its own idempotency key. Host validates the
same pinned manifest plus live grants independently of SDK checks.

F2 extensions bind Go functions directly. `OnChainPriority` publishes deterministic
priority metadata; `NameRecognition` is a strictly bounded hint contract, never a
file or media identity. The hello plugin demonstrates opt-in recognition rewrites
(`recognizeTitle`/`moduleTitle`), interception, `DownloadAdded`, JSON model binding,
two dashboard keys and an ESM entry importing `@mediatool/plugin-ui`. The host
bundles this facade and plugin sources in **one** Vite graph; no runtime download,
string rewriting, eval or duplicate Vue runtime is permitted. Same-page ESM is
trusted administrator-approved code, not a security sandbox.


## F3-A typed acquisition Host clients

`host/downloads`, `host/search`, `host/sites`, `host/subscriptions` share the
credential-free DTOs and exact method/scope/level table in `host/business`.
Pass `sdk.Host`'s `*host.Client` to `New`. Mutations require an explicit stable
`business.Operation{Key: "caller-request-id", Confirm: true}`; update/delete also
require the target revision where specified by the host. Retry an uncertain call
with the same key to inspect its recorded outcome, never invent a new key to
blindly repeat a side effect. Host checks declarations and live grants, not SDK
claims. `siteIds`, `hosts` and `downloaderIds` may narrow approved resources.

No credentials, tracker passkeys, download hints or arbitrary settings are DTOs.
`SiteRequest` yields status/matched/needs_action, not authentication headers or
raw HTML (Raw is explicitly unsupported). `downloads.retry` (a download is not an upload-transfer retry), `downloads.recheck` and
`subscriptions.completion` return `unsupported`; coverage only parses filenames,
not evidence that a TV subscription completed. `plugintest.Host.OnBusiness`
installs deterministic test behaviors and remains deny-by-default, rather than
pretending to be a real downloader/site.

Use `Downloads.UploadControl` with a transfer ID returned by `Downloads.Progress`; it preserves native upload-session/endpoint fencing. Manual unknown-submission duplicate-risk overrides remain administrator-only.

Site summaries expose only the normalized origin and compiled template ID needed for controlled requests; never username/password, credential-bearing query/path, Cookie, headers, or arbitrary integration settings.

Deletion takes `business.DeleteRequest`: `DeleteFiles:false` never removes payload files; `true` explicitly requests the native cleanup approval path. Both require original program ownership, imported state and original known-false HR evidence. There is no HR waiver.

### Media and file Host clients (F3-B)

`sdk/host/library`, `classification`, `metadata`, `files`, `storage`, `strm`
share the exact scope/level table in `sdk/host/business`. ContractVersion stays 2.
All mutations (including persisted previews) require a confirmed operation key;
revision-bearing native operations additionally require the matching revision.
File paths are root-relative, never storage/CDN URLs. `Storage.Mkdir/Relocate`
create previews only: inspect the returned digest, then explicitly `Apply` it.
Copy/upload/Resolve/STRM-content APIs are deliberately absent. `Files.ReadText`
returns unsupported (no audited safe text reader yet). Batch operations retain
native limits (20 file plans, 100 STRM files) even though the protocol ceiling is
200. Native-library ID selectors are not guessed: hosts reject unsupported
selector combinations rather than broadening access. Runtime capabilities are
the authority for availability; a typed method is not a guarantee of a configured
provider. See method golden and mock-host tests for live scope/revocation behavior.

### Media servers and playback boundary (F3-C)

The `mediaserver`, `catalog`, `mediausers`, `playback`, `mediaprobe` clients and
metadata native-details methods expose safe DTOs only. No playback/CDN URL,
stream resolver, HTTP image client, session activation or media credential
method exists. Catalog inputs deliberately have no userId/visibleLibraries.
Host-issued PluginPrincipal grants administrative catalog access only; it is
not a media user. Member Home/Discover/NextUp/UserData return unsupported until a
real authorized member actor is available. Emby catalog uses its mediaserver
adapter only; scoped item lookup or libraries without safe evidence fail closed.
User patches preserve omitted fields; ClearExpiry/ClearParentalRating explicitly
clear nullable values. They never accept passwords or create users. All writes
require a confirmed operation and native revisions when applicable. Playback and
player-access methods are read-only summaries. Probe retries only enqueue native
work; the SDK has no ffx/linksource dependency or execution interface.

### Platform Host APIs and release acceptance (F3-D)

Typed cache/integrations/jobs/maintenance/network/notifications/settings/system
clients share the method/scope table. Settings accepts only movieTemplate,
seriesTemplate, searchCacheSeconds, tmdbCacheSeconds (and approved settingKeys).
Jobs are instance-owned manual service wakes only; active work is not forcibly
cancelled. Cache invalidation needs explicit IDs in tmdb/search namespaces.
Deliveries/retry belong only to this instance's plugin.notice outbox events;
unknown send outcomes require duplicate-risk acknowledgement.
HTTP is exact-origin public GET/HEAD, profile=host, no caller auth/headers/body,
no redirects and no returned response body/headers. network + network.proxy.use
are both live-required. Proxy settings remain host-only. Mutable operations need
confirmed operation keys and native revisions; there is no SQL/Redis/backup/
115-auth/secret-settings/worker-control escape hatch.

Hello Go 1.3 adds a manually invoked `f3-read-check` action. It performs safe reads
across F3-A/B/C/D and records only booleans in its own `f3-read-check` KV. It also
asserts that an ungranted settings write is denied. It never runs automatically
or modifies business data. Expanded read permissions require administrator review.

### Manual tools (contractVersion 2)

Implement `sdk.ToolProvider.GetAgentTools` returning `sdk.Tool{Name, Description,
InputSchema, RequiredScopes, EffectKind, Handler}`. Declare names in
`manifest.extensions.agentTools` and `agent.tools` in permissions. Tool calls are
manual administrator operations, not an Agent/LLM loop. The host installs the
per-tool execution-scope ceiling and rechecks live grants on every Host call.
`effectKind` is `read`, `write`, or `dangerous`; it may only overstate, never
understate the shared Host method/scope effect table. Mutations require explicit
administrator confirmation. OS-level code is trusted as described by the host;
a descriptor is not a malicious-code sandbox.

Input is <=16KiB, output <=32KiB, schema <=8KiB (depth8/nodes128). Supported JSON
Schema keywords: type, title/description, properties/required/additionalProperties,
items, enum, minimum/maximum, min/maxLength and min/maxItems. No coercion,
remote `$ref`, regex, executable defaults, or ignored unknown keywords. The SDK
and host both validate input. Runtime invocation keeps stable handler IDs and
bounded stdio transport. `hello.echo` is pure; `hello.remember` writes only own KV.

Development order: change this repository, run vet/test/race/golden, push and wait
for its Actions; then bump the exact pin in mediatool, regenerate the resource
inventory, test and publish through main Actions -> immutable Docker Hub digest ->
backed-up deployment. Runtime compilation/download/upload is never supported.

## F6-α 验收插件

内置 `site-checkin`（宿主认证、布尔证据签到）和 `library-refresh-notify`（提交后事件、合并刷新及幂等通知）。
均默认不执行业务，使用 JSON 配置/状态页，不需要单独的源码上传、工具链或二进制包。
参见各插件 README 的 scope、unknown/needs_action、限额和合成测试说明。
`plugintest.Host.Client()` 可直接连接 typed clients；Start 仍用于真实 stdio 协议测试。
mock 非只读业务操作也保存幂等收据：相同 key/参数不重放，参数冲突拒绝，撤权即时拒绝。

通知 Send 的入队证据是 business.NotificationPending（pending），不是已送达；验收插件严格匹配该常量，unknown 不算成功。

## F6-β 验收插件

- `seeding-cleanup`：默认试运行、分享率/时长/标签/分类/站点规则；执行与删数据独立开关；
  由 Host `downloads.downloader.deleteOwned` 保持原所有权、HR、入库和清理证据保护。
- `watch-organize`：已授权根内单层有界补扫、持久静默队列；由 Host 预览/确认文件整理，
  不在 worker 里移动文件。具体边界见各插件 README。
- `library-refresh-notify` 1.1.0 的事件订阅使用 `Events.OnWindow(..., 7*24*time.Hour, handler)`：
  宿主投递、管理员 replay、排队执行都检查**原事实时间**。Host 副作用收据保留 7 日，
  插件仅清除已完成超过 8 日（含 24h 安全余量）、原事实超过重投窗且无 pending 的条目。
  200 条在窗内仍满则 needs_action。未声明窗口的其他订阅仍保留原有历史 replay 语义。
  `events.Envelope.OccurredAt` 由宿主 outbox/journal 提供，不由插件或回放请求指定。

## F6-γ 验收插件与窄 Host 契约

新增 wishlist-subscribe / torrent-assistant：默认只读预览，前者只接受公开来源 ID，
模糊匹配进入候选，唯一身份经 subscriptions.ensure 创建停用订阅；后者全部写操作经
宿主 Configure/Recheck 的 owned、HR、路径模式、读回、fence 和收据保护。
metadata.wishlist 不返回原始正文；私有 TMDB 收藏/OAuth 不支持，绝不借用登录态。
watch-organize 1.1.0 支持 depth 默认 3 / 最大 5 的有界递归，总条目仍最多 200，
只以叶文件 size/mtime 计算静默；链接/身份异常拒绝，不增加 worker 文件系统能力。
新插件仍必须在本库测试推送、CI 通过后，由主仓精确 pin，走 Actions→Hub 镜像发布；
任何插件均不支持上传、外部源码、运行时下载或编译。
