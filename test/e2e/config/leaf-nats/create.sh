#!/bin/bash
set -e

kubectl apply -f foo.yaml
kubectl apply -f bar.yaml
kubectl rollout status deployment leaf-nats-foo -n nats --timeout=90s
kubectl rollout status deployment leaf-nats-bar -n nats --timeout=90s
