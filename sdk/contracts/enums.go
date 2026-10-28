package contracts

type EventType string
type ChainEventType string

const (
	EventPluginReload             EventType      = "plugin.reloaded"
	EventPluginAction             EventType      = "plugin.action.invoked"
	EventPluginTriggered          EventType      = "plugin.custom.triggered"
	EventCommandExcute            EventType      = "command.executed"
	EventSiteDeleted              EventType      = "site.account.deleted"
	EventSiteUpdated              EventType      = "site.account.updated"
	EventSiteRefreshed            EventType      = "site.metrics.refreshed"
	EventTransferComplete         EventType      = "library.organized"
	EventTransferFailed           EventType      = "library.organize.failed"
	EventSubtitleTransferComplete EventType      = "library.subtitles.copied"
	EventSubtitleTransferFailed   EventType      = "library.subtitles.failed"
	EventAudioTransferComplete    EventType      = "transfer.audio.complete"
	EventAudioTransferFailed      EventType      = "transfer.audio.failed"
	EventDownloadAdded            EventType      = "download.added"
	EventHistoryDeleted           EventType      = "library.history.hidden"
	EventDownloadFileDeleted      EventType      = "download.files.deleted"
	EventDownloadDeleted          EventType      = "download.task.deleted"
	EventUserMessage              EventType      = "messaging.inbound"
	EventWebhookMessage           EventType      = "plugin.webhook.received"
	EventNoticeMessage            EventType      = "notification.queued"
	EventSubscribeAdded           EventType      = "subscription.created"
	EventSubscribeModified        EventType      = "subscription.updated"
	EventSubscribeDeleted         EventType      = "subscription.deleted"
	EventSubscribeComplete        EventType      = "subscription.coverage.completed"
	EventSystemError              EventType      = "system.operation.failed"
	EventMetadataScrape           EventType      = "metadata.scraped"
	EventModuleReload             EventType      = "module.registry.changed"
	EventConfigChanged            EventType      = "settings.changed"
	EventMessageAction            EventType      = "messaging.action.received"
	EventWorkflowExecute          EventType      = "workflow.started"
	EventAgentTokensUsage         EventType      = "agent.tokens.usage"
	ChainPluginDataReset          ChainEventType = "plugin.data.reset"
	ChainNameRecognize            ChainEventType = "name.recognize"
	ChainAuthVerification         ChainEventType = "auth.verification"
	ChainAuthIntercept            ChainEventType = "auth.intercept"
	ChainCommandRegister          ChainEventType = "command.register"
	ChainTransferRename           ChainEventType = "transfer.rename"
	ChainTransferRenameBuild      ChainEventType = "transfer.rename.build"
	ChainTransferIntercept        ChainEventType = "transfer.intercept"
	ChainTransferOverwriteCheck   ChainEventType = "transfer.overwrite.check"
	ChainResourceSelection        ChainEventType = "resource.selection"
	ChainResourceDownload         ChainEventType = "resource.download"
	ChainDiscoverSource           ChainEventType = "discover.source"
	ChainMediaRecognizeConvert    ChainEventType = "media.recognize.convert"
	ChainRecommendSource          ChainEventType = "recommend.source"
	ChainWorkflowExecution        ChainEventType = "workflow.execution"
	ChainStorageOperSelection     ChainEventType = "storage.operation"
	ChainAgentLLMProvider         ChainEventType = "agent.llm.provider"
	ChainSubscribeEpisodesRefresh ChainEventType = "subscribe.episodes.refresh"
	ChainSubscribeCompletionCheck ChainEventType = "subscribe.completion.check"
)

func KnownEventType(t EventType) bool {
	switch t {
	case EventPluginConfigChanged, EventPluginReload, EventPluginAction, EventPluginTriggered, EventCommandExcute, EventSiteDeleted, EventSiteUpdated, EventSiteRefreshed, EventTransferComplete, EventTransferFailed, EventSubtitleTransferComplete, EventSubtitleTransferFailed, EventDownloadAdded, EventHistoryDeleted, EventDownloadFileDeleted, EventDownloadDeleted, EventUserMessage, EventWebhookMessage, EventNoticeMessage, EventSubscribeAdded, EventSubscribeModified, EventSubscribeDeleted, EventSubscribeComplete, EventSystemError, EventMetadataScrape, EventModuleReload, EventConfigChanged, EventMessageAction, EventWorkflowExecute:
		return true
	default:
		return false
	}
}
func KnownChainEventType(t ChainEventType) bool {
	switch t {
	case ChainPluginDataReset, ChainNameRecognize, ChainAuthVerification, ChainAuthIntercept, ChainCommandRegister, ChainTransferRename, ChainTransferRenameBuild, ChainTransferIntercept, ChainTransferOverwriteCheck, ChainResourceSelection, ChainResourceDownload, ChainDiscoverSource, ChainMediaRecognizeConvert, ChainRecommendSource, ChainWorkflowExecution, ChainStorageOperSelection, ChainSubscribeEpisodesRefresh, ChainSubscribeCompletionCheck:
		return true
	default:
		return false
	}
}

