#!/bin/bash

bundle install
bundle exec vara-download-artifacts --product-metadata=metadata/cf.yml
bundle exec vara-build-pivotal --product-metadata=metadata/cf.yml
