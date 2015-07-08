source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:delete_installation[$RELENG_ENV,1.5]'
