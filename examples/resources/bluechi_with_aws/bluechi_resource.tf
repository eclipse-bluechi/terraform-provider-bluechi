resource "bluechi_node" "main" {

  depends_on = [aws_internet_gateway.autosd_demo_ig]

  ssh = {
    host                     = "${aws_instance.ec2main.*.public_ip[0]}:22"
    user                     = var.ssh_user
    password                 = ""
    private_key_path         = var.ssh_key_pair[1]
    accept_host_key_insecure = true
  }

  bluechi_controller = {
    allowed_node_names = [
      var.bluechi_nodes[0], var.bluechi_nodes[1],
    ]
    manager_port = var.bluechi_manager_port
    log_level    = "DEBUG"
    log_target   = "stderr-full"
    log_is_quiet = false
  }

  bluechi_agent = {
    node_name          = var.bluechi_nodes[0]
    manager_host       = "127.0.0.1"
    manager_port       = var.bluechi_manager_port
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}

resource "bluechi_node" "worker1" {

  depends_on = [aws_internet_gateway.autosd_demo_ig]

  ssh = {
    host                     = "${aws_instance.ec2worker1.*.public_ip[0]}:22"
    user                     = var.ssh_user
    password                 = ""
    private_key_path         = var.ssh_key_pair[1]
    accept_host_key_insecure = true
  }

  bluechi_agent = {
    node_name          = var.bluechi_nodes[1]
    manager_host       = "${aws_instance.ec2main.*.private_ip[0]}"
    manager_port       = var.bluechi_manager_port
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}
