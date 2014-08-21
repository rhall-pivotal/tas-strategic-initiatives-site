#!/bin/bash
# Warning: Running CATs (this test) may pollute your deployment

set -x

export SCRIPT_DIR=$(dirname $0)
export SKIP_SSL_VALIDATION=true

function set_cats_sha() {
  export CF_RELEASE_TAG=$(ruby -r 'yaml' -e 'puts "v" + YAML.load_file(ENV["SCRIPT_DIR"] + "/../metadata_parts/binaries.yml").fetch("releases").find{ |release| release["name"] == "cf"}.fetch("version")')
  export CATS_DIR='src/acceptance-tests'

  cd ${GOPATH}
  mkdir -p src/github.com/cloudfoundry/
  cd src/github.com/cloudfoundry/
  if [ -d cf-release ]; then
    cd cf-release
    git fetch origin master
  else
    git clone git@github.com:cloudfoundry/cf-release.git
    cd cf-release
  fi

  export CATS_SHA=$(git ls-tree ${CF_RELEASE_TAG} ${CATS_DIR} | awk '{print $3}')
}

set_cats_sha

cd $GOPATH
mkdir -p src/github.com/cloudfoundry/
cd src/github.com/cloudfoundry/
if [ -d cf-acceptance-tests ]; then
  cd cf-acceptance-tests
  git pull origin master
else
  git clone git@github.com:cloudfoundry/cf-acceptance-tests.git
  cd cf-acceptance-tests
fi

git checkout $CATS_SHA

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

echo "WARNING: skipping Security Groups CATs" >&2
./bin/test -nodes=3 --noColor -skipPackage=security_groups
