#!/usr/bin/bash

curl -X POST \
  -H "Authorization: Bearer ${RDEPOT_TOKEN}" \
  -H 'Content-Type: multipart/form-data' \
  -H 'Accept: application/json' \
  -F 'file=@test/resources/oaColors_0.0.4.tar.gz;type=application/gzip' \
  -F 'repository=public' \
  -F 'replace=false' \
  localhost:8080/api/manager/packages/submit

