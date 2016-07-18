resource "aws_instance" "ops_manager" {
  depends_on = ["aws_security_group.ops_manager_security_group", "aws_subnet.public_subnet1"]
  ami = "${var.ops_manager_ami}"
  instance_type = "m3.medium"
  key_name = "${var.nat_key_pair_name}"
  vpc_security_group_ids = ["${aws_security_group.ops_manager_security_group.id}"]
  source_dest_check = false
  subnet_id = "${aws_subnet.public_subnet1.id}"

  root_block_device {
    volume_type = "gp2"
    volume_size = 50
  }

  tags {
    Name = "${var.env_name}-ops-manager"
  }
}

resource "aws_eip" "ops_manager_eip" {
  instance = "${aws_instance.ops_manager.id}"
  vpc = true
}

resource "aws_route53_zone" "pcf_zone" {
  name = "${var.env_name}.cf-app.com"

  tags {
    Name = "${var.env_name}-hosted-zone"
  }
}

resource "aws_route53_record" "ops_manager" {
  zone_id = "${aws_route53_zone.pcf_zone.id}"
  name = "pcf.${var.env_name}.cf-app.com"
  type = "A"
  ttl = "300"
  records = ["${aws_eip.ops_manager_eip.public_ip}"]
}

resource "aws_route53_record" "wildcard" {
  zone_id = "${aws_route53_zone.pcf_zone.id}"
  name = "*.${var.env_name}.cf-app.com"
  type = "CNAME"
  ttl = "5"
  records = ["${aws_elb.elb.dns_name}"]
}

resource "aws_route53_record" "ssh" {
  zone_id = "${aws_route53_zone.pcf_zone.id}"
  name = "ssh.${var.env_name}.cf-app.com"
  type = "CNAME"
  ttl = "5"
  records = ["${aws_elb.ssh_elb.dns_name}"]
}

resource "aws_route53_record" "tcp" {
  zone_id = "${aws_route53_zone.pcf_zone.id}"
  name = "tcp.${var.env_name}.cf-app.com"
  type = "CNAME"
  ttl = "5"
  records = ["${aws_elb.tcp_elb.dns_name}"]
}

output "pcf_zone_name_servers" {
  value = "${aws_route53_zone.pcf_zone.name_servers}"
}
