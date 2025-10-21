#!/bin/bash

set -e

echo "Generating RSA keys..."

mkdir -p config/keys

cd config/keys

echo "Generating private key..."
openssl genrsa -out private.pem 2048

echo "Generating public key..."
openssl rsa -in private.pem -pubout -out public.pem
chmod 600 private.pem
chmod 644 public.pem

echo "✅ RSA keys generated successfully!"
echo "📁 Private key: config/keys/private.pem"
echo "📁 Public key:  config/keys/public.pem"
echo ""
echo "⚠️  Keep private.pem secure! Do not commit to repository!"
