#!/bin/bash

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

source ${SCRIPTS_DIR}/runtime/download_pivotal_from_cache.sh

if [[ $NO_INSTALL = true ]]; then
  bundle exec rake --trace runtime:setup[${RELENG_ENV},${RUNTIME_PIVOTAL_FILE}]
else
  bundle exec rake --trace runtime[${RELENG_ENV},${RUNTIME_PIVOTAL_FILE}]
fi
