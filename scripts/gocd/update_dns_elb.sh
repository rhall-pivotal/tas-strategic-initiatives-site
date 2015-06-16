source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:update_dns_elb[$RELENG_ENV]'
