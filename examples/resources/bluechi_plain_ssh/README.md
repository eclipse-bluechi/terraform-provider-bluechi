This example requires a container setup as provided by [container-setup.sh](../../../container/container-setup.sh):

```bash
# build image with bluechi pre-installed or
$ bash container/container-setup.sh build_image bluechi
# use a plain centos container without bluechi
$ bash container/container-setup.sh build_image centos

# start all container, bluechi isn't setup, yet
$ bash container/container-setup.sh start bluechi

# build and install terraform-provider-bluechi
$ make install

# apply the terraform example
$ cd examples/resources/bluechi_plain_ssh/
$ tf init
$ tf apply
```
