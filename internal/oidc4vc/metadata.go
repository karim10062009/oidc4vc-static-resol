package oidc4vc

// IssuerMetadata は OpenID4VCI 認可サーバー/発行者のメタデータを表します。
// 規格: OpenID for Verifiable Credential Issuance
type IssuerMetadata struct {
	Issuer                         string                 `json:"credential_issuer"`
	AuthorizationServers           []string               `json:"authorization_servers,omitempty"`
	CredentialEndpoint             string                 `json:"credential_endpoint"`
	CredentialDefinitionsSupported []CredentialDefinition `json:"credential_definitions_supported"`
	Display                        []Display              `json:"display,omitempty"`
}

// AuthServerMetadata は OAuth 2.0 認和サーバーのメタデータを表します (RFC 8414)。
type AuthServerMetadata struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	GrantTypesSupported    []string `json:"grant_types_supported"`
}

// CredentialDefinition はサポートされる資格情報の定義です。
type CredentialDefinition struct {
	ID                                   string           `json:"id"` // oidc4vci では通常 configuration_id
	Format                               string           `json:"format"`
	CryptographicBindingMethodsSupported []string         `json:"cryptographic_binding_methods_supported,omitempty"`
	CredentialSigningAlgorithmsSupported []string         `json:"credential_signing_algorithms_supported,omitempty"`
	Display                              []Display        `json:"display,omitempty"`
	Claims                               map[string]Claim `json:"claims,omitempty"`
}

// Display は表示名やロゴなどの情報です。
type Display struct {
	Name            string `json:"name,omitempty"`
	Locale          string `json:"locale,omitempty"`
	Logo            *Logo  `json:"logo,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	TextColor       string `json:"text_color,omitempty"`
}

// Logo はロゴのURLや代替テキストです。
type Logo struct {
	URL     string `json:"url"`
	AltText string `json:"alt_text,omitempty"`
}

// Claim は資格情報に含まれる属性の定義です。
type Claim struct {
	Display []Display `json:"display,omitempty"`
}
