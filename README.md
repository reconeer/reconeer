# Reconeer CLI

A simple Python CLI to interact with the Reconeer API.

## Install
```bash
pip install -r requirements.txt
```

## Usage
```bash
# Get all subdomains for a domain
python3 reconeer_cli.py --domain netflix.com

# Get info for an IP
python3 reconeer_cli.py --ip 45.57.89.23

# Get info about a specific subdomain
python3 reconeer_cli.py --subdomain cdn.netflix.com

# Compare a local file to the domain in Reconeer
python3 reconeer_cli.py --domain post.ch --compare my_subdomains.txt

# Save to file
python3 reconeer_cli.py --domain post.ch --output results.json
```


