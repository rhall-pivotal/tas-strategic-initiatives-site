module "releng" {
  source = "./releng"

  env_name = "${var.env_name}"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
  availability_zone1 = "${var.availability_zone1}"
  availability_zone2 = "${var.availability_zone2}"
  nat_key_pair_name = "${var.nat_key_pair_name}"
  ops_manager_ami = "${var.ops_manager_ami}"
  rds_db_name = "${var.rds_db_name}"
  rds_db_username = "${var.rds_db_username}"
  rds_db_password = "${var.rds_db_password}"
  rds_instance_class = "${var.rds_instance_class}"
  rds_instance_count = "${var.rds_instance_count}"
}

output "ops_manager_bucket" {
  value = "${module.releng.ops_manager_bucket}"
}

output "ert_buildpacks_bucket" {
  value = "${module.releng.ert_buildpacks_bucket}"
}

output "ert_droplets_bucket" {
  value = "${module.releng.ert_droplets_bucket}"
}

output "ert_packages_bucket" {
  value = "${module.releng.ert_packages_bucket}"
}

output "ert_resources_bucket" {
  value = "${module.releng.ert_resources_bucket}"
}

output "elb_dns_name" {
  value = "${module.releng.elb_dns_name}"
}

output "ssh_elb_dns_name" {
  value = "${module.releng.ssh_elb_dns_name}"
}

output "tcp_elb_dns_name" {
  value = "${module.releng.tcp_elb_dns_name}"
}

output "iam_user_name" {
  value = "${module.releng.iam_user_name}"
}

output "iam_user_access_key" {
  value = "${module.releng.iam_user_access_key}"
}

output "iam_user_secret_access_key" {
  value = "${module.releng.iam_user_secret_access_key}"
}

output "rds_address" {
  value = "${module.releng.rds_address}"
}

output "rds_port" {
  value = "${module.releng.rds_port}"
}

output "rds_username" {
  value = "${module.releng.rds_username}"
}

output "rds_password" {
  value = "${module.releng.rds_password}"
}

output "rds_db_name" {
  value = "${module.releng.rds_db_name}"
}

output "ops_manager_security_group_id" {
  value = "${module.releng.ops_manager_security_group_id}"
}

output "vms_security_group_id" {
  value = "${module.releng.vms_security_group_id}"
}

output "public_subnet1_id" {
  value = "${module.releng.public_subnet1_id}"
}

output "public_subnet2_id" {
  value = "${module.releng.public_subnet2_id}"
}

output "public_subnet1_availability_zone" {
  value = "${module.releng.public_subnet1_availability_zone}"
}

output "public_subnet2_availability_zone" {
  value = "${module.releng.public_subnet2_availability_zone}"
}

output "private_subnet1_id" {
  value = "${module.releng.private_subnet1_id}"
}

output "private_subnet2_id" {
  value = "${module.releng.private_subnet2_id}"
}

output "private_subnet1_availability_zone" {
  value = "${module.releng.private_subnet1_availability_zone}"
}

output "private_subnet2_availability_zone" {
  value = "${module.releng.private_subnet2_availability_zone}"
}

output "vpc_id" {
  value = "${module.releng.vpc_id}"
}

module "shared_dns" {
  source = "./shared_dns"

  env_name = "${var.env_name}"
  access_key = "${var.shared_dns_access_key}"
  secret_key = "${var.shared_dns_secret_key}"
  region = "${var.region}"
  shared_dns_hosted_zone_id = "${var.shared_dns_hosted_zone_id}"
  releng_dns_hosted_zone_name_servers="${module.releng.pcf_zone_name_servers}"
}
