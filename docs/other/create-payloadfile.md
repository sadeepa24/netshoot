# Generate Payload File

Netshoot වෙතින් Inbuild tool එකක් දීලා තියෙනවා payload file එකක් හදාගන්න. `generate` command එකෙන් මේක Use කරන්න පුලුවන්.

[payload file එක මොකක්ද දන්නෙ නැත්තන් අනිවාරෙන් මේක කියවන්න](../guide/payload-file.md)

මේ payload file එක generate කරගන්න අවශ්‍ය Payload සහ අනෙකුත් දේ ලබාදෙන්න. file එකක් use කරනවා ඒකෙ විස්තර අනුව තමයි Payload file එක gen වෙන්නෙ. ඒක json format එකෙන් තමයි තියෙන්නෙ.

Example file :=

```json
{
  "format": "hex",
  "output": "test.pd",
  "payloads": [
    {
      "payload": "474554202f20485454502f312e310d0a486f73743a203c2d2d6e657473686f6f742d2d3e0d0a557067726164653a20776562736f636b65740d0a436f6e6e656374696f6e3a20557067726164650d0a5365632d576562536f636b65742d4b65793a2078334a4a484d62444c31457a4c6b68394742685844773d3d0d0a5365632d576562536f636b65742d56657273696f6e3a2031330d0a0d0a",
      "response": "485454502f312e312031303120537769746368696e672050726f746f636f6c730d0a557067726164653a20776562736f636b65740d0a436f6e6e656374696f6e3a20557067726164650d0a5365632d576562536f636b65742d4163636570743a2048536d726330734d6c59556b41476d6d354f50704732486147576b3d0d0a0d0a",
      "name": "test",
      "skip": false
    }
  ]
}
```

### `format`

payloads & response තියෙන format එක.

- hex
- raw

binary payload දෙනවනන් තමා hex format Use කරන්නෙ ගොඩක් වෙලවාට http payloads දෙන නිසා එව්වා ඔක්කොම raw format එකෙන් දෙන්න පුලුවන් Raw file එකකට Example 👇👇

```json
{
  "format": "raw",
  "output": "test.pd",
  "payloads": [
    {
      "payload": "GET /chat HTTP/1.1\r\nHost: <--netshoot-->\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==\r\nSec-WebSocket-Version: 13\r\n\r\n",
      "response": "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=\r\n\r\n",
      "name": "test"
    }
  ]
}
```

### `output`

output වෙන file එක. ( generate වෙන Payload file එකෙ නම)

### `payloads`

සියලුම payloads

## payload

```json
{
  "payload": "474554202f20485454502f312e310d0a486f73743a203c2d2d6e657473686f6f742d2d3e0d0a557067726164653a20776562736f636b65740d0a436f6e6e656374696f6e3a20557067726164650d0a5365632d576562536f636b65742d4b65793a2078334a4a484d62444c31457a4c6b68394742685844773d3d0d0a5365632d576562536f636b65742d56657273696f6e3a2031330d0a0d0a",
  "response": "485454502f312e312031303120537769746368696e672050726f746f636f6c730d0a557067726164653a20776562736f636b65740d0a436f6e6e656374696f6e3a20557067726164650d0a5365632d576562536f636b65742d4163636570743a2048536d726330734d6c59556b41476d6d354f50704732486147576b3d0d0a0d0a",
  "name": "test",
  "skip": false
}
```

### `payload`

payload එක

### `response`

payload එකට අදාල response එක

### `name`

payload එකේ නම

### `skip`

මේ value true වුනොත් gen වෙන payload file එකේ skip true වුන payload අන්තර්ගත වෙන්නෙ නැහැ.
