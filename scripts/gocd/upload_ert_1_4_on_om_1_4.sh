set -e
source ./releng_env.sh
./docker_run.sh scripts/runtime/download_pivotal_with_gof3r.sh
./docker_run.sh xvfb-run -a 'rake opsmgr:product:upload_add[$RELENG_ENV,1.4,$RUNTIME_DOT_PIVOTAL_FILE,cf]'
