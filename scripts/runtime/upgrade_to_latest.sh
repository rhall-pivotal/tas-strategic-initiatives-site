#!/bin/bash

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

bundle

source ${SCRIPTS_DIR}/runtime/download_pivotal_from_cache.sh
bundle exec rake --trace runtime:upgrade[${RELENG_ENV},${RUNTIME_PIVOTAL_FILE}]
