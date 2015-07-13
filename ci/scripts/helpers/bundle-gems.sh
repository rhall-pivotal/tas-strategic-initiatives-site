#!/bin/bash

set -ex

# Ignore ssh fingerprints
mkdir -p $HOME/.ssh
echo -e "Host *\n\tStrictHostKeyChecking no\n\tUserKnownHostsFile=/dev/null" >> $HOME/.ssh/config

eval $(ssh-agent)
chmod 0600 ./ci/secrets/id_rsa
ssh-add ./ci/secrets/id_rsa

bundle
