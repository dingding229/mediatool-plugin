package business

// Method is the single wire-method/scope/level manifest shared by SDK, host and tests.
type Method struct {
	Name      string `json:"name"`
	Scope     string `json:"scope"`
	Level     string `json:"level"`
	Supported bool   `json:"supported"`
}

func Methods() []Method {
	return []Method{
		{"downloads.list", "downloads.read", "R", true},
		{"downloads.get", "downloads.read", "R", true},
		{"downloads.progress", "downloads.read", "R", true},
		{"downloads.files", "downloads.read", "R", true},
		{"downloads.trackers", "downloads.read", "R", true},
		{"downloads.manualList", "downloads.read", "R", true},
		{"downloads.downloader.allHashes", "downloads.read", "R", true},
		{"downloads.downloader.inspect", "downloads.read", "R", true},
		{"downloads.cleanup.preview", "downloads.read", "R", true},
		{"downloads.addResource", "downloads.add", "D", true},
		{"downloads.manualCreate", "downloads.add", "D", true},
		{"downloads.pause", "downloads.control", "W", true},
		{"downloads.resume", "downloads.control", "W", true},
		{"downloads.retry", "downloads.control", "W", false},
		{"downloads.uploadControl", "downloads.control", "W", true},
		{"downloads.manualRetry", "downloads.control", "W", true},
		{"downloads.downloader.setPaused", "downloads.control", "W", true},
		{"downloads.cleanup.approve", "downloads.delete", "D", true},
		{"downloads.downloader.deleteOwned", "downloads.delete", "D", true},
		{"metadata.torrent.parse", "metadata.read", "R", true},
		{"search.list", "search.read", "R", true},
		{"search.get", "search.read", "R", true},
		{"search.resources", "search.read", "R", true},
		{"search.start", "search.run", "W", true},
		{"sites.list", "sites.read", "R", true},
		{"sites.get", "sites.read", "R", true},
		{"sites.metrics", "sites.read", "R", true},
		{"sites.rss.discover", "sites.read", "R", true},
		{"sites.request", "sites.request", "D", true},
		{"sites.checkin", "sites.request", "D", true},
		{"sites.metrics.collect", "sites.refresh", "W", true},
		{"sites.rss.refresh", "sites.refresh", "W", true},
		{"sites.templates.catalog", "sites.templates.read", "R", true},
		{"sites.templates.loaded", "sites.templates.read", "R", true},
		{"sites.templates.load", "sites.templates.select", "D", true},
		{"subscriptions.list", "subscriptions.read", "R", true},
		{"subscriptions.get", "subscriptions.read", "R", true},
		{"subscriptions.history", "subscriptions.read", "R", true},
		{"subscriptions.create", "subscriptions.write", "D", true},
		{"subscriptions.update", "subscriptions.write", "D", true},
		{"subscriptions.queue", "subscriptions.write", "D", true},
		{"subscriptions.delete", "subscriptions.write", "D", true},
		{"subscriptions.evaluateRules", "subscriptions.rules", "R", true},
		{"subscriptions.coverage", "subscriptions.rules", "R", true},
		{"downloads.recheck", "downloads.control", "W", false},
		{"subscriptions.completion", "subscriptions.read", "R", false},
	}
}
func Lookup(name string) (Method, bool) {
	for _, m := range Methods() {
		if m.Name == name {
			return m, true
		}
	}
	return Method{}, false
}
