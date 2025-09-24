#!/bin/bash
helm install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql \
  -n go-app \
  -f values.yaml
