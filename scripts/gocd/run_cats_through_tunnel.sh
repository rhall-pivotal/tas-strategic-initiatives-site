source ./releng_env.sh
./docker_run.sh bundle exec rake ert:run_cats[$RELENG_ENV]
