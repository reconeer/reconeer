# CI monitoring with Reconeer

Reconeer can be used to monitor domains continuously and alert on new assets.

## Example: daily enumeration (cron/CI)
```bash
reconeer -d example.com -silent | sort -u > subs_today.txt
# diff against yesterday and alert on new findings
```

Tips:
- Use `-rl` to set a conservative client-side rate limit.
- Provide an API key to isolate rate limits and support stable automation.
