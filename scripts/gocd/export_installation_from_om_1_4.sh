source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:export_installation[$RELENG_ENV,1.4,export]'
