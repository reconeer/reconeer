# Reconeer

Reconeer is a passive subdomain enumeration tool powered by the Reconeer API.

It provides fast and reliable subdomain discovery using Reconeer’s continuously updated dataset.

---

# Features

* Passive subdomain enumeration
* High-quality curated dataset
* JSON output support
* Optional API authentication
* Rate limiting support
* Concurrent scanning
* Automatic version update checks

---

# Installation

Direct install

```
go install -v github.com/reconeer/reconeer/cmd/reconeer@latest
```

or


Clone and build the CLI:

```bash
git clone https://github.com/reconeer/reconeer
cd reconeer
go build -o reconeer ./cmd/reconeer
```

Move it to your PATH if desired:

```bash
sudo mv reconeer /usr/local/bin/
```

Verify installation:

```bash
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

# Example API Request

You can also query the API directly.

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
https://www.reconeer.com/api/domain/example.com
```

Example response:

```json
{
  "subdomains": [
    {"subdomain":"api.example.com"},
    {"subdomain":"mail.example.com"}
  ]
}
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

