#!/bin/bash

set -ex

if [ $# -eq 0 ] || [ $# -gt 2 ]; then
  echo "Usage: $0 ENV_NAME [DEPLOYMENT_NAME] " >&2
  exit 1
fi

env_name=$1
deployment=${2:-cf}

$(bundle exec rake opsmgr:microbosh:target[$env_name])
bosh_command=$(bundle exec rake opsmgr:microbosh:command[$env_name])
deployment_name=$($bosh_command deployments | grep -Eoh "${deployment}-[0-9a-f]{8,}")
deployment_file="${TMPDIR}/${env_name}$$.yml"

$bosh_command download manifest $deployment_name $deployment_file
$bosh_command deployment $deployment_file
