package contracts

import "encoding/json"

type Initialize struct {
	ContractVersion       int      `json:"contractVersion"`
	InstanceID            string   `json:"instanceId"`
	Generation            string   `json:"generation"`
	ConfigRevision        int64    `json:"configRevision"`
	DefinitionFingerprint string   `json:"definitionFingerprint"`
	CatalogBuild          string   `json:"catalogBuild"`
	Capabilities          []string `json:"capabilities"`
}
type Initialized struct {
	ContractVersion    int      `json:"contractVersion"`
	SchemaFingerprint  string   `json:"schemaFingerprint"`
	CatalogBuild       string   `json:"catalogBuild"`
	ImplementedMethods []string `json:"implementedMethods"`
}
type Config struct {
	Value    json.RawMessage `json:"value"`
	Revision int64           `json:"revision"`
}
type Value = Config
type Revision struct {
	Revision int64 `json:"revision"`
}
type Message struct {
	Title     string `json:"title"`
	Text      string `json:"text"`
	ChannelID string `json:"channelId,omitempty"`
}
type HostCall struct {
	Method           string          `json:"method"`
	Args             json.RawMessage `json:"args"`
	ExpectedRevision *int64          `json:"expectedRevision,omitempty"`
	OperationKey     string          `json:"operationKey,omitempty"`
}
type Component struct {
	Component string                     `json:"component"`
	Key       string                     `json:"key,omitempty"`
	Props     map[string]json.RawMessage `json:"props,omitempty"`
	Model     string                     `json:"model,omitempty"`
	Children  []Component                `json:"children,omitempty"`
}
type Form struct {
	Content  []Component     `json:"content"`
	Defaults json.RawMessage `json:"defaults"`
}
type PageRequest struct {
	Key string `json:"key"`
}
type Page struct {
	Content []Component `json:"content"`
}
type DashboardRequest = PageRequest
type Dashboard = Page
type DashboardMeta struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Refresh int    `json:"refresh"`
	Border  bool   `json:"border"`
	Cols    int    `json:"cols"`
}
