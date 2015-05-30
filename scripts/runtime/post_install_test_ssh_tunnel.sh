#!/bin/bash

set -e -x

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

pushd ${SCRIPTS_DIR}/..
  DIRECTOR_IP=`rake opsmgr:microbosh:director_ip[$RELENG_ENV]`
  OPSMGR_HOSTNAME=`rake opsmgr:info:host[${RELENG_ENV}]`
  PRIVATE_KEY=${SCRIPTS_DIR}/private.key

  rake opsmgr:get_private_key[${RELENG_ENV},${PRIVATE_KEY}]

  ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i ${PRIVATE_KEY} -f ubuntu@${OPSMGR_HOSTNAME} -L 25555:${DIRECTOR_IP}:25555 -N
  rm ${PRIVATE_KEY}

  export DIRECTOR_IP_OVERRIDE=localhost
  ${SCRIPTS_DIR}/runtime/post_install_test.sh
popd
