# 下载器助手

默认 active=false、write=false、recheck=false。启用后每五分钟检查配置批量（1–10）任务。
只读状态可经 page、dashboard、GET admin /status、downloader-status 命令查看；assist 动作
仅运行当前已批准规则，默认不写。write 与 recheck 都是明确的危险开关；不重启下载器进程。

所有变更经 Host downloads.downloader.configure / downloads.recheck，实时 scope、fence、
审计、稳定收据；foreign 只读。append 标签不删除 owner，宿主追加本实例 plugin: 标记。
分类只可选择已有、安全 save path 且 auto_tmm 明确 false 的分类，成功读回后同步冻结 category。
正在整理、共享 hash、其他规则/实例保护、清理中、源恢复/未知意图都拒绝。

限速使用 KiB/s，0 表示不修改而非解除限制；做种小时/比例也以 0 表示不修改。
已有限制只允许收窄，HR 必须明确 false，其他维度和 inactive seeding 策略保留。
远程部分成功或读回无法确认，保留意图并 needs_action，不自动重试/回滚外部副作用。
错误或持续无进度默认发通知；只有显式允许才重新校验，requested 不冒充修复完成。
