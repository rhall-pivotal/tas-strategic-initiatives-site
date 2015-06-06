set -e
./docker_run.sh bundle exec rake --trace krafa:claim_without_internet[${GO_PIPELINE_NAME}-${GO_PIPELINE_COUNTER}]
echo "export RELENG_ENV=$(./docker_run.sh bundle exec rake --trace krafa:print_environment_name[${GO_PIPELINE_NAME}-${GO_PIPELINE_COUNTER}])" > releng_env.sh
