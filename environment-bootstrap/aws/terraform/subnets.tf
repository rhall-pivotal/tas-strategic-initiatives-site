resource "aws_subnet" "public_subnet1" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.0.0/24"
  availability_zone = "us-west-1b" # TODO: don't hardcode this value

  tags {
    Name = "public-subnet1"
  }
}

resource "aws_subnet" "public_subnet2" {
  depends_on = ["aws_vpc.vpc"]
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.0.1.0/24"
  availability_zone = "us-west-1c" # TODO: don't hardcode this value

  tags {
    Name = "public-subnet2"
  }
}

output "public_subnet1_id" {
  value = "${aws_subnet.public_subnet1.id}"
}

output "public_subnet2_id" {
  value = "${aws_subnet.public_subnet2.id}"
}
