#!/bin/bash
helm -n go-app install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql -f ~/projects/web-server-gin/deploy/k8s/base/values.yaml
