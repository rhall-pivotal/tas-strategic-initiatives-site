#!/bin/bash

set -ex

if [ $# -ne 4 ]; then
  echo "Usage: $0 IAAS STEMCELL_VERSION ENVIRONMENT OM_VERSION" >&2
  exit 1
fi

iaas=$1
stemcell_version=$2
environment=$3
om_version=$4

bundle exec rake opsmgr:download_stemcell[$iaas,$stemcell_version]
bundle exec rake opsmgr:product:import_stemcell[$environment,$om_version,$(cat stemcell_reference.txt),cf]
