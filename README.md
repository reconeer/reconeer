<h1 align="center">
  <img src="static/reconeer-logo.png" alt="Reconeer" width="200px">
  <br>
</h1>

<h4 align="center">Fast subdomain enumeration client for reconeer.com API.</h4>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/reconeer/reconeer"><img src="https://goreportcard.com/badge/github.com/reconeer/reconeer"></a>
  <a href="https://github.com/reconeer/reconeer/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
  <a href="https://github.com/reconeer/reconeer/releases"><img src="https://img.shields.io/github/release/reconeer/reconeer"></a>
  <a href="https://twitter.com/reconeer"><img src="https://img.shields.io/twitter/follow/reconeer.svg?logo=twitter"></a>
  <a href="https://discord.gg/reconeer"><img src="https://img.shields.io/discord/123456789.svg?logo=discord"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Install</a> •
  <a href="#running-reconeer">Usage</a> •
  <a href="#api-setup">API Setup</a> •
  <a href="#reconeer-go-library">Library</a> •
  <a href="https://discord.gg/reconeer">Join Discord</a>
</p>

---

`Reconeer` is a subdomain enumeration tool designed to discover valid subdomains for websites using the reconeer.com API. It features a modular architecture optimized for speed and efficiency. `Reconeer` is built for one purpose—passive subdomain enumeration—and it excels at it.

The tool complies with the usage policies of the reconeer.com API. Its passive approach ensures rapid and discreet enumeration, making it ideal for penetration testers and bug bounty hunters.

# Features

<h1 align="left">
  <img src="static/reconeer-run.png" alt="Reconeer" width="700px">
  <br>
</h1>

- Fast and efficient subdomain enumeration via reconeer.com API
- **Curated** integration with reconeer.com for reliable results
- Multiple input and output options (file, stdout)
- Optimized for speed and lightweight resource usage
- **STDIN/OUT** support for seamless workflow integration

# Usage

```sh
reconeer -h
```

This will display help for the tool. Here are the supported switches:

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

# Installation

`Reconeer` requires **go1.24** to install successfully. Run the following command to install the latest version:

```sh
go install -v github.com/reconeer/reconeer@latest
```

## API Setup

After installation, configure your reconeer.com API key. Learn more here: https://docs.reconeer.com/setup.

## Running Reconeer

Learn about how to run Reconeer here: https://docs.reconeer.com/running.

## Reconeer Go Library

Reconeer can be used as a library. A minimal example is available [here](cmd/reconeer/examples/main.go).

### Resources

- [Recon with Reconeer!](https://reconeer.com/blog/recon-guide)

# License

`Reconeer` is made with 🖤 by the [reconeer](https://reconeer.com) team. See the **[THANKS.md](https://github.com/reconeer/reconeer/blob/main/THANKS.md)** for contributions.

Read the usage disclaimer at [DISCLAIMER.md](https://github.com/reconeer/reconeer/blob/main/DISCLAIMER.md) and [contact us](mailto:support@reconeer.com) for support.
```
