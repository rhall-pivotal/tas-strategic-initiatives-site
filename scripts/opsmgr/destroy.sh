#!/bin/bash -e

bundle
bundle exec rake --trace opsmgr:destroy:opsmgr[${RELENG_ENV}]
