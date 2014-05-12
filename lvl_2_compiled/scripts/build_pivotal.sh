#!/bin/bash

set -e

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../../scripts && pwd )"

PRODUCT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"

source ${SCRIPTS_DIR}/shared.sh
