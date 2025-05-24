# Payload

දැනට තියෙන්නෙ මේ type (`payload`) එකේ nodes විතරයි.

## 🔰 **ක්‍රියා කරන අයුරු**

Node එකකින් තමයි host එක චෙක් වෙන method එක define වෙන්නෙ. එකේත් එක විශේශ විදියක් තමයි payload node කියන්නෙ. මේ Node එක්කින් ඔයාට කැමති Method එකකට Host එක check වෙන්න දෙන්න පුලුවන්. ඒ වගේම payload Node එකකට (client /server දෙකටම) අනිවාරෙන් Payload file එකක් අවශ්‍ය වෙනවා.

payload file එකේ තියෙන payload තමයි client - server අතර send recive වෙන්නෙ ([වැඩි විස්තර](./payload-file.md))

මේ payload මගින් ඔයාට ඕනම Method එකක් හදාගන්න පුලුවන්

- උදා - v2ray/ws වැඩ කරනවද කියලා බලාගන්න විදියට payload එකක් හදගන්න අවශ්‍යනම්. v2ray/ws වැඩ කරන විදිය පොඩ්ඩක් කියන්නම් v2ray ws connection එකක් establish වෙද්දි මුලින්ම http GET req එකක් තමයි යවන්නෙ. ඊට පස්සෙ තමයි websocket වලට switch වෙන්නෙ. අපිට අවශ්‍ය වෙන්නෙ connection එක establish වෙනවද බලාගන්න විතරනෙ. ඒ හින්දා http GET req එකට අදාල Payload & Response එක තිබ්බාම ඇති අදාල payload එකෙන් host එකක් වැඩ කලොත් ඒක අනිවාරෙන් v2ray/ws වලට වැඩ කරනව. පහල payload එකයි response එකයි මේ සඳහා use කරන්න පුලුවන්.

```
GET /chat HTTP/1.1\r\n
Host: <--netshoot-->\r\n
Upgrade: websocket\r\n
Connection: Upgrade\r\n
Sec-WebSocket-Key: x3JJHMbDL1EzLkh9GBhXDw==\r\n
Sec-WebSocket-Version: 13\r\n
\r\n
```

```
HTTP/1.1 101 Switching Protocols\r\n
Upgrade: websocket\r\n
Connection: Upgrade\r\n
Sec-WebSocket-Accept: HSmrc0sMlYUkAGmm5OPpG2HaGWk=\r\n
\r\n
```

server node tls සමඟ listen කලොත් ඒවා tls නොමැති connection වලටත් support කරනවා.

- **උදා server node එක port 443 tls සමඟ listen කලොත්, client side එකෙන් tls disable payload or tls enable payload දෙකම test කරන්න පුලුවන්.**

note - tls සමඟ test කරනවනන් payload කිහිපයක් test කිරීම කිසිම තේරුමක් නැ. එකක් පමණක් tls වලින් check කරගන්න. (tls දැම්ම ගමන් isp ට අදාල වෙන්නෙ client hello & server hello මේ වගේ Tls authnicate වෙන්න share වෙන details නැතුව Payload එක එයාලට පේන්නෙ නැ. ඔයා payload මොක දැම්මත් client hello, server hello මේ දේවල් එකමයි.)

## 🔰 විශේශ කරුණු

- client node එකට අදාලව ගැලපෙන server node එකක් තිබිය යුතුයි.
- client / server payload file වල payload ගැලපිය යුතුයි. ( [වැඩි විස්තර](./payload-file.md))
- server tls node මගින් tls නොමැති connection වුවද accept කරයි. ( tls disable කරලා වුනත් test කරන්න පුලුවන්.)

## 🔰 Timeouts

මේ Node එකේ ප්‍රදාන වශයෙන් timeout කිහිපයක් භාවිතා කරනවා.

### client timeouts

- dialer_timeout - tcp handshake එක වෙන්න ලබාදෙන උපරිම වෙලාව.
- read_timeout - payload එක read කරන්න ලබාදෙම උපරිම වෙලාව
- write_timeout - payload එක write කරන්න ලබාදෙන උපරිම වෙලාව
- tls.auth_timeout - tls enable කරලනම් තියෙන්නෙ tls handshake එක වෙන්න ලබාදෙන උපරිම වෙලාව

### server timeouts

- read_timeout - payload එක read කරන්න ලබාදෙම උපරිම වෙලාව.
- write_timeout - payload එක write කරන්න ලබාදෙන උපරිම වෙලාව.
- tls.timeout - tls enable කරලනම් තියෙන්නෙ tls handshake එක වෙන්න ලබාදෙන උපරිම වෙලාව

මෙම timeouts තමන්ගෙ network capacity එකේ හැටියට ලබා දෙන්න ඕන. ඒ වගේම `max_concurrent` ගනන් වැඩි වෙන්න වැඩි වෙන්න timeouts ත් පොඩ්ඩක් වැඩි කරන්න වෙනවා. fiber වගේ connection එකක් තියෙන කෙනෙක්ට ගොඩක් අඩු timeouts use කරන්න පුලුවන්.
