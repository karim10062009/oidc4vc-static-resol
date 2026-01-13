package generator

import (
	"fmt"
	"os"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OpenID4VC Static Issuer Demo</title>
    <script src="https://cdn.jsdelivr.net/npm/qrcode@1.5.1/build/qrcode.min.js"></script>
    <style>
        body { font-family: sans-serif; display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #f0f2f5; }
        .card { background: white; padding: 2rem; border-radius: 1rem; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
        #qrcode { margin: 1.5rem 0; }
        .btn { background: #007bff; color: white; padding: 0.5rem 1rem; border-radius: 0.5rem; text-decoration: none; display: inline-block; margin-top: 1rem; }
    </style>
</head>
<body>
    <div class="card">
        <h1>VC発行デモ</h1>
        <p>ウォレットアプリで以下のQRコードをスキャンして、<br>デジタル資格情報を取得してください。</p>
        <div id="qrcode"></div>
        <p><small>Issuer: %s</small></p>
        <a href="%s" class="btn">直接リンクを開く (モバイル用)</a>
    </div>
    <script>
        const issuanceUrl = "%s";
        QRCode.toCanvas(document.getElementById('qrcode'), issuanceUrl, { width: 256 }, function (error) {
            if (error) console.error(error);
        });
    </script>
</body>
</html>
`

// GenerateHTML は QRコードを含む index.html を出力します。
func GenerateHTML(outputPath string, issuerURL string) error {
	// 規格に基づいたカスタムスキーム URL (openid-initiate-issuance://)
	// credential_configuration_id を指定して、どの資格情報のデモかを示す
	issuanceUrl := fmt.Sprintf("openid-initiate-issuance://?issuer=%s&credential_configuration_id=UniversityDegree", issuerURL)

	content := fmt.Sprintf(htmlTemplate, issuerURL, issuanceUrl, issuanceUrl)

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}
	return nil
}
