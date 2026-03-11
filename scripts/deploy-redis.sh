#!/bin/bash
helm install casbin-redis oci://registry-1.docker.io/bitnamicharts/redis --namespace go-app --set architecture=standalone --set auth.enabled=true --set auth.password=password --set master.persistence.enabled=false
