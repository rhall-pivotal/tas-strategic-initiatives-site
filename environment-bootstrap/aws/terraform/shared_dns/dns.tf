resource "aws_route53_record" "pcf-ns" {
  zone_id = "${var.shared_dns_hosted_zone_id}"
  name = "${var.env_name}.cf-app.com"
  type = "NS"
  ttl = "300"
  records = ["${var.releng_dns_hosted_zone_name_servers}"]
}
