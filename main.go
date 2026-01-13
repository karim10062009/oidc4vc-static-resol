package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/user/oidc4vc-static-resolver/internal/config"
	"github.com/user/oidc4vc-static-resolver/internal/generator"
	"github.com/user/oidc4vc-static-resolver/internal/oidc4vc"
	"github.com/user/oidc4vc-static-resolver/internal/vc"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: oidc4gen <command>")
		fmt.Println("Commands: build")
		return
	}

	command := os.Args[1]
	switch command {
	case "build":
		if err := runBuild(); err != nil {
			log.Fatalf("Build failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func runBuild() error {
	// 1. 設定ファイルの読み込み
	configFile := "issuer.yaml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Printf("Building metadata for issuer: %s\n", cfg.IssuerURL)

	// 2. 出力ディレクトリの作成
	wellKnownDir := filepath.Join("public", ".well-known")
	if err := os.MkdirAll(wellKnownDir, 0755); err != nil {
		return fmt.Errorf("failed to create .well-known directory: %w", err)
	}

	// 3. メタデータの初期化
	issuerMeta := oidc4vc.IssuerMetadata{
		Issuer:                            cfg.IssuerURL,
		AuthorizationServers:              []string{cfg.IssuerURL}, // 自身を認可サーバーとして模倣
		CredentialEndpoint:                cfg.IssuerURL + "/credentials/push-metadata",
		CredentialConfigurationsSupported: make(map[string]oidc4vc.CredentialConfiguration),
	}

	authMeta := oidc4vc.AuthServerMetadata{
		Issuer:                 cfg.IssuerURL,
		AuthorizationEndpoint:  cfg.IssuerURL + "/authorize", // 静的サイトなので実際には動作しない
		TokenEndpoint:          cfg.IssuerURL + "/token",     // 静的サイトなので実際には動作しない
		ResponseTypesSupported: []string{"code"},
		GrantTypesSupported:    []string{"urn:ietf:params:oauth:grant-type:pre-authorized_code"},
	}

	for _, d := range cfg.CredentialDefinitions {
		configID := d.ID
		var displays []oidc4vc.Display
		for _, disp := range d.Display {
			displays = append(displays, oidc4vc.Display{
				Name:            disp.Name,
				BackgroundColor: disp.BackgroundColor,
			})
		}
		issuerMeta.CredentialConfigurationsSupported[configID] = oidc4vc.CredentialConfiguration{
			Format:  d.Format,
			Display: displays,
		}
	}

	// 4. Pre-authorized VC の生成 (push-metadata として出力)
	credentialsDir := filepath.Join("public", "credentials")
	if err := os.MkdirAll(credentialsDir, 0755); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	signer, err := vc.NewSDJWT()
	if err != nil {
		return fmt.Errorf("failed to initialize signer: %w", err)
	}

	// 署名済みVCを生成
	signedVC, err := signer.Sign(cfg.IssuerURL, "did:example:123")
	if err != nil {
		return fmt.Errorf("failed to sign VC: %w", err)
	}

	// [修正] ファイル名を push-metadata に変更
	if err := os.WriteFile(filepath.Join(credentialsDir, "push-metadata"), []byte(signedVC), 0644); err != nil {
		return fmt.Errorf("failed to write VC file: %w", err)
	}

	// 5. JSONファイルの書き出し
	if err := writeJSON(filepath.Join(wellKnownDir, "openid-credential-issuer"), issuerMeta); err != nil {
		return err
	}

	// [追加] oauth-authorization-server の書き出し
	if err := writeJSON(filepath.Join(wellKnownDir, "oauth-authorization-server"), authMeta); err != nil {
		return err
	}

	// [追加] did.json の生成と書き出し
	didDoc, err := oidc4vc.GenerateDIDWeb(cfg.IssuerURL)
	if err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(wellKnownDir, "did.json"), didDoc); err != nil {
		return err
	}

	// 6. index.html の生成
	if err := generator.GenerateHTML(filepath.Join("public", "index.html"), cfg.IssuerURL); err != nil {
		return err
	}

	fmt.Println("Build completed successfully.")
	return nil
}

func writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON for %s: %w", path, err)
	}
	return nil
}
