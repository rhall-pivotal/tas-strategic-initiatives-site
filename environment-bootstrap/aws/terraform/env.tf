variable "env_name" {}
variable "access_key" {}
variable "secret_key" {}
variable "region" {}
variable "nat_key_pair_name" {}
variable "ssl_certificate_arn" {}

variable "internal_cidr_block" {
  default = "10.0.0.0/16"
}
