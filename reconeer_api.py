import requests

BASE_URL = 'https://www.reconeer.com/api'

def fetch_domain(domain):
    try:
        r = requests.get(f'{BASE_URL}/domain/{domain}', timeout=60)
        r.raise_for_status()
        return r.json()
    except Exception as e:
        print(f"[!] Failed to fetch domain data: {e}")
        return None

def fetch_ip(ip):
    try:
        r = requests.get(f'{BASE_URL}/ip/{ip}', timeout=60)
        r.raise_for_status()
        return r.json()
    except Exception as e:
        print(f"[!] Failed to fetch IP data: {e}")
        return None

def fetch_subdomain(subdomain):
    try:
        r = requests.get(f'{BASE_URL}/subdomain/{subdomain}', timeout=60)
        r.raise_for_status()
        return r.json()
    except Exception as e:
        print(f"[!] Failed to fetch subdomain data: {e}")
        return None
