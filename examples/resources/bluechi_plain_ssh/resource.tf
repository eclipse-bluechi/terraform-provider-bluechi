// SPDX-License-Identifier: MIT-0
resource "bluechi_node" "main" {

  ssh = {
    host                     = "127.0.0.1:2020"
    user                     = "root"
    password                 = ""
    private_key_path         = "~/.ssh/id_rsa"
    accept_host_key_insecure = true
  }

  bluechi_controller = {
    allowed_node_names = [
      "main", "worker1", "worker2", "worker3",
    ]
    manager_port = 3030
    log_level    = "DEBUG"
    log_target   = "stderr-full"
    log_is_quiet = false
  }

  bluechi_agent = {
    node_name          = "main"
    manager_host       = "127.0.0.1"
    manager_port       = 3030
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}

resource "bluechi_node" "worker1" {

  ssh = {
    host                     = "127.0.0.1:2021"
    user                     = "root"
    password                 = ""
    private_key_path         = "~/.ssh/id_rsa"
    accept_host_key_insecure = true
  }

  bluechi_agent = {
    node_name          = "worker1"
    manager_host       = "127.0.0.1"
    manager_port       = 3030
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}

resource "bluechi_node" "worker2" {

  ssh = {
    host                     = "127.0.0.1:2022"
    user                     = "root"
    password                 = ""
    private_key_path         = "~/.ssh/id_rsa"
    accept_host_key_insecure = true
  }

  bluechi_agent = {
    node_name          = "worker2"
    manager_host       = "127.0.0.1"
    manager_port       = 3030
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}

resource "bluechi_node" "worker3" {

  ssh = {
    host                     = "127.0.0.1:2023"
    user                     = "root"
    password                 = ""
    private_key_path         = "~/.ssh/id_rsa"
    accept_host_key_insecure = true
  }

  bluechi_agent = {
    node_name          = "worker3"
    manager_host       = "127.0.0.1"
    manager_port       = 3030
    manager_address    = ""
    heartbeat_interval = 5000
    log_level          = "DEBUG"
    log_target         = "stderr-full"
    log_is_quiet       = false
  }
}
