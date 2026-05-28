# Reconeer

Passive subdomain enumeration powered by the Reconeer API (https://www.reconeer.com).

Fast, no brute-force, no active DNS probes — just observed data from a continuously updated dataset.

---

## Get started in 30 seconds

1. Create a free account at https://www.reconeer.com/register
2. Copy your API key from https://www.reconeer.com/dashboard
3. Install the CLI: `go install github.com/reconeer/reconeer/cmd/reconeer@latest`
4. Run: `reconeer -d example.com -k YOUR_API_KEY`

That's it. Or skip the CLI entirely and use it with **Subfinder** (default source since [PR #1694](https://github.com/projectdiscovery/subfinder/pull/1694)):

```bash
export RECONEER_API_KEY=YOUR_API_KEY
subfinder -d example.com
```

---

# Features

* Passive subdomain enumeration — no active DNS, no brute force, no portscans
* High-quality curated dataset (continuously refreshed)
* Three ways to call: Reconeer CLI, subfinder integration, raw curl/JSON
* JSON Lines output (`-jsonl`) — pipes cleanly into `httpx`, `jq`, etc.
* Free tier: 10 queries/day. Premium: unlimited + CSV/TSV export + bulk endpoint

---

# Installation

```bash
go install -v github.com/reconeer/reconeer/cmd/reconeer@latest
```

Or build from source:

```bash
git clone https://github.com/reconeer/reconeer
cd reconeer
go build -o reconeer ./cmd/reconeer
sudo mv reconeer /usr/local/bin/   # optional
reconeer -version
```

---

# Usage

Basic scan:

```bash
reconeer -d example.com
```

Scan multiple domains:

```bash
reconeer -d example.com -d example.org
```

Scan domains from file:

```bash
reconeer -dL domains.txt
```

Increase concurrency:

```bash
reconeer -dL domains.txt -t 10
```

Write results to file:

```bash
reconeer -d example.com -o results.txt
```

JSON output:

```bash
reconeer -d example.com -jsonl
```

Silent output:

```bash
reconeer -d example.com -silent
```

---

# Authentication

Reconeer supports optional API keys for higher rate limits.

You can provide your API key in two ways.

### CLI flag

```bash
reconeer -d example.com -k YOUR_API_KEY
```

### Environment variable

```bash
export RECONEER_API_KEY=YOUR_API_KEY
reconeer -d example.com
```

Internally the CLI sends requests using:

```
Authorization: Bearer <API_KEY>
```

---

# Direct API access

The API accepts your key in any of three ways:

```bash
# 1. Authorization header (the standard)
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://www.reconeer.com/api/domain/example.com

# 2. X-API-KEY header (used by Subfinder)
curl -H "X-API-KEY: YOUR_API_KEY" \
  https://www.reconeer.com/api/domain/example.com

# 3. Query parameter (for quick browser tests)
curl "https://www.reconeer.com/api/domain/example.com?api_key=YOUR_API_KEY"
```

Example response:

```json
{
  "subdomains": [
    {"subdomain":"api.example.com","ip":"93.184.216.34","country":"US","reverse_resolves":true},
    {"subdomain":"mail.example.com","ip":"93.184.216.34","country":"US","reverse_resolves":false}
  ],
  "ipStats": [{"ip":"93.184.216.34","count":2}]
}
```

## Premium endpoints

```bash
# CSV / TSV export — drops a spreadsheet-ready file
curl -H "Authorization: Bearer YOUR_API_KEY" \
  "https://www.reconeer.com/api/domain/target.com?format=csv" \
  -o target.com_subs.csv

# Bulk query — up to 50 domains per request
curl -X POST -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"domains":["a.com","b.com","c.com"]}' \
  https://www.reconeer.com/api/domain/bulk
```

# Subfinder integration

Reconeer is a **default source** in [projectdiscovery/subfinder](https://github.com/projectdiscovery/subfinder) since [PR #1694](https://github.com/projectdiscovery/subfinder/pull/1694) — every subfinder user can pull from Reconeer with zero config changes:

```bash
export RECONEER_API_KEY=YOUR_API_KEY
subfinder -d example.com
```

Or use *only* Reconeer (skip other sources):

```bash
subfinder -d example.com -s reconeer
```

You can also put the key in `~/.config/subfinder/provider-config.yaml`:

```yaml
reconeer:
  - YOUR_API_KEY
```

---

# Output

Default output:

```
api.example.com
mail.example.com
dev.example.com
```

JSONL output:

```json
{"subdomain":"api.example.com"}
{"subdomain":"mail.example.com"}
```

---

# Rate Limits

Free tier:

* 10 queries per day

Using an API key provides higher limits and isolation.

Register for a free API key:

https://www.reconeer.com/register

---

# Updating

Check for updates:

```bash
reconeer -update
```

You can also download the latest release manually:

https://github.com/reconeer/reconeer/releases

---

# CLI Flags

| Flag       | Description             |
| ---------- | ----------------------- |
| `-d`       | Domain to enumerate     |
| `-dL`      | File containing domains |
| `-k`       | API key                 |
| `-rl`      | Requests per second     |
| `-t`       | Number of workers       |
| `-o`       | Output file             |
| `-jsonl`   | JSONL output            |
| `-silent`  | Only print subdomains   |
| `-v`       | Verbose logging         |
| `-version` | Print version           |
| `-update`  | Check for updates       |

---

# License

MIT License

---

# Links

Website
https://www.reconeer.com

API documentation
https://www.reconeer.com/docs

Register for API key
https://www.reconeer.com/register

