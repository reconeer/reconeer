<h1 align="center">
  <img src="static/reconeer-logo.png" alt="reconeer" width="200px">
  <br>
</h1>

<h4 align="center">Fast passive subdomain + IP data enumeration tool.</h4>

<p align="center">
<a href="https://goreportcard.com/report/github.com/reconeer/reconeer"><img src="https://goreportcard.com/badge/github.com/reconeer/reconeer"></a>
<a href="https://github.com/reconeer/reconeer/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://github.com/reconeer/reconeer/releases"><img src="https://img.shields.io/github/release/reconeer/reconeer"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#api-setup">API Setup</a> •
  <a href="#library">Library</a>
</p>

---

`reconeer` is a subdomain + IP data discovery tool that retrieves known subdomains and resolution data via the [https://reconeer.com](https://reconeer.com) API.

## Features

- Simple JSON-based recon API client
- Subdomain/IP enumeration from the Reconeer platform
- Support for domain or file input
- Output to file or stdout
- Verbose/debug modes

## 📦 Installation

```bash
go install -v github.com/reconeer/reconeer/cmd/reconeer@latest
```

Requires **Go 1.21+**.

## ⚙️ Usage

```bash
reconeer -d example.com
reconeer -dL domains.txt
```

### Available Flags

```text
  -d,   --domain string[]     domains to fetch from the API
  -dL,  --list string         file with list of domains
  -o,   --output string       write output to file
  -json                      output full JSON records
  -ip                        include IP info in flat output
  -v                         verbose mode
  -silent                    only print subdomains
  -timeout int              timeout per request (default 30s)
```

## 🔐 API Setup

reconeer uses the public reconeer.com API. No key required at this stage.

## 📚 Library

The tool is modular. To use as a Go package, import:

```go
import "github.com/reconeer/reconeer/pkg/client"
```

and call:

```go
result, err := client.Fetch("example.com")
```

---

## License

MIT License © 2025 reconeer.com