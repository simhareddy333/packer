resource “aws_instance” “web” {
 ami = “${lookup(var.amis,var.region)}”
 count = “${var.count}”
 key_name = “${var.key_name}”
 vpc_security_group_ids = [“${aws_security_group.instance.id}”]
 source_dest_check = false
 instance_type = “t2.micro”
tags {
 Name = “${format(“web-%03d”, count.index + 1)}”
 }
}
