source ./releng_env.sh
./docker_run.sh ./scripts/util/collect_bosh_logs.sh $RELENG_ENV ../bosh_logs
