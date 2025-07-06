#!/bin/bash

# 定义Header
HEADER='{"alg":"HS256","typ":"JWT"}'

# 定义Payload
PAYLOAD='{"exp":1739078005,"iat":1735478005,"nbf":1735478005,"x-user-id":"user-w6irkg"}'

# 定义Secret
SECRET="Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5"

# 1.Base64编码Header
HEADER_BASE64=$(echo -n "${HEADER}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 2.Base64编码Payload
PAYLOAD_BASE64=$(echo -n "${PAYLOAD}" | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

# 3.拼接Header和PayLoad为签名数据
SIGNING_INPUT="${HEADER_BASE64}.${PAYLOAD_BASE64}"

# 4.使用HMAC SHA256算法生成签名
SIGNATURE=$(echo -n "${SIGNING_INPUT}" | openssl dgst -sha256 -hmac "${SECRET}" -binary | openssl base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n' )

# 5.拼接最终的JWT Token
JWT="${SIGNING_INPUT}.${SIGNATURE}"

echo "Generated JWT Token:"
echo "${JWT}"