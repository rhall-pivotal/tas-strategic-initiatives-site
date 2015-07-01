./docker_run.sh xvfb-run -a "rake opsmgr:bucket:get[vcd/${INITIAL_VCD},${INITIAL_VCD}]"

source ./releng_env.sh
./docker_run.sh xvfb-run -a "rake opsmgr:install[$RELENG_ENV,$INITIAL_VCD]"
