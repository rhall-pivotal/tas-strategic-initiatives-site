source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:download[aws,stable]'

# Temporarily extract us-east-1 AMI until we can bump vm_shepherd
cat ops_man_image | grep us-east-1 | sed s/us-east-1\:\ // > ops_man_image.new
mv ops_man_image ops_man_image.old
mv ops_man_image.new ops_man_image
