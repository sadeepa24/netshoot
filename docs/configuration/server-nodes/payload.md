```json
{
  "type": "payload",
  "disabled": false,
  "payload_file": "http-payloads.dt",
  "read_timeout": "1s",
  "write_timeout": "1s",
  "listen_conf": {
    "listen": "0.0.0.0:443",
    "tls": {
      "enabled": true,
      "cert": "./ssl/certificate.pem",
      "key": "./ssl/private_key.pem",
      "timeout": "900ms"
    }
  }
}
```

### type

type eka `payload`

### disabled

true or false

false වුනොත් මේ Node එක on වෙන්නෙ නැ.

### payload_file

payload file path එක

### read_timeout

connection එකක් establish වුනාම payload එක read කරන්න ගන්න පුලුවන් උපරිම වෙලාව

### write_timeout

response write කරන්න ගන්න පුලුවන් උපරිම වෙලාව.

### listen_conf.listen

listen වෙන socket එක

### listen_conf.tls.enabled

true or false

### listen_conf.tls.cert

tls certficate Path එක

### listen_conf.tls.key

tls privkey Path එක

### listen_conf.tls.timeout

tls handhsake timeout
