#!/bin/bash

set -ex

if [ ! -z $GO_ENVIRONMENT_NAME ] && [ $GO_ENVIRONMENT_NAME == "vSphere" ]; then
    echo "Error: On a vSphere agent"
    exit 1
fi
