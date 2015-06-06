source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:install[$RELENG_ENV,ops_man_image]'
