---
description:  How to run a Wasp node in using Docker. Build the image, configure it, run it. 
image: /img/logo/WASP_logo_dark.png
keywords:
- ISCP
- Smart Contracts
- Running a node
- docker
- image
- build
- configure
- arguments
---
# Docker

This page describes the configuration of the Wasp node in combination with Docker. If you followed the instructions in [Running a Node](../guide/chains_and_nodes/running-a-node.md), you can skip to [Configuring wasp-cli](../guide/chains_and_nodes/wasp-cli.md).

## Introduction

The Dockerfile is separated into several stages which effectively splits Wasp into four small pieces:

* Testing
    * Unit testing
    * Integration testing
* Wasp CLI
* Wasp Node

## Running a Wasp Node

Checkout the project, switch to 'develop' and build the main image:

```
git clone -b develop https://github.com/iotaledger/wasp.git
cd wasp
docker build -t wasp-node .
```

The build process will copy the docker_config.json file into the image, which will use it when the node gets started. 

By default, the build process will use `-tags rocksdb,builtin_static` as a build argument. This argument can be modified with `--build-arg BUILD_TAGS=<tags>`.

Depending on the use case, Wasp requires a different GoShimmer hostname which can be changed at this part inside the [docker_config.json](https://github.com/iotaledger/wasp/blob/develop/docker_config.json) file:

```
  "nodeconn": {
    "address": "goshimmer:5000"
  },
```

After the build process has finished, you can start your Wasp node by running:

```
docker run wasp-node
```

### Configuration

After the build process has been completed, it is still possible to inject a different configuration file into a new container by running: 

```
docker run -v $(pwd)/alternative_docker_config.json:/etc/wasp_config.json wasp-node
```

You can also add further configuration using arguments:

```
docker run wasp-node --nodeconn.address=alt_goshimmer:5000 
```

To get a list of all available arguments, run the node with the argument '--help'

```
docker run wasp-node --help
```
