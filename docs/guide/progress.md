# Progress

පලමු වර run කිරීමෙන් පසුව config හි `result.progress_file` අනුව file එකක් හැදිලා තියෙනව. මේකෙ තමා දැනට ටෙස්ට් වෙමින් පවතින Host file එකේ Progress එක තියෙන්නෙ.
මේකෙන් වෙන්නෙ Netshoot stop කරලා නැවත start කරහම කලින් නවත්තපු තැන ඉදලා ආපහු run වෙන එක.

( භාවිතයෙ පහසුවට මෙය prog.json විදිහට මෙතන් සිට සලකමු. )

පහල තියෙන්නෙ prog.json file එකට example එකක්

```json
{
  "checked_host": 2880,
  "total_tcp_fail": 10,
  "last_host": "example.com",
  "total_success": 20
}
```

ඉහත example එකට අනුව දැනට host 2880 ක් check වෙලා තියෙනවා. අවස්තා 20 success වෙලා තියෙනවා. (එකම හොස්ට් එක Payload කිහිපයකින් සාර්තක වෙන්නත් පුලුවන්.)

prog.json file එකේ `checked_host` & `total_tcp_fail` හැර ඉතිරි දෙක Netshoot වලින් භාවිතා කරන්නෙ නැ Information දෙන්න විහරයි තියෙන්නෙ.

ඒ වගේම `checked_host` එකේ ගාණ අනුව තමයි ඊලගට පටන් ගන්න ස්තානය තෝරන්නෙ. ඒ කියන්නෙ ඒක 0 කලොත් ආයෙ මුල ඉදන් තමයි යන්නෙ. 20 කලොත් 20 වෙනි host එකේ ඉදන් තමයි යන්නෙ. මේක වෙනස් කරන්න අවශ්‍ය වෙන්නෙ නැ ගොඩක් වෙලාවට.
