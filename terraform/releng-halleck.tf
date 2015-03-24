provider "aws" {
  region = "us-east-1"
  access_key = "AKIAI43577ARGHXPWIMQ"
  secret_key = "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"
}

resource "aws_vpc" "releng-halleck" {
  cidr_block = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_hostnames = true

  tags {
    Name = "releng-halleck"
  }
}

/******** PUBLIC SUBNET *********/

resource "aws_internet_gateway" "releng-halleck" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
}

resource "aws_subnet" "releng-halleck-public" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  cidr_block = "10.0.1.0/24"
  availability_zone = "us-east-1d"

  tags {
    Name = "releng-halleck-public"
  }
}

resource "aws_route_table" "releng-halleck-public" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.releng-halleck.id}"
  }

  tags {
    Name = "releng-halleck-public"
  }
}

resource "aws_route_table_association" "releng-halleck-public" {
  subnet_id = "${aws_subnet.releng-halleck-public.id}"
  route_table_id = "${aws_route_table.releng-halleck-public.id}"
}

/******** NAT STUFF *********/
resource "aws_security_group" "releng-halleck-nat" {
  name = "releng-halleck-nat"
  description = "Ops Manager releng-halleck for NAT box"
  vpc_id = "${aws_vpc.releng-halleck.id}"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  ingress {
    self = true
    from_port = 0
    to_port = 65535
    protocol = "udp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  ingress {
    self = true
    from_port = -1
    to_port = -1
    protocol = "icmp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  tags {
    Name = "releng-halleck-nat"
  }
}

resource "aws_instance" "releng-halleck-nat" {
    ami = "ami-184dc970"
    key_name = "bosh"
    instance_type = "t2.micro"
    subnet_id = "${aws_subnet.releng-halleck-public.id}"
    associate_public_ip_address = true
    private_ip = "10.0.1.4"
    security_groups = ["${aws_security_group.releng-halleck-nat.id}"]
    source_dest_check = false
    tags {
        Name = "releng-halleck-nat"
        do_not_terminate = "true"
    }
}

/******** PRIVATE SUBNET *********/

resource "aws_subnet" "releng-halleck-private" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  cidr_block = "10.0.2.0/24"
  availability_zone = "us-east-1d"
  tags {
    Name = "releng-halleck-private"
  }
}

resource "aws_route_table" "releng-halleck-private" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  route {
    cidr_block = "0.0.0.0/0"
    instance_id = "${aws_instance.releng-halleck-nat.id}"
  }
  tags {
    Name = "releng-halleck-private"
  }
}

resource "aws_route_table_association" "releng-halleck-private" {
  subnet_id = "${aws_subnet.releng-halleck-private.id}"
  route_table_id = "${aws_route_table.releng-halleck-private.id}"
}

/******** BOSH VM security group *********/

resource "aws_security_group" "releng-halleck" {
  name = "releng-halleck"
  description = "Ops Manager releng-halleck"
  vpc_id = "${aws_vpc.releng-halleck.id}"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    self = true
    from_port = 0
    to_port = 65535
    protocol = "tcp"
  }

  ingress {
    self = true
    from_port = 0
    to_port = 65535
    protocol = "udp"
  }

  ingress {
    self = true
    from_port = -1
    to_port = -1
    protocol = "icmp"
  }

  tags {
    Name = "releng-halleck"
  }
}

/******** ELASTIC IP *********/

resource "aws_eip" "releng-halleck-ops-manager" {
  vpc = true
}

/******** ELB/RDS SUBNETS *********/

resource "aws_subnet" "releng-halleck-elb-rds1" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  cidr_block = "10.0.3.0/25"
  availability_zone = "us-east-1d"

  tags {
    Name = "releng-halleck-elb-rds1"
  }
}

resource "aws_route_table" "releng-halleck-elb-rds1" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.releng-halleck.id}"
  }

  tags {
    Name = "releng-halleck-elb-rds1"
  }
}

resource "aws_route_table_association" "releng-halleck-elb-rds1" {
  subnet_id = "${aws_subnet.releng-halleck-elb-rds1.id}"
  route_table_id = "${aws_route_table.releng-halleck-elb-rds1.id}"
}

resource "aws_subnet" "releng-halleck-elb-rds2" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  cidr_block = "10.0.3.128/25"
  availability_zone = "us-east-1e"

  tags {
    Name = "releng-halleck-elb-rds2"
  }
}

resource "aws_route_table" "releng-halleck-elb-rds2" {
  vpc_id = "${aws_vpc.releng-halleck.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.releng-halleck.id}"
  }

  tags {
    Name = "releng-halleck-elb-rds2"
  }
}

resource "aws_route_table_association" "releng-halleck-elb-rds2" {
  subnet_id = "${aws_subnet.releng-halleck-elb-rds2.id}"
  route_table_id = "${aws_route_table.releng-halleck-elb-rds2.id}"
}

/******** ELB *********/

resource "aws_security_group" "releng-halleck-open" {
  name = "releng-halleck-open"
  description = "Allow all inbound traffic"
  vpc_id = "${aws_vpc.releng-halleck.id}"

  ingress {
      from_port = 0
      to_port = 65535
      protocol = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
  }
}


resource "aws_elb" "releng-halleck" {
  name = "releng-halleck"
  subnets = ["${aws_subnet.releng-halleck-elb-rds1.id}", "${aws_subnet.releng-halleck-elb-rds2.id}"]
  security_groups = ["${aws_security_group.releng-halleck-open.id}"]

  listener {
    instance_port = 80
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  listener {
    instance_port = 22
    instance_protocol = "tcp"
    lb_port = 2222
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 2
    target = "TCP:22"
    interval = 5
  }
}

/******** S3 *********/

resource "aws_s3_bucket" "releng-halleck" {
    bucket = "releng-halleck"
    acl = "public-read"
}

/******** RDS *********/

resource "aws_db_subnet_group" "releng-halleck" {
    name = "releng-halleck"
    description = "RDS subnets for releng-halleck"
    subnet_ids = ["${aws_subnet.releng-halleck-elb-rds1.id}", "${aws_subnet.releng-halleck-elb-rds2.id}"]
}

resource "aws_db_instance" "releng-halleck" {
    identifier = "releng-halleck-bosh"
    allocated_storage = 5
    engine = "mysql"
    engine_version = "5.6.21"
    instance_class = "db.t2.small"
    name = "bosh"
    publicly_accessible = true
    username = "boshuser"
    password = "boshpass1234"

    skip_final_snapshot = true

    vpc_security_group_ids = ["${aws_security_group.releng-halleck-open.id}"]
    db_subnet_group_name = "releng-halleck"
}

