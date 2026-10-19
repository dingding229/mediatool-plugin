# 站点自动签到（site-checkin 1.0.0）

默认 `active=false`，初始化不请求站点。每实例绑定一个宿主已配置的账号；多个账号用分身。
配置 `accountId`、精确 `origin`、相对 `path`、`GET/POST`、`successMarker`、宿主 Asia/Shanghai 五段 cron。
可选 `precheckPath` 先只读检查当日状态：匹配后不再签到；验证码、非 2xx 或不可判定失败需人工处理。
Cookie、passkey、认证 header 只能配置在宿主站点账号内，不要放进插件配置的 path/body。

## 有界判断而非返回正文

不新增原始 body、选择器文本或任意正则接口。`sites.request/checkin` 使用已有 `sites.request`
危险 scope：宿主绑定已配置账号的精确 origin，禁止跳转和调用方 header，最多读取 1 MiB，
先移除宿主持有的认证值，再在宿主内比较显式声明的 **≤200 字节字面量** `successMarker`。
插件仅得到 `status/state/matched`，得不到 HTML、Cookie、Set-Cookie 或认证 URL。
这已满足本站点签到证据判断；若站点需要动态 CSRF/验证码，明确 needs_action，不猜测成功、不绕过验证。

同账号同 UTC 日的操作 key 稳定；宿主收据不重放 unknown，本实例也持久保留当日结果。
失败/验证码先保存 `status=needs_action` 并向现有通知 outbox 排队，随后返回 SDK `NeedsAction`。
宿主在完成服务/事件收据后将实例停用并保留 needs_action，必须管理员显式重新启用。
通知 pending（已持久入队）不代表已送达。正常日最多一条已入队通知；同日失败不会自动再签到。
Host 收据现有保留期 7 日，本实例当日 checkpoint 在跨进程重启时仍有效。

权限：自身 config/KV、ui.json、services.register、危险 sites.request、notifications.send。
测试只用 mock Host 和合成 HTTP 上游，不接真实站点。`testdata/` 固定 manifest/descriptor golden。

cron 时区沿用宿主 scheduling.Timezone（Asia/Shanghai），不是 worker 环境中的 TZ=UTC；每日幂等周期明确按 UTC 日，不能据此推断任意站点的日界。

checkpoint 包含请求身份摘要，不会把同日旧账号的成功套到新账号；同账号请求参数变更仍由宿主原操作 key 的参数摘要冲突保护，不能借配置编辑重放不确定副作用。
