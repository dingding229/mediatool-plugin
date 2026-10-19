# 目录监控实时整理（1.0.0）

默认停用。管理员指定已配置、已授权、可写 rootId 内的 directory 和 destination。
每 30 秒用 Host fs.media.list 补扫 **单层** 目录，最多 200 项；不是 fsnotify 递归监听。
超过上限进入 needs_action，不静默漏项。暂停后没有后台 watcher 或 goroutine。

对普通媒体文件按 root/path/identity/size/mtime 生成稳定身份，静默期 30–3600 秒，
队列持久化到加密 KV；同一文件的变化合并。重启继续补扫，单轮最多提交 10 项。
先调 fs.media.preview 固定身份/size/mtime，再把 digest 交给 fs.media.apply。
**所有移动均由原 filemanager Engine 完成**：root、冲突、seeding source、owner、
符号链接逃逸、原始文件证据和 fence 继续强制。插件没有直接文件系统能力。

本插件提交 core organize；命名解析和模板由宿主决定，不猜 TMDB、不自动认领影视身份、
不自动确认 library 候选；需要智能命名或递归监控的场景不伪装支持。
不读写 STRM，不替换 Resolve，不提供播放／115 直链。未知结果保留执行意图并停用待处理；
崩溃后只读宿主既有 plan 的 succeeded 证据，绝不重放未知 move。

当前只有提供可靠 mtime 的本地根满足静默证据；CloudDrive2 列表没有 mtime，明确返回 unsupported，
不以 size-only 冒充稳定文件，也不新增取流或 115 请求。
