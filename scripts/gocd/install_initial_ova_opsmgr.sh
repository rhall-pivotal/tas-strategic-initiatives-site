./docker_run.sh xvfb-run -a "rake opsmgr:bucket:get[ova/${INITIAL_OVA},${INITIAL_OVA}]"

source ./releng_env.sh
./docker_run.sh xvfb-run -a "rake opsmgr:install[$RELENG_ENV,$INITIAL_OVA]"
