set -e
source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:add_first_user[$RELENG_ENV,1.4]'
./docker_run.sh xvfb-run -a 'rake opsmgr:microbosh:configure[$RELENG_ENV,1.4]'
