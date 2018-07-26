#!/bin/bash

set -eu
#set -eux

function remove_srt_instance_groups() {
  grep -v "control.yml" | grep -v "database.yml" | grep -v "blobstore.yml" | grep -v "compute.yml"
}

function remove_placeholder_jobs() {
  grep -v 'placeholder'
}

function pluck_instance_group() {
  sed 's|instance_groups/\(.*\).yml:- $( job ".*" )|\1|'
}

function pluck_job_name() {
  sed 's|instance_groups/.*.yml:- $( job "\(.*\)" )|\1|'
}

function lookup_release_name() {
  local job_name
  job_name="$1"

  grep "release: " "jobs/${job_name}.yml" | sed 's/release: //'
}

# the default IFS is whitespace and we want `line` to be a complete line.
IFS=$'\n'
for line in $(grep 'job' instance_groups/*.yml | remove_srt_instance_groups | remove_placeholder_jobs); do
  INSTANCE_GROUP=$(echo "$line" | pluck_instance_group)
  JOB_NAME=$(echo "$line" | pluck_job_name)
  RELEASE_NAME=$(lookup_release_name "$JOB_NAME")

  echo "$INSTANCE_GROUP,$JOB_NAME,$RELEASE_NAME"
done
