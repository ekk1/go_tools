package main

const (
	clashYamlBase = `
---
port: 7890
socks-port: 10099
redir-port: 7892
allow-lan: true
mode: rule
log-level: info
ipv6: false
external-controller: 127.0.0.1:9090
dns:
  enable: true
  listen: 0.0.0.0:53
  default-nameserver:
    - 114.114.114.114
  enhanced-mode: fake-ip # or fake-ip
  fake-ip-range: 198.18.0.1/16 # Fake IP addresses pool CIDR
  nameserver:
    - 114.114.114.114 # default value

proxies:
  # socks5
  - name: "socks"
    type: socks5
    server: 127.0.0.1
    port: 10099
    udp: true
%s

proxy-groups:
%s

rules:
%s
`
)
