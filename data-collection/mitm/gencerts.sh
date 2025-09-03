openssl genrsa -out cert.key 2048
openssl req -new -key cert.key -out cert.csr -config san.cnf
openssl x509 -req -in cert.csr -signkey cert.key -out cert.crt -extensions v3_req -extfile san.cnf
cat cert.key cert.crt > cert.pem
