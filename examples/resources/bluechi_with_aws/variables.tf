// SPDX-License-Identifier: MIT-0

/*********************/
/* Variables for AWS */
/*********************/

# tuple(public key path, private key path)
variable "ssh_key_pair" {
  type    = tuple([string, string])
  default = ["~/.ssh/bluechi_aws.pub", "~/.ssh/bluechi_aws"]
}

variable "ssh_user" {
  type    = string
  default = "ec2-user"
}

# Developer AMI, changes frequently
variable "autosd_ami" {
  type    = string
  default = "ami-0b2337f1f6379076e"
}

variable "instance_type" {
  type    = string
  default = "t3a.micro"
}

/*************************/
/* Variables for BlueChi */
/*************************/

variable "bluechi_manager_port" {
  type    = number
  default = 3030
}

variable "bluechi_nodes" {
  type    = list(string)
  default = ["main", "worker1"]
}
