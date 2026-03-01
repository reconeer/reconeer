# Bug bounty recon pipeline with Reconeer

This guide shows a simple passive-first recon pipeline that scales well for bug bounty work.

## 1) Enumerate passive subdomains
```bash
export RECONEER_API_KEY="your_key_here"
reconeer -d example.com -silent | sort -u > subs.txt
```

## 2) Probe for live hosts
```bash
cat subs.txt | httpx -silent > live.txt
```

## 3) Optional: snapshot tech stack
```bash
cat live.txt | httpx -title -server -tech-detect -silent > tech.txt
```

## Free vs Premium
If you run this daily across many targets, the free tier (10 queries/day) will be limiting.
Premium unlocks unlimited queries and is designed for continuous automation:
https://www.reconeer.com/pricing?utm_source=github&utm_medium=doc&utm_campaign=oss_funnel
