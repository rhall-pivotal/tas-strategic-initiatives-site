source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake ert:update_dns_elb[$RELENG_ENV]'
