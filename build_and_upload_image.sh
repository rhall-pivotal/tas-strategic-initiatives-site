#!/bin/bash
set -e
set -x

docker_image_tag=`git describe --dirty`

docker_image_name=releng
docker_image_ref=${docker_image_name}:${docker_image_tag}
docker_file=Dockerfile_gocd
image_folder=.
docker_registry_vsphere=docker.vsphere.gocd.cf-app.com:5000
docker_registry_aws=docker.gocd.cf-app.com:5000

wget https://github.com/rlmcpherson/s3gof3r/releases/download/v0.4.10/gof3r_0.4.10_linux_amd64.tar.gz -O ${image_folder}/include/gof3r.tar.gz
docker build -t ${docker_registry_vsphere}/${docker_image_ref} -f ${docker_file}  ${image_folder}
echo "Build finished for ${docker_image_ref} in vsphere"

docker build -t ${docker_registry_aws}/${docker_image_ref} -f ${docker_file} ${image_folder}
echo "Build finished for ${docker_image_ref} in aws"

docker push ${docker_registry_vsphere}/${docker_image_ref}
echo "Finished uploading image ${docker_image_ref} to ${docker_registry_vsphere}"

docker push ${docker_registry_aws}/${docker_image_ref}
echo "Finished uploading image ${docker_image_ref} to ${docker_registry_aws}"
