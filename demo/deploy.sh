#!/bin/bash

echo "Creating certificates"
mkdir certs
openssl req -nodes -new -x509 -keyout certs/ca.key -out certs/ca.crt -subj "/CN=Admission Controller Demo"
openssl genrsa -out certs/admission-tls.key 2048
openssl req -new -key certs/admission-tls.key -subj "/CN=admission-server.default.svc" | openssl x509 -req -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/admission-tls.crt

echo "Creating k8s Secret"
kubectl create secret tls admission-tls \
    --cert "certs/admission-tls.crt" \
    --key "certs/admission-tls.key"

echo -n "Creating k8s admission deployment"

if test -f "secrets.yaml"; then
    "Creating registry secrets"
    kubectl apply -f secrets.yaml
fi

echo -n "Deploying certificates"
kubectl apply -f cert.yaml

if [ -x "$(command -v "ko")" ]; then
  echo -n "Found ko, building image from source"
  ko apply -f ../config/
else
  echo -n "Running with pre-built image"
  kubectl apply -f deployment.yaml
fi

echo -n "Creating k8s webhooks"
CA_BUNDLE=$(cat certs/ca.crt | base64 | tr -d '\n')
sed -e 's@${CA_BUNDLE}@'"$CA_BUNDLE"'@g' <"webhooks.yaml" | kubectl create -f -
