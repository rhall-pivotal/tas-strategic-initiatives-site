#!/bin/bash
set -e

s3cmd -c /s3cfg.riak put $*
