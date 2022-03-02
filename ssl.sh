#/bin/bash

SERVER_CN=localhost

# Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -out ssl/ca.key 4096
openssl req -new -x509 -days 365 -key ssl/ca.key -out ssl/ca.crt -subj "/CN=${SERVER_CN}"

# Generate server private key
openssl genrsa -out ssl/server.key 4096

# Get Certificate Signing Request (CSR) from the CA (server.csr)
openssl req -new -key ssl/server.key -out ssl/server.csr -subj "/CN=${SERVER_CN}"

# Sign the certificate with the CA we created (server.crt)
# legacy Common Name field is ignored.
# https://pkg.go.dev/crypto/x509#Certificate.VerifyHostname
echo "subjectAltName=DNS:${SERVER_CN}" > altsubj.ext
openssl x509 -extfile altsubj.ext -req -days 365 -in ssl/server.csr -CA ssl/ca.crt -CAkey ssl/ca.key -set_serial 01 -out ssl/server.crt

# Convert the server certificate to PEM format
openssl pkcs8 -topk8 -nocrypt -in ssl/server.key -out ssl/server.pem
