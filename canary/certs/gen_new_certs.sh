# gen ca
openssl genrsa 2048 > ca.key
openssl req -new -x509 -nodes -days 365000 -key ca.key -out ca.crt

# gen cert
openssl req -newkey rsa:2048 -nodes -days 3650 \
   -keyout client.key \
   -out client-req.crt
openssl x509 -req -days 365 -set_serial 01 \
   -in client-req.crt \
   -out client.crt \
   -CA ca.crt \
   -CAkey ca.key

# upload ca as secret
# kubectl create secret generic ca-cert --from-file=ca.crt=ca.crt -n tasks
