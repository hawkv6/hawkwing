# HawkWing
The wings which brings your packets to the destination

```yaml
---
hawkeye:
  hostname: hawkeye.hawk.net
  port: 1234
services:
  wb.hawk.net:
    - intent: high-bandwidth
      port: 80
      sid: 
        - fcbb:bb00:1::2
        - fcbb:bb00:3::2
    - intent: low-bandwidth
      port: 8080
  wc.hawk.net:
    - intent: high-bandwidth
      port: 1433
      sid: 
        - fcbb:bb00:2::2
        - fcbb:bb00:4::2
```

```
bpftool net show
bpftool prog tracelog
trace add virtio-input 10
mount --make-shared /sys/fs/bpf
```