type Scope string

const (
	ScopePluginConfigRead         Scope = "plugin.config.read"
	ScopePluginConfigWrite        Scope = "plugin.config.write"
	ScopePluginDataRead           Scope = "plugin.data.read"
	ScopePluginDataWrite          Scope = "plugin.data.write"
	ScopePluginStateRead          Scope = "plugin.state.read"
	ScopePluginInstancesRead      Scope = "plugin.instances.read"
	ScopeFsData                   Scope = "fs.data"
	ScopeFsMediaRead              Scope = "fs.media:read"
	ScopeFsMediaWrite             Scope = "fs.media:write"
	ScopeSystemExec               Scope = "system.exec"
	ScopeSystemShell              Scope = "system.shell"
	ScopeNetwork                  Scope = "network"
	ScopeNetworkPrivate           Scope = "network.private"
	ScopeNetworkProxyUse          Scope = "network.proxy.use"
	ScopeLibraryRead              Scope = "library.read"
	ScopeLibraryScan              Scope = "library.scan"
	ScopeLibraryRecognize         Scope = "library.recognize"
	ScopeLibraryOrganize          Scope = "library.organize"
	ScopeLibraryRefresh           Scope = "library.refresh"
	ScopeLibraryHistoryWrite      Scope = "library.history.write"
	ScopeClassificationRead       Scope = "classification.read"
	ScopeClassificationWrite      Scope = "classification.write"
	ScopeMetadataRead             Scope = "metadata.read"
	ScopeMetadataRefresh          Scope = "metadata.refresh"
	ScopeSubscriptionsRead        Scope = "subscriptions.read"
	ScopeSubscriptionsWrite       Scope = "subscriptions.write"
	ScopeSubscriptionsRules       Scope = "subscriptions.rules"
	ScopeSearchRead               Scope = "search.read"
	ScopeSearchRun                Scope = "search.run"
	ScopeDownloadsRead            Scope = "downloads.read"
	ScopeDownloadsAdd             Scope = "downloads.add"
	ScopeDownloadsControl         Scope = "downloads.control"
	ScopeDownloadsDelete          Scope = "downloads.delete"
	ScopeSitesRead                Scope = "sites.read"
	ScopeSitesRequest             Scope = "sites.request"
	ScopeSitesContentRaw          Scope = "sites.content.raw"
	ScopeSitesRefresh             Scope = "sites.refresh"
	ScopeSitesTemplatesRead       Scope = "sites.templates.read"
	ScopeSitesTemplatesSelect     Scope = "sites.templates.select"
	ScopeIntegrationsRead         Scope = "integrations.read"
	ScopeStorageRead              Scope = "storage.read"
	ScopeStorageWrite             Scope = "storage.write"
	ScopeStrmRead                 Scope = "strm.read"
	ScopeStrmRepair               Scope = "strm.repair"
	ScopeMediaserverRead          Scope = "mediaserver.read"
	ScopeMediaserverRefresh       Scope = "mediaserver.refresh"
	ScopeMediaserverUsersRead     Scope = "mediaserver.users.read"
	ScopeMediaserverUsersWrite    Scope = "mediaserver.users.write"
	ScopeMediaserverSettingsWrite Scope = "mediaserver.settings.write"
	ScopeMediaprobeRead           Scope = "mediaprobe.read"
	ScopeMediaprobeRetry          Scope = "mediaprobe.retry"
	ScopeNotificationsSend        Scope = "notifications.send"
	ScopeNotificationsRead        Scope = "notifications.read"
	ScopeNotificationsRetry       Scope = "notifications.retry"
	ScopeSettingsRead             Scope = "settings.read"
	ScopeSettingsWrite            Scope = "settings.write"
	ScopeJobsRead                 Scope = "jobs.read"
	ScopeJobsCancel               Scope = "jobs.cancel"
	ScopeCacheRead                Scope = "cache.read"
	ScopeCacheInvalidate          Scope = "cache.invalidate"
	ScopeSystemRead               Scope = "system.read"
	ScopeMaintenanceRead          Scope = "maintenance.read"
	ScopeEventsPublish            Scope = "events.publish"
	ScopeUiJson                   Scope = "ui.json"
	ScopeUiTrusted                Scope = "ui.trusted"
	ScopeUiSidebar                Scope = "ui.sidebar"
	ScopeApiRegister              Scope = "api.register"
	ScopeApiPublic                Scope = "api.public"
	ScopeApiPublicWrite           Scope = "api.public.write"
	ScopeServicesRegister         Scope = "services.register"
	ScopeCommandsRegister         Scope = "commands.register"
	ScopeCommandsOverride         Scope = "commands.override"
	ScopeActionsRegister          Scope = "actions.register"
	ScopeWorkflowsRun             Scope = "workflows.run"
	ScopeAuthProvider             Scope = "auth.provider"
	ScopeAgentTools               Scope = "agent.tools"
)

