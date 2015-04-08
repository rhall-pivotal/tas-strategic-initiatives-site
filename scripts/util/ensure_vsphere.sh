#!/bin/bash

set -ex

if [ -z $GO_ENVIRONMENT_NAME ]; then
    echo "Not on ANY environment"
    exit 1
else
    echo "GO_ENVIRONMENT_NAME is ${GO_ENVIRONMENT_NAME}"
    if [ $GO_ENVIRONMENT_NAME != "vSphere" ]; then
        exit 1
    fi
fi
