import argparse
from reconeer_api import fetch_domain, fetch_ip, fetch_subdomain
import os
import requests
from pprint import pprint

def main():
    parser = argparse.ArgumentParser(description="Reconeer API Client")
    parser.add_argument('--domain', help="Get subdomains and stats for a domain")
    parser.add_argument('--ip', help="Get subdomains linked to an IP")
    parser.add_argument('--subdomain', help="Get info about a specific subdomain")
    parser.add_argument('--compare', help="Compare a local file of subdomains to Reconeer DB for the given domain")
    parser.add_argument('--output', help="Save result to file")

    args = parser.parse_args()
    result = None

    if args.domain and not args.compare:
        result = fetch_domain(args.domain)

    elif args.ip:
        result = fetch_ip(args.ip)

    elif args.subdomain:
        result = fetch_subdomain(args.subdomain)

    elif args.compare:
        if not args.domain:
            print("[!] --domain is required when using --compare")
            return
        if not os.path.isfile(args.compare):
            print(f"[!] File not found: {args.compare}")
            return

        with open(args.compare, 'rb') as f:
            files = {'subdomainfile': (os.path.basename(args.compare), f)}
            url = f'https://www.reconeer.com/api/upload/{args.domain}'
            try:
                r = requests.post(url, files=files, timeout=60)
                r.raise_for_status()
                result = r.json()
            except Exception as e:
                print(f"[!] Upload failed: {e}")
                return

    else:
        parser.print_help()
        return

    if args.output and result:
        with open(args.output, 'w') as f:
            import json
            f.write(json.dumps(result, indent=2))
        print(f"[+] Output saved to {args.output}")
    elif result:
        pprint(result)

if __name__ == '__main__':
    main()
