#!/bin/bash
# Warning: Running CATs (this test) may pollute your deployment

set -ex

export SCRIPT_DIR=$(dirname $0)
export SKIP_SSL_VALIDATION=true

PCF_RELEASE_REF=$(ruby -r 'yaml' -e 'puts "v" + YAML.load_file(ENV["SCRIPT_DIR"] + "/../metadata_parts/binaries.yml").fetch("releases").find{ |release| release["name"] == "cf"}.fetch("version")')

rm -rf pcf-release
git clone -b ${PCF_RELEASE_REF} git@github.com:pivotal-cf/pcf-release.git
cd pcf-release

DEST_DIR="${GOPATH}/src/github.com/cloudfoundry/cf-acceptance-tests/"

git submodule update --init --recursive -- 'src/acceptance-tests'
rm -rf ${DEST_DIR}
mkdir -p ${DEST_DIR}
mv src/acceptance-tests/* ${DEST_DIR}
cd ${DEST_DIR}

cat > integration_config.json <<EOF
{
  "api": "api.${SYSTEM_DOMAIN}",
  "admin_user": "${ADMIN_USER}",
  "admin_password": "${ADMIN_PASSWORD}",
  "apps_domain": "${APPS_DOMAIN}",
  "skip_ssl_validation": ${SKIP_SSL_VALIDATION}
}
EOF

export CONFIG=$PWD/integration_config.json

echo "WARNING: only running applications and services CATs; other suites skipped" >&2
./bin/test -nodes=3 --noColor apps services && ./bin/test_operator -nodes=3