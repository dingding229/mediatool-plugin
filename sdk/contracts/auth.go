package contracts

// Broker assertions are not sessions. No roles or permissions are accepted from
// a provider. IsAdmin is decoded solely to prove that it has no authority.
type AuthRequest struct {
	ProviderID string `json:"providerId"`
	Challenge  string `json:"challenge"`
}
type AuthAssertion struct {
	Subject   string `json:"subject"`
	Issuer    string `json:"issuer"`
	Audience  string `json:"audience"`
	ExpiresAt int64  `json:"expiresAt"`
	IsAdmin   bool   `json:"isAdmin,omitempty"`
}
type AuthPolicy struct {
	AccountID string `json:"accountId"`
	Allow     bool   `json:"allow"`
}
