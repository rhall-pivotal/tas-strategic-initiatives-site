set -e
export RUNTIME_DOT_PIVOTAL_FILE=$(cat ../cf-pivotal-artifacts-metadata/cf-pivotal.blobkey)
./docker_run.sh xvfb-run -a 'rake opsmgr:bucket:get[${RUNTIME_DOT_PIVOTAL_FILE},new_runtime.pivotal]'

source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:product:upload_upgrade[$RELENG_ENV,1.5,new_runtime.pivotal,cf]'
