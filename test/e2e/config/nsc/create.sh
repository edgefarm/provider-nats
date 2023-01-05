#!/bin/bash
set -e

echo "Pushing nats user/account credentials"
FREE_RANDOM_PORT=$(comm -23 <(seq 49152 65535 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1)
f=$(mktemp)
kubectl port-forward -n nats svc/nats ${FREE_RANDOM_PORT}:4222 &
export pid=$!
echo "Waiting for port-forward to be ready"
until [ $(netstat -l | grep ${FREE_RANDOM_PORT} | wc -l) -gt 0 ]; do sleep 1; echo "."; done
echo " done"
nsc push --system-account SYS -A -u nats://localhost:${FREE_RANDOM_PORT} --config-dir nsc/accounts --data-dir nsc/accounts/nats --keystore-dir nsc/nkeys
kill ${pid}
kubectl create -n nats secret generic leaf-nats-foo-creds --from-file=nsc/nkeys/creds/myoperator/acc1/user.creds --dry-run=client -o yaml | kubectl apply -f -
kubectl create -n nats secret generic leaf-nats-bar-creds --from-file=nsc/nkeys/creds/myoperator/acc2/user.creds --dry-run=client -o yaml | kubectl apply -f -