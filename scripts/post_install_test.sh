#!/bin/bash
# Warning: Running CATs (this test) may pollute your deployment

export CATS_SHA=9fc2aa8
export SKIP_SSL_VALIDATION=true

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
  "api": "api.${APPS_DOMAIN}",
  "admin_user": "${ADMIN_USER}",
  "admin_password": "${ADMIN_PASSWORD}",
  "apps_domain": "${APPS_DOMAIN}",
  "skip_ssl_validation": ${SKIP_SSL_VALIDATION}
}
EOF

export CONFIG=$PWD/integration_config.json

./bin/test -nodes=3
