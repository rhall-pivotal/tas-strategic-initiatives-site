#!/bin/bash

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../../scripts && pwd )"

PRODUCT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"

rm -f ${PRODUCT_DIR}/*.pivotal
rm -f ${PRODUCT_DIR}/*.pivotal.yml
rm -f ${PRODUCT_DIR}/*.pivotal.md5

RC_VALUE="-build${BUILD_NUMBER:--local}-precompiled"

source ${SCRIPTS_DIR}/shared.sh
