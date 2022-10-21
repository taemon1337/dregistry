# dregistry
A containerd distributed image registry

## Overview

The design and purpose of this image registry is to manage container images on all nodes in your cluster in the most efficient way possible.

## Why?



## How?

This is achieved by installing a `dregistry` as a local service on each cluster node and setting it as the default image registry.  DRegistry will interact with the `containerd:image:API` to List and Get container images and layers.  The Gossip protocol is used to send events across the cluster to discover if any other node has downloaded an Image Layer needed and if so, the node will fetch from the node which already has it.  If no node is found with that image, the node itself (or a special egress node) will download it and provide it to the cluster.

