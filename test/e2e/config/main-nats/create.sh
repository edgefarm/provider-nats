#!/bin/bash
set -e
kubectl create namespace nats
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm upgrade --install nats nats/nats --namespace nats --values helm/values.yaml
until kubectl rollout status --watch --timeout=60s statefulset -n nats nats 2>/dev/null; do echo -n "." && sleep 2; done
