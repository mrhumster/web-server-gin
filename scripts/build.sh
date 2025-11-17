#!/bin/bash
VER=$1
if [ -z "$1" ]; then
  VER=$(git rev-parse HEAD)
fi
PROJECTDIR=/home/xomrkob/projects/web-server-gin/
echo "🔨 docker build..."
docker build \
  --build-arg VERSION=$VER \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t xomrkob/web-server-gin:$VER \
  -t xomrkob/web-server-gin:latest \
  $PROJECTDIR

if [ $? -eq 0 ]; then
  echo "🟢 Build success"
else
  echo "🔴 Build failed. Exited."
  exit 1
fi

echo "🚀 docker push..."
docker push xomrkob/web-server-gin:$VER

if [ $? -eq 0 ]; then
  echo "🟢 Push success"
else
  echo "🔴 Push to docker hub failed. Exited."
  exit 1
fi

echo "👷🏻‍♂️ k8s update image..."
kubectl -n go-app set image deployment/web-server-gin web-server-gin=xomrkob/web-server-gin:$VER

if [ $? -eq 0 ]; then
  echo "🟢 Image update  success"
else
  echo "🔴 K8s update failed. Exited."
  exit 1
fi
