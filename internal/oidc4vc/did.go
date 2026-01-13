package oidc4vc

import (
	"strings"
)

// DIDDocument は did:web 用の did.json 構造を表します。
type DIDDocument struct {
	Context            []string             `json:"@context"`
	ID                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod,omitempty"`
	AssertionMethod    []string             `json:"assertionMethod,omitempty"`
}

// VerificationMethod は公開鍵情報などを表します。
type VerificationMethod struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Controller   string                 `json:"controller"`
	PublicKeyJWK map[string]interface{} `json:"publicKeyJwk,omitempty"`
}

// GenerateDIDWeb はドメインから did:web ドキュメントを生成します。
func GenerateDIDWeb(issuerURL string, pubKeyJWK map[string]interface{}) (*DIDDocument, error) {
	// 1. スキームの除去 (https:// -> "")
	did := strings.TrimPrefix(issuerURL, "https://")
	did = strings.TrimPrefix(did, "http://") // 安全のため
	did = strings.TrimSuffix(did, "/")

	// 2. パスの ":" への置換 (example.com/repo -> example.com:repo)
	did = strings.ReplaceAll(did, "/", ":")

	didID := "did:web:" + did
	keyID := didID + "#key-1"

	return &DIDDocument{
		Context: []string{"https://www.w3.org/ns/did/v1", "https://w3id.org/security/suites/jws-2020/v1"},
		ID:      didID,
		VerificationMethod: []VerificationMethod{
			{
				ID:           keyID,
				Type:         "JsonWebKey2020",
				Controller:   didID,
				PublicKeyJWK: pubKeyJWK,
			},
		},
		AssertionMethod: []string{keyID},
	}, nil
}
