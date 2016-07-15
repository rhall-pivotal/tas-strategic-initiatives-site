resource "aws_subnet" "public_subnet1" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.0.0/24"
  availability_zone = "us-west-1b" # TODO: don't hardcode this value

  tags {
    Name = "${var.env_name}-public-subnet1"
  }
}

resource "aws_subnet" "public_subnet2" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.1.0/24"
  availability_zone = "us-west-1c" # TODO: don't hardcode this value

  tags {
    Name = "${var.env_name}-public-subnet2"
  }
}

resource "aws_subnet" "private_subnet1" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.16.0/20"
  availability_zone = "us-west-1b" # TODO: don't hardcode this value

  tags {
    Name = "${var.env_name}-private-subnet1"
  }
}

resource "aws_subnet" "private_subnet2" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.32.0/20"
  availability_zone = "us-west-1c" # TODO: don't hardcode this value

  tags {
    Name = "${var.env_name}-private-subnet2"
  }
}

output "public_subnet1_id" {
  value = "${aws_subnet.public_subnet1.id}"
}

output "public_subnet2_id" {
  value = "${aws_subnet.public_subnet2.id}"
}
