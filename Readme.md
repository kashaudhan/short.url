# Url shortener using Go, gin & redis

It uses `Base62` encoding.


Post request to add a url to redis db with custom url

```bash
curl --location 'http://localhost:8080/add' \
--header 'Content-Type: application/json' \
--data '{
    "url": "https://kashaudhan.in",
    "short": "kashaudhan"
}'
```

response
```bash
{
    "custom_short": "/kashaudhan",
    "expiry": 24,
    "url": "https://kashaudhan.in",
    "x-rate-limiting": -2,
    "x-rate-limiting-reset": 0
}
```

Get request to redirect to the url

```bash
curl --location 'http://localhost:8080/kashaudhan'
```
