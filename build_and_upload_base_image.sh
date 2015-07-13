#!/bin/bash
set -e
set -x

docker_image_name=releng
docker_image_ref=releng:releng_base_ruby_22
base_docker_file=Dockerfile_base
image_folder=.
docker_registry_vsphere=docker.vsphere.gocd.cf-app.com:5000

wget https://github.com/rlmcpherson/s3gof3r/releases/download/v0.4.10/gof3r_0.4.10_linux_amd64.tar.gz -O ${image_folder}/include/gof3r.tar.gz
docker build -t ${docker_registry_vsphere}/${docker_image_ref} -f ${base_docker_file} ${image_folder}
echo "Build finished for ${docker_image_ref} in vsphere"

docker push ${docker_registry_vsphere}/${docker_image_ref}
echo "Finished uploading image ${docker_image_ref} to ${docker_registry_vsphere}"
