#!/bin/bash -xe

PUBKEYPATH=~/.ssh/id_rsa.pub
PUBKEY=$( cat $PUBKEYPATH )

SCRIPT_DIR=$( realpath "$0"  )
SCRIPT_DIR=$(dirname "$SCRIPT_DIR")

CONTAINER_NAMES=(main worker1 worker2 worker3)

function build_image(){
    if [[ "$1" == "bluechi" ]]; then
        podman build -t localhost/bluechi -f $SCRIPT_DIR/bluechi.image
    elif [[ "$1" == "centos" ]]; then
        podman build -t localhost/centos -f $SCRIPT_DIR/centos.image
    else
        echo "Unknown image: '$1'"
    fi
}

function start(){
    if [[ "$1" != "bluechi" && "$1" != "centos" ]]; then
        echo "Unknown container image: '$1'"
        exit 1
    fi

    port=2020
    for name in ${CONTAINER_NAMES[@]}; do
        # start all containers
        podman run -dt --rm --name $name --network host localhost/$1:latest
        # inject public key
        podman exec $name bash -c "echo $PUBKEY >> ~/.ssh/authorized_keys"
        # change the port for the ssh config
        podman exec $name bash -c "echo 'Port $port' >> /etc/ssh/sshd_config"
        podman exec $name bash -c "systemctl restart sshd"
        let port++
    done
}

function stop() {
    for name in ${CONTAINER_NAMES[@]}; do
        podman stop $name
    done
}

$1 $2
