set -e
claimed_environment=$(./docker_run.sh bundle exec krafa_client claim --iaas=${IAAS:-vsphere} --has-internet --claim-reason="${GO_PIPELINE_NAME}-${GO_PIPELINE_COUNTER}")
echo "export RELENG_ENV=$claimed_environment" > releng_env.sh
