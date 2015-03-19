#!/bin/bash

set -e

if [ -z "${RUNTIME_DOT_PIVOTAL_FILE}" ]; then
  echo 'You must set the RUNTIME_DOT_PIVOTAL_FILE environment variable!' >&2
  exit 1
fi

PIVOTAL_FILE_CACHE_KEY="runtime/${RUNTIME_DOT_PIVOTAL_FILE}"

bundle
bundle exec rake artifacts:cf:clean[.]
bundle exec rake artifacts:cf:retrieve[${PIVOTAL_FILE_CACHE_KEY}]
bundle exec rake md5:validate[${RUNTIME_DOT_PIVOTAL_FILE},${RUNTIME_DOT_PIVOTAL_FILE}.md5]
bundle exec rake --trace runtime[${RELENG_ENV},${RUNTIME_DOT_PIVOTAL_FILE}]