func KnownScope(s string) bool {
	switch s {
	case "plugin.config.read", "plugin.config.write", "plugin.data.read", "plugin.data.write", "plugin.state.read", "plugin.instances.read", "fs.data", "fs.media:read", "fs.media:write", "system.exec", "system.shell", "network", "network.private", "network.proxy.use", "library.read", "library.scan", "library.recognize", "library.organize", "library.refresh", "library.history.write", "classification.read", "classification.write", "metadata.read", "metadata.refresh", "subscriptions.read", "subscriptions.write", "subscriptions.rules", "search.read", "search.run", "downloads.read", "downloads.add", "downloads.control", "downloads.delete", "sites.read", "sites.request", "sites.content.raw", "sites.refresh", "sites.templates.read", "sites.templates.select", "integrations.read", "storage.read", "storage.write", "strm.read", "strm.repair", "mediaserver.read", "mediaserver.refresh", "mediaserver.users.read", "mediaserver.users.write", "mediaserver.settings.write", "mediaprobe.read", "mediaprobe.retry", "notifications.send", "notifications.read", "notifications.retry", "settings.read", "settings.write", "jobs.read", "jobs.cancel", "cache.read", "cache.invalidate", "system.read", "maintenance.read", "events.publish", "ui.json", "ui.trusted", "ui.sidebar", "api.register", "api.public", "api.public.write", "services.register", "commands.register", "commands.override", "actions.register", "workflows.run", "auth.provider", "agent.tools":
		return true
	case "events.subscribe.plugin.reloaded":
		return true
	case "events.subscribe.plugin.action.invoked":
		return true
	case "events.subscribe.plugin.custom.triggered":
		return true
	case "events.subscribe.command.executed":
		return true
	case "events.subscribe.site.account.deleted":
		return true
	case "events.subscribe.site.account.updated":
		return true
	case "events.subscribe.site.metrics.refreshed":
		return true
	case "events.subscribe.library.organized":
		return true
	case "events.subscribe.library.organize.failed":
		return true
	case "events.subscribe.library.subtitles.copied":
		return true
	case "events.subscribe.library.subtitles.failed":
		return true
	case "events.subscribe.download.added":
		return true
	case "events.subscribe.library.history.hidden":
		return true
	case "events.subscribe.download.files.deleted":
		return true
	case "events.subscribe.download.task.deleted":
		return true
	case "events.subscribe.messaging.inbound":
		return true
	case "events.subscribe.plugin.webhook.received":
		return true
	case "events.subscribe.notification.queued":
		return true
	case "events.subscribe.subscription.created":
		return true
	case "events.subscribe.subscription.updated":
		return true
	case "events.subscribe.subscription.deleted":
		return true
	case "events.subscribe.subscription.coverage.completed":
		return true
	case "events.subscribe.system.operation.failed":
		return true
	case "events.subscribe.metadata.scraped":
		return true
	case "events.subscribe.module.registry.changed":
		return true
	case "events.subscribe.plugin.config.changed":
		return true
	case "events.subscribe.settings.changed":
		return true
	case "events.subscribe.messaging.action.received":
		return true
	case "events.subscribe.workflow.started":
		return true
	case "events.chain.plugin.data.reset":
		return true
	case "events.chain.name.recognize":
		return true
	case "events.chain.auth.verification":
		return true
	case "events.chain.auth.intercept":
		return true
	case "events.chain.command.register":
		return true
	case "events.chain.transfer.rename":
		return true
	case "events.chain.transfer.rename.build":
		return true
	case "events.chain.transfer.intercept":
		return true
	case "events.chain.transfer.overwrite.check":
		return true
	case "events.chain.resource.selection":
		return true
	case "events.chain.resource.download":
		return true
	case "events.chain.discover.source":
		return true
	case "events.chain.media.recognize.convert":
		return true
	case "events.chain.recommend.source":
		return true
	case "events.chain.workflow.execution":
		return true
	case "events.chain.storage.operation":
		return true
	case "events.chain.subscribe.episodes.refresh":
		return true
	case "events.chain.subscribe.completion.check":
		return true
	case "module.metadata.name_recognize":
		return true
	case "module.metadata.identify":
		return true
	case "module.notifications.channel":
		return true
	case "module.downloads.provider":
		return true
	case "module.sites.search":
		return true
	case "module.mediaserver.refresh":
		return true
	case "module.library.output_plan":
		return true
	case "module.storage.selector":
		return true
	default:
		return false
	}
}
