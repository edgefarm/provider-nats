#!/bin/bash
set -e

echo "Pushing nats user/account credentials"
if [ -f nsc/accounts/nats/nsc.json ]; then
    rm nsc/accounts/nats/nsc.json
fi
tar cfz nsc.tar.gz nsc
base64 nsc.tar.gz > nsc.tar.gz.base64
kubectl -n nats create secret generic nsc --from-file=nsc.tar.gz.base64 --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f manifest.yaml
kubectl wait --for=condition=Complete job -l app=nsc -n nats --timeout=60s
kubectl create -n nats secret generic leaf-nats-foo-creds --from-file=nsc/nkeys/creds/myoperator/acc1/user.creds --dry-run=client -o yaml | kubectl apply -f -
kubectl create -n nats secret generic leaf-nats-bar-creds --from-file=nsc/nkeys/creds/myoperator/acc2/user.creds --dry-run=client -o yaml | kubectl apply -f -