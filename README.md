# HawkWing
The wings which brings your packets to the destination

## Configuration
```yaml
---
hawkeye:
  hostname: hawkeye.hawk.net
  port: 5001
services:
  service1:
    domain_name: service1.hawk.net
    applications:
      - port: 80
        intents:
          - intent: sfc
            functions:
              - function1
              - function2
    sid:
      - fcbb:bb00:1::2
      - fcbb:bb00:3::2
  service2:
    domain_name: service2.hawk.net
    applications:
      - port: 80
        intents:
          - intent: low-latency
            min_value: 10
            max_value: 20
```

## Development Network
![Development Network](docs/network.png)
