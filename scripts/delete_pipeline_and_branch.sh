#!/bin/bash
set -e

ci_info=`ruby -e "require 'yaml'; require 'json'; y=YAML.load_file('/Users/pivotal/.flyrc'); puts y.to_json"`

url=`echo $ci_info | jq '.targets.ci.api' | sed -e 's/"//g'`
username=`echo $ci_info | jq '.targets.ci.username' | sed -e 's/"//g'`
password=`echo $ci_info | jq '.targets.ci.password' | sed -e 's/"//g'`

pipelines=`curl --silent -k -u ${username}:${password}  ${url}/api/v1/pipelines |\
  jq '.[]|.name' |\
  sed -e 's/"//g'`

PS3="Select pipeline: "
select PIPELINE in $pipelines QUIT;
do
  case $PIPELINE in
    "QUIT")
      break
      ;;
    *)
       branch_name=`echo $PIPELINE | sed -e 's/::/\//g'`
       echo $PIPELINE has $branch_name
       fly -t ci d $PIPELINE
       git branch -d $branch_name
       git push origin --delete $branch_name
      ;;
  esac
done
