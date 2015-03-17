#!/bin/bash

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

bundle

source ${SCRIPTS_DIR}/runtime/get_published_pivotal.sh
bundle exec rake --trace runtime[${RELENG_ENV},${RUNTIME_DOT_PIVOTAL_FILE}]
