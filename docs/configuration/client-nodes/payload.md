```json
{
  "type": "payload",
  "disabled": false,
  "dialer_timeout": "2s",
  "handshake_retry": 2,
  "server_addr": "141.148.194.37:443",
  "payload_file": "http-payloads.dt",
  "speedtest_size": 1,
  "local_addr": "192.168.60.198",
  "net_iface": "Ethernet 4",
  "read_timeout": "4s",
  "write_timeout": "4s",
  "tls": {
    "enabled": true,
    "auth_timeout": "4s",
    "min_version": 1.0,
    "max_version": 1.3,
    "next_proto": [],
    "insecure": true
  }
}
```

### `type`

type එක `payload`

### `disabled`

true or false

false වුනොත් මේ Node එකෙන් check වෙන්නෙ නැ.

### `dialer_timeout`

tcp handshake timeout

### `handshake_retry`

handhsake එක retry කරන වාර ගණන

### `server_addr`

server එකෙ address එක

### `payload_file`

payload file path එක

### `speedtest_size`

speed test එකේදි download වෙන size එක

### `local_addr`

local interface ip address

### `net_iface`

local interface නම, interface ip address එක ලබාගන්න `local_addr` field එකට දීලනම් තියෙන්නෙ මේක අදාල නැහැ.

`local`

### `read_timeout`

Payload response read timeout

### `write_timeout`

Payload write timeout

### tls.enabled

true or false

### tls.auth_timeout

tls handshake timeout

### tls.min_version

tls minimum version 1.0

### tls.max_version

tls maximum version 1.3

### tls.next_proto

ALPN `h2`, `h3`, `http/1.1`

### tls.insecure

true or false
