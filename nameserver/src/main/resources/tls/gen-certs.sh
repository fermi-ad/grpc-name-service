#!/bin/bash

openssl genrsa -out server.key 2048
openssl req -new -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj "/C=US/ST=Illinois/L=West Chicago/O=Fermilab/OU=ACORN/CN=Nameserver" -addext "subjectAltName = DNS:localhost,DNS:Nameserver,IP:127.0.0.1"
