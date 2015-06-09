source ./releng_env.sh
./docker_run.sh bundle exec krafa_client release_claim --environment-name=$RELENG_ENV
