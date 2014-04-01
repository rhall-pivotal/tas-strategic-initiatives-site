#!/bin/bash

bundle install
bundle exec vara-build-metadata --product-dir=$PWD
bundle exec vara-download-artifacts --product-metadata=metadata/cf.yml
bundle exec vara-build-pivotal --product-metadata=metadata/cf.yml --rc="-build${BUILD_NUMBER}"
