# TMDB/豆瓣想看订阅

默认 active=false、create=false。公开来源限 tmdb-list、douban-wish、douban-subject，
只填 ID，不接受 URL、Cookie、会话或凭据。TMDB app token 只在宿主；私有收藏/OAuth 不支持。
豆瓣仅公开固定页面，登录/验证码/页面格式变化返回 uncertain，不绕过站点限制。

宿主只返回白名单身份 DTO，不返回正文、图片或播放 URL。每次一页最多 20 项；
配置每轮 1–20 项，最短间隔 30 分钟；失败指数退避 5 分钟至 6 小时，待确认最多 200 条。
TMDB 明确的 type/id 可直接确认；豆瓣仅标题和年份都精确且唯一时自动确认。
其他情况保留候选，不打分猜测；状态页提供 sourceId/candidates，管理员执行 resolve 动作，
JSON 填 sourceId、tmdbId、mediaType，只能选择先前提供的候选组合。

创建经 subscriptions.ensure，existing/已入库不重复；默认创建停用的 search 订阅，
不立即搜索、不下载、不触发 RSS。用户应另在订阅页审核后启用。配置 accountIds 为宿主已有 PT 账户 ID。
副作用失败/unknown 不增加“确保数”；重新尝试使用稳定身份操作 key 和永久业务重复检查。
服务每 30 分钟唤醒，受 nextAt 配置间隔/退避控制；停用没有后台请求。
