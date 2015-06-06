set -e
KEY=`grep -o -e cf.* cf-pivotal.blobkey`
../p-runtime/scripts/run_in_docker.sh gof3r cp s3://releng-products/untested/${KEY} ${KEY}
../p-runtime/scripts/run_in_docker.sh gof3r cp ${KEY} s3://releng-products/runtime/${KEY}

../p-runtime/scripts/run_in_docker.sh s3cmd -c s3cfg.s3 cp s3://releng-products/untested/${KEY}.md5 s3://releng-products/runtime/${KEY}.md5

../p-runtime/scripts/run_in_docker.sh s3cmd -c s3cfg.s3 cp s3://releng-products/untested/${KEY}.yml s3://releng-products/runtime/${KEY}.yml
