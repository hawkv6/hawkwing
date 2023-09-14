#!/bin/bash

# Killing vpp instances
kill $(pidof vpp)

# Killing dns and webserver instances
kill $(pidof ./dns/dns)
kill $(pidof ./web/webserver)

# Bring down host namespaces
ip netns del ns-host-a
ip netns del ns-host-b
ip netns del ns-host-c
ip netns del ns-dns
ip netns del ns-beyond-ctrl

# Bring down bridges
ip link set dev br-12 down
ip link set dev br-13 down
ip link set dev br-24 down
ip link set dev br-35 down
ip link set dev br-46 down
ip link set dev br-56 down
ip link set dev br-a1 down
ip link set dev br-b6 down
ip link set dev br-c6 down

# Remove bridges
ip link del dev br-12
ip link del dev br-13
ip link del dev br-24
ip link del dev br-35
ip link del dev br-46
ip link del dev br-56
ip link del dev br-a1
ip link del dev br-b6
ip link del dev br-c6

# Remove veth pairs
ip link del dev host-a
ip link del dev host-b
ip link del dev host-c
ip link del dev dns