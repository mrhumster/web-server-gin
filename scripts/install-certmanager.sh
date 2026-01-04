#!/bin/bash
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.14.5 \
  --set installCRDs=true \
  --set image.repository=registry.cn-hangzhou.aliyuncs.com/google_containers/cert-manager-controller \
  --set webhook.image.repository=registry.cn-hangzhou.aliyuncs.com/google_containers/cert-manager-webhook \
  --set cainjector.image.repository=registry.cn-hangzhou.aliyuncs.com/google_containers/cert-manager-cainjector
