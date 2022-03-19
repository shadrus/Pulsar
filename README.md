# **Pulstar is a tool to ship network checks to Prometheus.**


The pulsar is network testing tool. For now it can perfrom checks:


Test  | Description
--- | ---
Http | We can perform GET and POST requests.
Certificate | Check date and validity for site tls certificate.
Port | Check if post on remote resource is opened. (TODO)
Ping | Ping remote resource. (TODO)


All check results are available as Prometheus metrics on endpoint `/metrics`.


## Metrics
Name  | Description
--- | ---
pulsar_http_request_seconds | Time to get response from the endpoint in seconds.
pulsar_days_to_expire_cert | Days for site tls certificate expire.

All metrics also have label `success` that can be configured

## Configuration example

```
---
log_level: "debug"
certificate_config:
  - target:
      endpoint: "ya.ru"
      interval: 60
      timeout: 30
    days_for_warn: 3    
http_config:
  - target:
      endpoint: "https://ya.ru"
      interval: 60
      timeout: 30
    headers:
      Awasome-Header: something
    method: GET
    success_status: 200
    check_text: "yandex"
```

sections:

`certificate_config`: 

**days_to_warn** - success label will be set to false if days to certificate expire ale less than in config

**target.endpoint** - must be without protocol

`http_config`: 

**success_status** - if server has returned not this code, `succes` will be set to false

**check_text** - we can search string in response body to set `succes` flag







