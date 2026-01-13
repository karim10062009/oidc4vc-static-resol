package config

// Config は issuer.yaml の構造を定義します。
type Config struct {
	IssuerURL             string                 `yaml:"issuer_url"`
	CredentialDefinitions []CredentialDefinition `yaml:"credential_definitions"`
}

// CredentialDefinition は発行する資格情報の定義を保持します。
type CredentialDefinition struct {
	ID                                  string    `yaml:"id"`
	Format                              string    `yaml:"format"`
	CryptographicBindingMethodsSupported []string `yaml:"cryptographic_binding_methods_supported"`
	Display                             []Display `yaml:"display"`
}

// Display は資格情報の表示に関する情報を保持します。
type Display struct {
	Name            string `yaml:"name"`
	BackgroundColor string `yaml:"background_color"`
}
