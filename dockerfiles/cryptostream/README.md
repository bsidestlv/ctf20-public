curl -s --header Content-Type: application/json --data '{"message":"e56rduytfcgkvj BSidesTLV2020{Slip steaming all around} o789tuygkf"}' -k https://cryptostream.ctf.bsidestlv.com/encrypt
{"result":"2c155e13121059000e06471817074562300609071a3a2d225b5f5c10143d4c0409001f0102060c0e0b466941040d5604521b1d0b440e410252185a1b181b0e05077b6660612f60612f62762f637a6868","Signature":"d38f","StatusCode":200}

curl -s --header Content-Type: application/json --data '{"message":"2c155e13121059000e06471817074562300609071a3a2d225b5f5c10143d4c0409001f0102060c0e0b466941040d5604521b1d0b440e410252185a1b181b0e05077b6660612f60612f62762f637a6868","signature":"d38f"}' -k https://cryptostream.ctf.bsidestlv.com/decrypt
{"result":"e56rduytfcgkvj BSidesTLV2020{Slip steaming all around} o789tuygkf","StatusCode":200}

curl -s --header "Content-Type: application/json" --data '{"message":"what am i doing here?", "key":"1245678"}' -k http://localhost:8080/encrypt
{"result":"3e48091556044d540145441c080302000b0a1f07567d72677a7c7d337c7d337e6a337f6674747274","Signature":"8f63","StatusCode":200}

curl -s --header "Content-Type: application/json" --data '{"message":"what am i doing here?", "key":"1245678"}' -k https://cryptostream.ctf.bsidestlv.com/encrypt
{"result":"3e48091556044d540145441c080302000b0a1f07567d72677a7c7d337c7d337e6a337f6674747274","Signature":"8f63","StatusCode":200}
