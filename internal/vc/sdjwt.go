package vc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// SDJWT は SD-JWT VC の生成と署名を担当します。
type SDJWT struct {
	PrivateKey *ecdsa.PrivateKey
}

// NewSDJWT は新しい署名器を作成します。
func NewSDJWT() (*SDJWT, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return &SDJWT{PrivateKey: key}, nil
}

// Sign は指定されたクレームに署名し、SD-JWT形式のトークンを返します。
func (s *SDJWT) Sign(issuer string, subject string) (string, error) {
	// 1. Disclosure (情報の開示) の生成
	// 本来は各属性に対して Salt を生成しハッシュ化しますが、ここでは簡略化して1つの属性をDisclosure化します。
	salt := "random_salt_123"
	claimName := "given_name"
	claimValue := "Toshiki"

	disclosureData := []interface{}{salt, claimName, claimValue}
	disclosureJSON, err := json.Marshal(disclosureData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal disclosure: %w", err)
	}
	disclosureB64 := base64.RawURLEncoding.EncodeToString(disclosureJSON)

	hash := sha256.Sum256([]byte(disclosureB64))
	disclosureHash := base64.RawURLEncoding.EncodeToString(hash[:])

	// 2. JWT ペイロードの作成 (_sd フィールドにハッシュを入れる)
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"_sd": []string{disclosureHash},
		"cnf": map[string]interface{}{
			"jwk": map[string]interface{}{
				"kty": "EC",
				"crv": "P-256",
				"x":   base64.RawURLEncoding.EncodeToString(s.PrivateKey.PublicKey.X.Bytes()),
				"y":   base64.RawURLEncoding.EncodeToString(s.PrivateKey.PublicKey.Y.Bytes()),
			},
		},
	})

	signedJWT, err := token.SignedString(s.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// 3. SD-JWT 形式の組み立て (JWT ~ Disclosure1 ~ Disclosure2 ~ ...)
	sdjwt := fmt.Sprintf("%s~%s~", signedJWT, disclosureB64)

	return sdjwt, nil
}

// PublicKeyJWK は公開鍵を JWK 形式のマップで返します。
func (s *SDJWT) PublicKeyJWK() map[string]interface{} {
	return map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   base64.RawURLEncoding.EncodeToString(s.PrivateKey.PublicKey.X.Bytes()),
		"y":   base64.RawURLEncoding.EncodeToString(s.PrivateKey.PublicKey.Y.Bytes()),
	}
}
