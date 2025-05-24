# PAYLOAD FILE

## **_✅ ක්‍රීයා කරන අයුරු_**

- client node එකට දීලා තියෙන payload file එකේ සියලුම payload භාවිතා කරලා එකම Host එක test වෙනවා.

       - Example - node එකට අදාලව තියෙන payload file එකේ payload 5 ක් තියේනම්, එකම Host එක පස් වතාවක් ඒ ඒ payload එකට අනුව test වෙනවා.

- Payload එකේ `<--netshoot-->` parameter එක තියෙන තැන් test වෙන Host එකෙන් replace වෙනවා. ඊට පස්සෙ තමා payload එක send වෙන්නෙ.

````md
```
GET / HTTP/1.1\r\n
Host: <--netshoot-->\r\n
User-Agent: netshoot/1.0\r\n
Accept: */*\r\n
Connection: close\r\n
\r\n
```

example.com කියන හොස්ට් එක ඉහත payload එක යොදාගෙන test කරොත් send වෙන payload එක වෙන්නෙ 👇

```
GET / HTTP/1.1\r\n
Host: example.com\r\n
User-Agent: netshoot/1.0\r\n
Accept: */*\r\n
Connection: close\r\n
\r\n
```

ඉහත Payload එකෙ `<--netshoot-->` කියන එක example.com ලෙස වෙනස් වෙලා තියෙනවා.
````

````md
ඒ වගෙම `<--netshoot-->` මේ Parameter එක Payload එකෙ ඕනම තැනක භාවිතා කරන්න පුලුවන්. (මුල සහ අග හැර) 👇

```
GET / HTTP/1.1\r\n
<--netshoot--> Host: <--netshoot-->\r\n
User-Agent: netshoot/1.0\r\n
Accept: */*\r\n <--netshoot--> <--netshoot-->
<--netshoot--> Connection: close\r\n
\r\n
```

example.com මගින් replace වීමෙන් පසුව 👇

```
GET / HTTP/1.1\r\n
example.com Host: example.com\r\n
User-Agent: netshoot/1.0\r\n
Accept: */*\r\n example.com example.com
example.com Connection: close\r\n
\r\n
```
````

## **_✅ ප්‍රධාන වශයෙන් සැලකිය යුතු කරුනු_**

- test කරන Payload එක server & client file දෙකේම තියෙන්න ඕන.
- payload එකේ මුලට සහ අගට Parameter දාන්න බැහැ.
- Payload එකට අදාලව අනිවාරෙ response payload එකකුත් තියෙන්න ඕනෙ. (පහල Header එක බලන්න)
- Payload file structure එක පහල විදියටම තියෙන්න ඕන.
- Server & Client දෙකටම එකම Payload FIle structure එක භාවිතා කරන්නෙ ( පහල බලන්න )
- Server එකෙයි client ගෙයි එකම payload file තියෙන්න ඕනෙ නැහැ. නමුත් client test කරන Payload සේරම server එකේ payload file එකේ තියෙන්න ඕනෙ.
- Binary payload එකක් වුනත් use කරන්න පුලුවන්.

## 📝 **_ භාවිතයට උපදෙස් _**

- එකවර payload file කිහිපයක් වුනත් client or server side භාවිතා කරන්න පුලුවන්

      - Example - ඔයාගෙ Client side එකේ Node තුනක් හදලා තියෙනවනන් ඒ තුනෙ එකිනෙකට වෙනස් payload file තුනක් භාවිතා කරන්න පුලුවන්. ( හැබැයි ඒ තුනම එකම server node එකටනම් connect වෙන්නෙ server node එකට දීලා තියෙන Payload file එකේ අර client side තියෙන payload file තුනේම තියෙන payload තියෙන්න ඕනෙ.)

- Connect වෙන server Node එකට දීලා තියෙ payload file එකේ ඕනෙ ගනනක් payload තියෙන්න පුලුවන්. client ටෙස්ට් කරන්න පුලුවන් වෙන්නෙත් ඒ payload ටිකමයි

      - Example - හිතන්න server node එකකට දීලා තියෙනවා 100 payloads තියෙන file එකක්. client side එකෙන් මේ node එක Use කරලා Host check කරනවනන් Use කරන්න පුලුවන් වෙන්නෙ අර payload 100 විතරයි. client node එකෙ Payload ගනන 2 ක් වුනත් කමක් නැ. හැබැයි server node එකට දීලා තියන payload file එකේ අදාල payload දෙක තියෙන්න ඕනෙ.

- Server node එකෙයි Client node එකෙයි payload file වල Payload ගනන කිසි විටක සමාන වියයුතු නැහැ. හැබැයි අනිවාරෙන් test වෙන payload node දෙකේම තියෙන්න ඕනෙ.

      - Example - server side file එකේ payload 20 ක් තියෙනවා.client side file එකේ payload 5 ක් තියෙනවා. මෙහෙම තිබ්බට ප්‍රශ්නයක් නැහැ. හැබැයි client side තියෙන payload 5 අනිවාරෙන් server side 20 ඇතුලෙ තියෙන්න ඕන.

## **_📝 Payload File එකක් සකසන අයුරු_**

### 📦 Payload File Header

| 2 bytes       | 1 byte      | X bytes      | 4 bytes        | 4 bytes         | X bytes | X bytes  |
| ------------- | ----------- | ------------ | -------------- | --------------- | ------- | -------- |
| Payload Count | Name Length | Payload Name | Payload Length | Response Length | Payload | Response |

- Payload File එකක් ඔයා manual හදාගන්නවනන් උඩ තියෙන structure එකට අනුව තමා හදන්න ඕනෙ. **අනිවාරෙන්ම සියලුම Length Bigendian Format එකට තමා තියෙන්නෙ ඕනෙ.**
- netshoot වෙතින් inbuild tool එකක් දීලා තියෙනවා මේ payload file හදාගන්න. [(වැඩි විස්තර)](../other/create-payloadfile.md)
