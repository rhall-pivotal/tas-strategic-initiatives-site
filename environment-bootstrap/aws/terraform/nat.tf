resource "aws_instance" "nat" {
  depends_on = ["aws_security_group.nat_security_group", "aws_subnet.public_subnet1"]
  ami = "ami-7da94839" # TODO: Determine this based on region
  instance_type = "t2.medium"
  key_name = "${var.nat_key_pair_name}"
  vpc_security_group_ids = ["${aws_security_group.nat_security_group.id}"]
  source_dest_check = false
  subnet_id = "${aws_subnet.public_subnet1.id}"

  tags {
    Name = "${var.env_name}-nat"
  }
}

resource "aws_eip" "nat_eip" {
  instance = "${aws_instance.nat.id}"
  vpc = true
}
