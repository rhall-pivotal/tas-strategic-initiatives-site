#!/bin/bash

set -e

ci_info=`ruby -e "require 'yaml'; require 'json'; y=YAML.load_file('/Users/pivotal/.flyrc'); puts y.to_json"`

url=`echo $ci_info | jq '.targets.ci.api' | sed -e 's/"//g'`
username=`echo $ci_info | jq '.targets.ci.username' | sed -e 's/"//g'`
password=`echo $ci_info | jq '.targets.ci.password' | sed -e 's/"//g'`

curl --silent -k -u ${username}:${password} ${url}/api/v1/builds> ~/b-new.json
newer_than=`jq '. | first| .id' ~/b.json`

jq ".[]" ~/b-new.json | \
 jq "select(.id > ${newer_than})" | \
 jq 'select((.status | contains("succeeded")) | not) ' | \
 jq 'select((.status | contains("started")) | not) ' | \
 jq 'select((.status | contains("pending")) | not) '

cat ~/b-new.json | \
 jq 'map(select((.status | contains("succeeded")) | not)) ' | \
 jq 'map(select((.status | contains("started")) | not)) ' | \
 jq 'map(select((.status | contains("pending")) | not)) ' > ~/b.json
