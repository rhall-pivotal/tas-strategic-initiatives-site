#!/bin/bash
set -xe

if [ $# -ne 1 ]; then
  echo 'Usage: $0 env_name'
  exit 1
fi

env=$1

openssl genrsa -out $env.pem 2048
openssl req -sha256 -new -key $env.pem -out $env-csr.pem \
  -subj "/C=US/ST=California/L=San\ Francisco/O=Pivotal/OU=Release\ Engineering/CN=*.$env.cf-app.com/emailAddress=cf-release-engineering@pivotal.io"

openssl x509 -req -days 365 -in $env-csr.pem -signkey $env.pem -out $env-cert.pem

read -p "hit enter to upload $env SSL certificate to AWS $1 "

aws iam upload-server-certificate --server-certificate-name $env \
--certificate-body file://$env-cert.pem --private-key file://$env.pem 
