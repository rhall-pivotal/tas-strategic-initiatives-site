#!/bin/bash
set -e

pipelines=`fly -t ci ps | tail +1 | cut -f1 -d' '`
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
       fly -t ci dp -p $PIPELINE
       git branch -d $branch_name
       git push origin --delete $branch_name
      ;;
  esac
done
