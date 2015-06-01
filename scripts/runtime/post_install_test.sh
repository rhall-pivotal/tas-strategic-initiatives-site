#!/bin/bash

set -e -x

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"
if [ "x${INTERNETLESS}" != "x" ]
then
  errand_name=acceptance-tests-internetless
else
  errand_name=acceptance-tests
fi

${SCRIPTS_DIR}/runtime/set_bosh_deployment.sh "$RELENG_ENV"
$(bundle exec rake opsmgr:microbosh:command[$RELENG_ENV]) run errand $errand_name
