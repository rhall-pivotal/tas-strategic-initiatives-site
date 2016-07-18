variable "env_name" {}
variable "access_key" {}
variable "secret_key" {}
variable "region" {}
variable "availability_zone1" {}
variable "availability_zone2" {}
variable "nat_key_pair_name" {}
variable "ops_manager_ami" {
  default = "ami-2e02454e"
}

variable "rds_db_name" {
  default = "bosh"
}
variable "rds_db_username" {}
variable "rds_db_password" {}
variable "rds_instance_class" {
  default = "db.m4.large"
}
variable "rds_instance_count" {
  default = 1
}
variable "shared_dns_access_key" {}
variable "shared_dns_secret_key" {}
variable "shared_dns_hosted_zone_id" {}
