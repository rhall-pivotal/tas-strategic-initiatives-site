#!/bin/bash

set -eux

grep 'job' instance_groups/*.yml | grep -v "control.yml" | grep -v "database.yml" | grep -v "blobstore.yml" | grep -v "compute.yml" | grep -v 'placeholder' | sed 's|instance_groups/\(.*\).yml:- $( job "\(.*\)" )|\1,\2|' >iteration0.csv

for line in $( cat iteration0.csv ); do echo "$line,$( grep "release: " "jobs/$( echo $line | awk -F, '{print $2}' ).yml" | sed 's/release: //')" ; done >iteration1.csv
