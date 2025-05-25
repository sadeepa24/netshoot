# Output

```json
"result": {
"output_file": "out.json",
"progress_file": "prog.json",
"tcp_fail_threshold": 1000
},
```

පලමු වර Netshoot run කිරීමෙන් පසුව result.progress_file සඳහා දී ඇති නමට අනුව json file එකක් හැදිලා තියෙවි. ඒ file එකේ තමයි සියලුම host ව්ලාට අදාල test details තියෙන්නෙ.

පහල output file එකක example එකක් පහල තියෙනවා.

## Example file

```json
[
  {
    "Host": "www.facebook.com",
    "TotalTcpFail": 0,
    "Success": 2,
    "Maybe": 4,
    "MaxSpeed": 0.919445,
    "Err": "",
    "PayloadInfo": [
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-GET",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "www.facebook.com"
        },
        "MaxSpeed": 0.788401
      },
      {
        "Success": false,
        "Maybe": false,
        "Error": "",
        "PayloadName": "raw-http-HEAD",
        "TpcFailed": false,
        "Tls": {
          "Failed": true,
          "Error": "EOF",
          "Servername": "www.facebook.com"
        },
        "MaxSpeed": 0.0
      },
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-PUT",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "www.facebook.com"
        },
        "MaxSpeed": 0.919445
      },
      {
        "Success": false,
        "Maybe": true,
        "Error": "speedtest err: EOF",
        "PayloadName": "raw-http-POST",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "www.facebook.com"
        },
        "MaxSpeed": 0.0
      },
      {
        "Success": false,
        "Maybe": true,
        "Error": "speedtest err: EOF",
        "PayloadName": "raw-http-PATCH",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "www.facebook.com"
        },
        "MaxSpeed": 0.0
      }
    ]
  },
  {
    "Host": "m.facebook.com",
    "TotalTcpFail": 0,
    "Success": 4,
    "Maybe": 5,
    "MaxSpeed": 2.153463,
    "Err": "",
    "PayloadInfo": [
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-GET",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "m.facebook.com"
        },
        "MaxSpeed": 0.835495
      },
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-HEAD",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "m.facebook.com"
        },
        "MaxSpeed": 1.548489
      },
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-PUT",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "m.facebook.com"
        },
        "MaxSpeed": 2.153463
      },
      {
        "Success": true,
        "Maybe": true,
        "Error": "",
        "PayloadName": "raw-http-POST",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "m.facebook.com"
        },
        "MaxSpeed": 1.668666
      },
      {
        "Success": false,
        "Maybe": true,
        "Error": "speedtest err: EOF",
        "PayloadName": "raw-http-PATCH",
        "TpcFailed": false,
        "Tls": {
          "Failed": false,
          "Error": "",
          "Servername": "m.facebook.com"
        },
        "MaxSpeed": 0.0
      }
    ]
  }
]
```

## How to Read

ඉහත file එක සම්පූර්ණ Array එකක් කියලා බලාගන්න පුලුවන්. මේකෙ එක item එකක් අදාල වෙන්නෙ එක Node එහෙකට ඒ කියන්නෙ එක host එකක් check වෙද්දි ඒ host එක අදාල Node එකෙන් check වුන විදිය තමයි තියෙන්නෙ. Node කිහිපයක් තියෙනම් එකම Host එක ඒ node ඕක්කොගෙන්ම check වෙනවා.

ඒ වගේම එක item එකක් ඇතුලෙ PayloadInfo කියලා තව array එකක් තියෙනවා. මේ තියෙන්නෙ අර Payload file එකෙ තිබුන එ ඒ Payload වලට අදාලව result එක, payload එකට අදාල නමෙන් payload එක හොයාගන්න පුලුවන්.

### Extract Info

මේකෙන් information extract කරගන්න අවශ්‍යනම් ඒ කියන්නෙ වැඩ කරන host file ටික වෙනම වගේ එකවර ගන්න ඕනනන්ම්, json path වගෙ use කරන්න පුලුවන්.
online [JSONPATH](https://jsonpath.com/) වගේ tool එකක් use කරන්න පුලුවන්.

- Example එකක් විදිහට
  මේකට output file එක upload කරලා json path එකට `$[?(@.PayloadInfo[?(@.MaxSpeed > 0)])].Host` දැම්මොත් speedtest එක 0 ට වැඩියෙන් ආපු සියලුම Host list එක ගන්න පුලුවන්.

### More Jsonpath example to extract Info

| JsonPath                                             | Info                                               |
| ---------------------------------------------------- | -------------------------------------------------- |
| `$[?(@.PayloadInfo[?(@.MaxSpeed > 0)])].Host`        | එක Payload එකකින් හෝ speed එක 0 ට වැඩි සියලුම host |
| `$[?(@.PayloadInfo[?(@.MaxSpeed > 0)])].PayloadInfo` | speedtest එක 0 ට වැඩි සියලුම Host Payload Info සමඟ |
