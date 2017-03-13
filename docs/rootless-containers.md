# Rootless containers in Garden

With the latest release of Garden it is now possible to create and run processes
in containers without requiring root privileges.

## Component Overview

Following is a brief overview of the components required to enable rootless containers
with Garden.

* `gdn` - an all-in-one, standalone version of [Garden-runC](https://github.com/cloudfoundry/garden-runc-release),
the container engine powering [Cloud Foundry](https://www.cloudfoundry.org/) and [Concourse CI](http://concourse.ci/).
* `grootfs` - a daemonless container image manager.
* `runc` - a CLI tool for spawning and running containers according to the OCI specification.

# Getting Started

The following documents the process of installing, configuring and running Garden
as a non-root user on an ubuntu xenial machine.

If you run into any issues or would like any further into, feel free to chat to us on the
`#garden` channel of the [Cloud Foundry Slack](http://slack.cloudfoundry.org/).

## Prerequisites

* An Ubuntu Xenial machine (with kernel version 4.4+)

## Step 1: Install gdn and grootfs

The first step is to download and install `gdn` and `grootfs`. Note that `runc`
does not need to be installed separately as it is bundled together as part of the
`gdn` binary.

NB: All commands in Step 1 need to be run as the root user

```
sudo su
curl "https://raw.githubusercontent.com/cloudfoundry/garden-runc-release/wip-140759953/scripts/install-rootless-gdn" | bash
gdn setup
```

## Step 2: Start the `gdn` server

```
su - rootless
gdn server \
  --image-plugin /usr/local/bin/grootfs \
  --image-plugin-extra-arg=--store \
  --image-plugin-extra-arg=/var/lib/grootfs/btrfs \
  --network-plugin /bin/true \
  --skip-setup
```

As shown above, `gdn` is configurable and extensible via plugins. At the moment `gdn` provides
a plugin interface for image and network management.

## Step 3: Create a container

The `gaol` CLI can be used to interact with the `gdn` server.

```
gaol create -n my-rootless-container -r docker:///busybox
```

## Step 4: Run some processes in the container

```
gaol run my-rootless-container -a -c "echo Hello Rootless :D"
gaol run my-rootless-container -a -c "sh -c 'exit 13'"
```

## Step 5: Destroy the container
