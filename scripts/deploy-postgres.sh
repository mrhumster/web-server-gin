#!/bin/bash
helm install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql \
  -n go-app \
  -f /home/xomrkob/projects/web-server-gin/k8s/base/values.yaml
