# Configurations

Full Config Structure

```json
{
  "client": {},
  "server": {},
  "result": {},
  "host": {},
  "log": {}
}
```

| Field                   | Info               |
| ----------------------- | ------------------ |
| [`client`](./client.md) | Client Side        |
| [`server`](./server.md) | Server Side        |
| [`result`](#result)     | Output Result Info |
| [`host`](#host)         | Host File          |
| [`log`](#log)           | Logger             |

## Result

```json
{
  "output_file": "out.json",
  "progress_file": "prog.json",
  "tcp_fail_threshold": 1000
}
```

### `output_file`

test සියලුම output save විය යුතු file එක, json format එකට තමා save වෙන්නෙ.

### `progress_file`

progress information save විය යුතු file එක, මේකත් json

### `tcp_fail_threshold`

max tcp fail වෙන්න පුලුවන් ගණන, මෙ ගානට වඩා fail වුනු ගමන් Netshoot Auto ම stop වෙනවා.

## Host

```json
{
  "max_concurrent": 15,
  "host_file": "host.txt",
  "interval": "10ms"
}
```

### `max_concurrent`

එකවර test වෙන host ගනන්, Network capacity එක අනුව මේක දෙන්න.

### `host_file`

test කරන host file එක

### `interval`

host set එකක් check කරලා ඊලග සෙට් එක චෙක් වෙන්න කලින් තියෙන time එක ( 0 දෙන්න)

## log

```json
{
  "level": "debug",
  "paths": ["stdout"],
  "encode": "console"
}
```

### `level`

log level එක

- debug
- info
- warn
- error

### `paths`

log output paths

### `encode`

encode methods

- console
- json
