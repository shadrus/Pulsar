# **Pulsar is a tool to ship network checks to Prometheus.**


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
pulsar_cert_not_after | The date after which a peer certificate expires. Expressed as a Unix Epoch Time.
pulsar_cert_not_before | The date before which a peer certificate is not valid. Expressed as a Unix Epoch Time.

## Configuration

Tool can be configured through config.yml file. It must be placed at the same dir with app.

### Configuration example

```
---
log_level: "debug"
certificate_config:
  - target:
      endpoint: "ya.ru"
      interval: 60
      timeout: 30
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

**target.endpoint** - must be without protocol

`http_config`: 

**success_status** - if server has returned not this code, `succes` will be set to false

**check_text** - we can search string in response body to set `succes` flag







