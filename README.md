# Reconeer

<div align="center">
  <img src="static/reconeer-logo.png" alt="Reconeer" width="200">
</div>

<h4 align="center">A high-performance subdomain enumeration client for the reconeer.com API.</h4>

<div align="center">
  <a href="https://goreportcard.com/report/github.com/reconeer/reconeer"><img src="https://goreportcard.com/badge/github.com/reconeer/reconeer" alt="Go Report Card"></a>
  <a href="https://github.com/reconeer/reconeer/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="Contributions Welcome"></a>
  <a href="https://github.com/reconeer/reconeer/releases"><img src="https://img.shields.io/github/release/reconeer/reconeer" alt="GitHub Release"></a>
  <a href="https://x.com/reconeerx"><img src="https://img.shields.io/twitter/follow/reconeer.svg?logo=twitter" alt="Twitter Follow"></a>
  <a href="https://discord.gg/reconeer"><img src="https://img.shields.io/discord/123456789.svg?logo=discord" alt="Discord"></a>
</div>

<div align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#api-setup">API Setup</a> •
  <a href="#reconeer-go-library">Go Library</a> •
  <a href="https://discord.gg/reconeer">Join Discord</a>
</div>

---

Reconeer is a specialized subdomain enumeration tool designed to identify valid subdomains for target domains using the reconeer.com API. It employs a modular and efficient architecture to deliver rapid results. Focused exclusively on passive subdomain enumeration, Reconeer excels in providing discreet and high-speed discovery, making it an essential asset for security professionals, penetration testers, and bug bounty researchers.

While the reconeer.com website offers a user-friendly graphical interface (UI) for subdomain enumeration, Reconeer serves as a command-line interface (CLI) client that interacts directly with the reconeer.com API to fetch data programmatically. This API-driven approach ensures compliance with reconeer.com's usage policies and enables seamless integration into automated workflows. For more information on reconeer.com's services, visit [reconeer.com](https://reconeer.com).

<div align="left">
  <img src="static/reconeer-run.png" alt="Reconeer in action" width="700">
</div>

## Features

- High-speed subdomain enumeration powered by the reconeer.com API
- Curated API integration for accurate and reliable results
- Support for multiple input and output formats (files, stdout)
- Lightweight design with minimal resource consumption
- STDIN/STDOUT compatibility for easy pipeline integration

## Usage

To view the help menu, run:

```sh
reconeer -h
```

The tool supports the following options:

```yaml
Usage:
  ./reconeer [flags]

Flags:
INPUT:
  -d, -domain string[]  domains to enumerate subdomains for
  -dL, -list string     file containing list of domains for enumeration

RATE-LIMIT:
  -rl, -rate-limit int  maximum number of API requests per second

OUTPUT:
  -o, -output string    file to write output to

CONFIGURATION:
  -config string        config file (default "$CONFIG/reconeer/config.yaml")

DEBUG:
  -silent             show only subdomains in output
  -version            show version of reconeer
  -v                  show verbose output
  -nc, -no-color      disable color in output
```

## Installation

Reconeer requires Go version 1.24 or higher for installation. Install the latest version using the following command:

```sh
go install -v github.com/reconeer/reconeer/cmd/reconeer@latest
```

## API Setup

After installation, configure your reconeer.com API key to enable access to the subdomain enumeration service. For detailed setup instructions, refer to the API documentation available on [reconeer.com](https://reconeer.com).

## Running Reconeer

To get started with running Reconeer, consult the usage guidelines and examples provided on [reconeer.com](https://reconeer.com).

## Reconeer Go Library

Reconeer can also be utilized as a Go library for custom integrations. A minimal example is available [here](cmd/reconeer/examples/main.go).

### Resources

- [Recon with Reconeer!](https://reconeer.com/blog/recon-guide)

## License

Reconeer is developed with passion by the [reconeer.com](https://reconeer.com) team. We extend our gratitude to all contributors—see **[THANKS.md](https://github.com/reconeer/reconeer/blob/main/THANKS.md)** for details.

Please review the usage disclaimer in **[DISCLAIMER.md](https://github.com/reconeer/reconeer/blob/main/DISCLAIMER.md)**. For support inquiries, contact us at [support@reconeer.com](mailto:support@reconeer.com).
