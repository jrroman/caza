#!/usr/bin/env bash

DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

NAME="caza"

# Check if the namespace exists, if it does not create it
kubectl get namespace "$NAME" || \
    kubectl create namespace "$NAME"

# Apply deployment manifest
kubectl apply -n "$NAME" -f <( \
    NAME="${NAME}" \
    envsubst '$NAME' < \
    "$DIR/deployment.yaml"
)
