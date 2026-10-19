# 自动做种清理（1.0.0）

默认 active=false / execute=false / deleteFiles=false。每十分钟（宿主时区）扫描最多
200 个本程序下载任务，游标分页；按分享率、做种小时、分类、标签、站点账户 ID 筛选，
每批最多 20 项。状态页展示 preview/retired/pending/unknown；pending 不是文件已删除。

宿主 live Inspect 的 CleanupEligible 只是一份预览：实际 DeleteOwned 重新验证原 endpoint、
hash、category、save path 和所有权 tag、HR=false 来源快照、imported、无在途 transfer、
无共享 hash、无其他清理请求／保护标记。foreign 只读。keep/protect/no-cleanup、其他任务
UUID、plugin: 或 rule: tag 都是不可覆盖的保护。文件删除额外走原有云端/共享/路径/lease 证据链。

执行意图先写加密 KV，稳定 key 交给宿主幂等收据。崩溃中的 executing/unknown 不自动重放，
需人工核对后处理。变更危险开关仅代表管理员配置授权，不绕过宿主检查。无后台 goroutine。
