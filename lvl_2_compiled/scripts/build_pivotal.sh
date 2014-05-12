#!/bin/bash

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../../scripts && pwd )"

PRODUCT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"

RC_VALUE="-build${BUILD_NUMBER:--local}-precompiled"

source ${SCRIPTS_DIR}/shared.sh
