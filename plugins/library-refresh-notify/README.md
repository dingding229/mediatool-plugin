# 媒体库刷新与入库通知（library-refresh-notify 1.0.0）

默认 `active=false`。配置 `active=true`、`provider=builtin|emby`、`intervalSeconds=60..3600`。
不自己启动线程/定时器：宿主 event lane 和持久 service 调度，停用后没有后台任务。

订阅 TransferComplete（`library.organized`）：只接收提交后、带合法 assetId 的媒体入库事实。
旧 files.* 的一般文件操作没有 assetId，直接忽略，绝不猜测资产或构造播放/取流 URL。
主程序在确认输出和恢复完成的事务中发新版身份事实；payload 只有 operationId/assetId/state。

每个 eventId 的通知通过稳定的 `notifications.publish` 操作 key 和宿主收据去重，
本实例先持久保存待办，再执行副作用，最后保存 checkpoint；崩溃后相同 key 不重复通知。
同一资产的待刷新请求合并；持久 service 全局每个配置间隔最多提交一个资产，立即手动再运行不突破节流。
`mediaserver.refresh` 只排现有业务任务，结果 requested/queued 不冒充刷新成功。
通知说明的是“入库已确认、进入插件待刷新队列”，不是“刷新已完成”。

队列/去重账本最多 200 个事件身份，不自动淘汰身份，防止旧事件在 Host 7 日收据过期后变成新副作用。
达到上限或刷新结果未知：通知、needs_action、实例停用，需要管理员处理；不静默丢任务或自动重放。
同库 SDK mock Host 有相同的 key/参数冲突和未知收据规则；测试覆盖重复、合并、节流、丢 checkpoint、撤权。

权限：自身 config/KV、ui.json、services.register、精确 events.subscribe.library.organized、
mediaserver.refresh、notifications.send。没有播放、115、原文 STRM 或媒体文件写权限。

1.1.0：取代旧“不淘汰”策略。该订阅由宿主强制七日原事实重投窗口，
完成超过八日且无 pending 引用的条目才可淘汰；未完成、旧版无时间证据、
或窗口内满载仍需处理。人工 replay 也不能越过该窗口，其他订阅不受影响。
