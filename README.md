# payPDL-api

- simpan session dengan menggunakan :mksession
- buka session pake nvim -S payPDL-api
- untuk start, restart, stop service menggunakan "sudo systemctl start/restart/stop payPDL-api.service"
- untuk melihat logs pakai "journalctl -u payPDL-api.service"
- pakai curl :
  curl -X POST -H "Content-type: application/json" \                                                                                                                       ✔  17s 
-d '{"Nop":"3512aaabbbcccdddde",
        "Masapajak":"20xx",
        "DateTime":"2020-08-25 14:11:05",
        "Merchant":"6010",
        "KodeInstitusi":"001011",
        "NoHp":"08123456789",
        "Email":"a@a.com"}' \
localhost/inquiry | jq

