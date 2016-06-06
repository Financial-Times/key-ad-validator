# Configure the AWS Provider
provider "aws" {}

variable "ldap_password"{}
variable "ldap_user"{}
variable "ldap_server"{}
variable "ldap_port"{}
variable "keys_uri"{}
variable "deploy_user"{}

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
	# The connection block tells our provisioner how to
	# communicate with the resource (instance)
	connection {
		# Connect to private IP, 22 blocked for public
		host = "${self.private_ip}"
		# AD username to deploy the artifact
		user = "${var.deploy_user}"
		# Timeout after 5 minutes
		timeout = "5m"

		# The connection will use the local SSH agent for authentication.
	}

	provisioner "remote-exec" {
  		inline = [
  			"sudo yum install -y docker",
  			"sudo service docker start",
			"sudo docker kill key-ad-validator || echo Not running",
			"export KEYS_URI=${var.keys_uri}",
			"export LDAP_SERVER=${var.ldap_server}",
			"export LDAP_PORT=${var.ldap_port}",
			"export LDAP_USER=${var.ldap_user}",
			"export LDAP_PASSWORD=${var.ldap_password}",
			"sudo docker run -d -p 80:8080 --name key-ad-validator coco/key-ad-validator"
  		]
  	}
}
