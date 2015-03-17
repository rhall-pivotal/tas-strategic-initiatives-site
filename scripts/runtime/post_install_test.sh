#!/bin/bash

set -e -x

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

bundle
${SCRIPTS_DIR}/runtime/set_bosh_deployment.sh "$RELENG_ENV"
$(bundle exec rake opsmgr:bosh:command[$RELENG_ENV]) run errand acceptance-tests
