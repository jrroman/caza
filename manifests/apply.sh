#!/usr/bin/env bash

# DIR's value is the directory we are currently inside of
DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# Name of the application
NAME="caza"
NAMESPACE="$NAME"

# Check if the namespace exists, if it does not create it
kubectl get namespace "$NAMESPACE" || \
    kubectl create namespace "$NAMESPACE"

# Apply deployment manifest
kubectl apply -n "$NAMESPACE" -f <( \
    NAME="${NAME}" \
    envsubst '$NAME' < \
    "$DIR/deployment.yaml"
)

# Set region if it is passed in, otherwise default to "us-east-1"
# e.g. REGION=us-west-2 ./apply.sh
REGION="${REGION:-us-east-1}"

# Apply Configmap
kubectl apply -n "$NAMESPACE" -f <( \
    NAME="${NAME}" \
    REGION="${REGION}" \
    envsubst '$NAME $REGION' < \
    "$DIR/config.yaml"
)

# Apply Service Account
# the role's json policy is located at ./aws/iam-policy.json
IAM_ROLE_ARN="placeholder" # this is temporary
kubectl apply -n "$NAMESPACE" -f <( \
    NAME="${NAME}" \
    ROLE_ARN="${IAM_ROLE_ARN}" \
    envsubst '$NAME $ROLE_ARN' < \
    "$DIR/service-account.yaml"
)
