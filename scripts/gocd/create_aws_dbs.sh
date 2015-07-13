source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake ert:create_aws_dbs[$RELENG_ENV]'
