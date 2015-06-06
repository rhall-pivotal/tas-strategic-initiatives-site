source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:download[vsphere,latest]'
