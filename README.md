# Reconeer CLI

<div align="center">
  <img src="static/reconeer-logo.png" alt="Reconeer" width="180">
</div>

<p align="center">
  <b>Passive subdomain enumeration</b> powered by the <a href="https://www.reconeer.com/docs.html">reconeer.com API</a>.
  Built for recon pipelines, bug bounty workflows, and automation.
</p>

<div align="center">
  <a href="https://github.com/reconeer/reconeer/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/reconeer/reconeer/ci.yml?branch=main" alt="CI"></a>
  <a href="https://goreportcard.com/report/github.com/reconeer/reconeer"><img src="https://goreportcard.com/badge/github.com/reconeer/reconeer" alt="Go Report Card"></a>
  <a href="https://github.com/reconeer/reconeer/releases"><img src="https://img.shields.io/github/v/release/reconeer/reconeer" alt="Release"></a>
  <a href="https://discord.gg/reconeer"><img src="https://img.shields.io/badge/Discord-join-5865F2?logo=discord&logoColor=white" alt="Discord"></a>
  <a href="https://x.com/reconeerx"><img src="https://img.shields.io/twitter/follow/reconeerx?logo=twitter" alt="X"></a>
  <a href="https://www.reconeer.com/signup?utm_source=github&utm_medium=badge&utm_campaign=oss_funnel"><img src="https://img.shields.io/badge/Get%20API%20Key-free-success" alt="Get API Key"></a>
</div>

---

## Why Reconeer?

Reconeer is **passive-only**: it returns intelligence from previously observed infrastructure (no brute forcing, no DNS guessing, no active probing).
That makes it safe to run continuously in automation.

- ✅ Fast passive subdomain enumeration
- ✅ CLI built for pipelines (STDIN/STDOUT)
- ✅ Optional API key support (`X-API-Key`) for higher limits and isolation
- ✅ Designed to integrate with recon tooling and custom scripts

---

## Quick start

### 1) Get a free API key (recommended)

Create a free account (10 API queries/day):

- https://www.reconeer.com/signup?utm_source=github&utm_medium=readme&utm_campaign=oss_funnel

Set your key:

```bash
export RECONEER_API_KEY="your_key_here"
```

### 2) Enumerate a domain

```bash
reconeer -d example.com
```

Write to a file:

```bash
reconeer -d example.com -o results.txt
```

### 3) Bulk enumerate many domains

```bash
reconeer -dL domains.txt -o results.txt
```

> If you hit the free daily limit, Reconeer will tell you and include a direct upgrade link.

<div align="left">
  <img src="static/reconeer-run.png" alt="Reconeer in action" width="700">
</div>

---

## Installation

### Option A: download a release binary
Grab the latest release from GitHub releases:
- https://github.com/reconeer/reconeer/releases

### Option B: install with Go

Requires Go **1.22+**.

```bash
go install -v github.com/reconeer/reconeer/cmd/reconeer@latest
```

---

## Usage

```text
reconeer -d example.com
reconeer -dL domains.txt
cat domains.txt | reconeer
```

Flags:

```text
INPUT
  -d, --domain <domain>       Domain to enumerate (repeatable)
  -dL, --list <file>          File with one domain per line; supports "-" for STDIN

AUTH
  -k, --api-key <key>         API key (or set RECONEER_API_KEY)

RATE LIMIT
  -rl, --rate-limit <n>       Max requests per second (client-side). Default: 3

OUTPUT
  -o, --output <file>         Write output to file (default: STDOUT)
  --jsonl                      Output JSON Lines (one object per subdomain)

MISC
  -silent                      Print only subdomains (no banners/logs)
  -v, --verbose                Verbose logging
  -version                     Print version
  -h, --help                   Help
```

---

## Free vs Premium (why upgrade?)

Reconeer pricing: https://www.reconeer.com/pricing?utm_source=github&utm_medium=readme&utm_campaign=oss_funnel

**Free ($0/mo)**  
- 10 API queries/day
- Basic subdomain enumeration
- CLI access

**Premium ($49/mo)**  
- Unlimited API queries
- Advanced analytics
- Priority integrations & early access
- Email support

---

## Integrations

### subfinder
Reconeer is designed to slot into recon pipelines. If you already use subfinder, add Reconeer as a passive source and provide your API key.

### Example pipeline
```bash
subfinder -d example.com -silent | sort -u | tee subs.txt
cat subs.txt | httpx -silent
```

---

## Developer docs

- API docs: https://www.reconeer.com/docs.html
- Guides:
  - docs/bugbounty-pipeline.md
  - docs/ci-monitoring.md

---

## Contributing

See **CONTRIBUTING.md** and open an issue. Good-first-issues are tagged.

---

## License & safety

This project is intended for legitimate security testing and research.  
See **DISCLAIMER.md** and **SECURITY.md**.
