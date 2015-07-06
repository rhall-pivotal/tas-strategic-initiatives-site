source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake ert:configure_external_dbs[$RELENG_ENV,1.4,1.4]'
