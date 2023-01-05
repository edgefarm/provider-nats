#!/bin/bash
set -e
kind create cluster --config kind-config.yaml --name provider-nats-e2e-test --image=kindest/node:v1.24.7
kubectl cluster-info --context kind-provider-nats-e2e-test