# Configure the AWS Provider
provider "aws" {}

variable "ldap_password"{}
variable "ldap_user"{}
variable "ldap_server"{}
variable "ldap_port"{}
variable "keys_uri"{}

# Remote state file
resource "terraform_remote_state" "state_file" {
	backend = "s3"
	config {
		bucket = "ft.terraform.state-files"
		key = "coco-key-ad-valid-svc/terraform.tfstate"
		region = "eu-west-1"
	}
}

# Userdata
resource "template_file" "userdata" {
	template = "${file("userdata.tpl")}"
	vars {
	        ldap_password = "${var.ldap_password}"
	        ldap_user = "${var.ldap_user}"
	        ldap_server = "${var.ldap_server}"
	        ldap_port = "${var.ldap_port}"
	        keys_uri = "${var.keys_uri}"
	}
}

# Create a server
resource "aws_instance" "key-ad-validator" {
	ami = "ami-b0ac25c3" # EU Ireland - HVM (SSD) EBS-Backed 64-bit
	vpc_security_group_ids = ["sg-1beb8b7e", "sg-1dec8c78"]
        instance_type = "t2.nano"
	user_data = "${template_file.userdata.rendered}"
        subnet_id = "subnet-c86755bc" # FT-Tech-Infra-Prod-AzA-Public
	associate_public_ip_address = true
        iam_instance_profile = "FT-Linux-Role"
	root_block_device {
		volume_size = 10
	}
	tags {
		Name = "CoCo Key AD Validator Service"
		environment = "p"
		ipCode = "P196"
		systemCode = "Universal-Publishing"
		description = "Universal-Publishing"
		teamDL = "universal.publishing.platform@ft.com"
		stopSchedule = "nostop"
	}
}
