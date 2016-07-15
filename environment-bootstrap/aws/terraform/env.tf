variable "env_name" {}
variable "access_key" {}
variable "secret_key" {}
variable "region" {}
variable "availability_zone1" {}
variable "availability_zone2" {}
variable "nat_key_pair_name" {}
variable "ssl_certificate_arn" {}
variable "rds_db_name" {
  default = "bosh"
}
variable "rds_db_username" {}
variable "rds_db_password" {}
variable "rds_instance_class" {
  default = "db.m4.large"
}
