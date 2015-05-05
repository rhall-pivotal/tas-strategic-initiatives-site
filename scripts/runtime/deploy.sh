#!/bin/bash

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

source ${SCRIPTS_DIR}/runtime/download_pivotal_from_cache.sh


bundle exec rake --trace ert:upload[${RELENG_ENV},1.4,${RUNTIME_PIVOTAL_FILE}]
bundle exec rake --trace ert:configure[${RELENG_ENV},1.5,1.4]
bundle exec rake --trace opsmgr:product:install[${RELENG_ENV}]
