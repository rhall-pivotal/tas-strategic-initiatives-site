source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:trigger_install[$RELENG_ENV,1.4,45]'
