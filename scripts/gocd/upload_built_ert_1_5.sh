export RUNTIME_DOT_PIVOTAL_FILE=$(cat ../cf-pivotal-artifacts-metadata/cf-pivotal.blobkey)
./docker_run.sh xvfb-run -a 'rake opsmgr:bucket:get[${RUNTIME_DOT_PIVOTAL_FILE},runtime.pivotal]'
source ./releng_env.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:product:upload_add[$RELENG_ENV,1.5,runtime.pivotal,cf]'
