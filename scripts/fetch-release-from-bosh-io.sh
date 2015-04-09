set -e

SCRIPT_DIR=$( cd "$( dirname $0 )" && pwd)
RELEASE_NAME=$1
BOSH_IO_RELEASE_NAME=$2

RELEASE_DEFINITON=$(grep -B1 -A3 "^  name: ${RELEASE_NAME}$" ${SCRIPT_DIR}/../metadata_parts/binaries.yml)
RELEASE_VERSION=$(echo "${RELEASE_DEFINITON}" | grep version | grep -o -E [0-9]+)
RELEASE_FILE=$(echo "${RELEASE_DEFINITON}" | grep file | grep -o -E "\S+$")

if [ -e ./releases/${RELEASE_FILE} ]
  then
    echo "${RELEASE_FILE} already exists"
    exit 0
fi

set -x
${SCRIPT_DIR}/run_in_docker.sh aria2c -x 5 --out=releases/${RELEASE_FILE} http://bosh.io/d/github.com/cloudfoundry/${BOSH_IO_RELEASE_NAME}?v=${RELEASE_VERSION}
