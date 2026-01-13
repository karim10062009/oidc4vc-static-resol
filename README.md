# oidc4vc-static-resolver

[Japanese version (日本語版)](README_ja.md)

## Overview

A CLI tool (`oidc4gen`) designed to generate and deploy OpenID4VCI (Issuance) and OpenID4VP (Presentation) compliant metadata and pre-signed SD-JWT Verifiable Credentials (VCs) on GitHub Pages. This project enables the setup of a serverless "VC Issuance Demo Site" in just a few minutes.

## Core Features

1. Metadata Generator
    - Generates .well-known directory structures compliant with RFC 9101 and OpenID4VCI standards.
    - Uses strong Go types to ensure accurate JSON schema generation.

2. SD-JWT VC Generation
    - Pre-generates "Pre-authorized" VCs at build time.
    - Supports SD-JWT (Selective Disclosure JWT) format with masked claims and disclosure hashes.

3. did:web Support
    - Automatically generates compliant did.json for the issuer.
    - Correctly handles domain and path mapping for did:web identifiers.

4. Interactive Demo HTML
    - Generates an index.html with a QR code containing openid-initiate-issuance:// custom schemes.
    - Allows direct VC acquisition from mobile wallet apps.

## System Architecture

The tool generates the following structure for GitHub Pages:

```text
/public
|-- .well-known
|   |-- openid-credential-issuer  # Issuer metadata
|   |-- did.json                  # did:web document
|   `-- oauth-authorization-server # Authorization server metadata
|-- credentials
|   `-- push-metadata             # Pre-signed SD-JWT VC
`-- index.html                    # GUI with QR code for wallet scanning
```

## Quick Start

1. Configure your issuer in `issuer.yaml`.
2. Run the build command:
   ```bash
   go run main.go build
   ```
3. Deploy the contents of the `public` directory to GitHub Pages.

## Technical Context

This project is built for the 2026 digital identity ecosystem, focusing on EUDI Wallet / eIDAS 2.0 interoperability and SD-JWT format standards.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
