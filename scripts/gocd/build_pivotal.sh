if [ -e ../vara_extra_flags.sh ]
then
  source ../vara_extra_flags.sh
fi

./docker_run.sh ./scripts/build_pivotal.sh
