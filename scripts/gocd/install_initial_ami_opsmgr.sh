echo $INITIAL_AMI > ami_reference.txt

source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:install[$RELENG_ENV,ami_reference.txt]'
