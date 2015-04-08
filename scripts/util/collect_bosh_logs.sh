#!/bin/bash

set -e

if [ $# -ne 2 ]; then
  echo "Usage: $0 ENV_NAME LOG_OUT_DIR" >&2
  exit 1
fi

env_name=$1
target_log_dir=$2
bosh_command=$(bundle exec rake opsmgr:bosh:command[$env_name])

mkdir -p ${target_log_dir}

tasks_file="${target_log_dir}/tasks.out"
echo "Retrieving list of recent BOSH tasks..." >&2
$bosh_command tasks recent --no-filter > ${tasks_file}

for task in $(grep -Eoh '^\|[[:space:]]+[[:digit:]]+' ${tasks_file} | awk '{print $2}'); do
  echo "Retrieving log for BOSH task ${task}..." >&2
  $bosh_command task ${task} --debug > ${target_log_dir}/${task}.debug.log
done
