provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region = "${var.region}"
}

resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_hostnames = true

  tags {
    Name = "${var.env_name}-vpc"
  }
}

output "vpc_id" {
  value = "${aws_vpc.vpc.id}"
}
