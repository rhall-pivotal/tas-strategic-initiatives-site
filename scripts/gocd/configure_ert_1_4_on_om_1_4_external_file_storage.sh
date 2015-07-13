source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake ert:configure_external_file_storage[$RELENG_ENV,1.4,1.4]'
