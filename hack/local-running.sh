#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


ROOT=$(dirname "${BASH_SOURCE[0]}")/..
HACK_ROOT=${ROOT}/hack

# build image
cd ${ROOT}
make image

# run
cd ${HACK_ROOT}
docker-compose up -d

