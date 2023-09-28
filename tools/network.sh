#!/bin/bash

# Define colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Define VPP variables
export VPP_BINARY_PATH=/home/ins/vpp/build-root/build-vpp-native/vpp/bin/vpp
export VPPCTL_BINARY_PATH=/home/ins/vpp/build-root/build-vpp-native/vpp/bin/vppctl

# Check for sudo privileges
if [[ $EUID -ne 0 ]]; then
    echo -e "${RED}This script must be run with sudo privileges. Please use 'sudo' to execute the script.${NC}"
    exit 1
fi

# Function to display help/usage message
usage() {
    echo -e "${YELLOW}Usage:${NC} $0 [-h] [-s] [-c] [-i] [-n] [-p]"
    echo "  -h                  Display this help message"
    echo "  -s                  Start the network"
    echo "  -c                  Clean the network"
    echo "  -i [process]        Interact with VPP"
    echo "  -n [namespace]      Open a shell in the specified namespace"
    echo "  -p                  Start the client application"
    echo "  -q                  Start the server application"
    exit 1
}

# Function to start the network
start_network() {
    echo -e "${GREEN}Starting the network...${NC}"
    bash network-setup.sh
    echo -e "${GREEN}Network started successfully${NC}"
    start_dns_server
    start_webservers
}

# Function to clean the network
clean_network() {
    echo -e "${RED}Cleaning the network...${NC}"
    bash network-clean.sh
    echo -e "${GREEN}Network cleaned successfully${NC}"
}

vpp_interaction() {
    local PROCESS=$1
    CMD="$VPPCTL_BINARY_PATH -s /run/vpp/cli.$PROCESS.sock"
    echo -e "${GREEN}Connecting to ${PROCESS}...${NC}"
    $CMD
}

open_namespace_shell() {
    local NS=$1
    echo -e "${GREEN}Entering namespace ${NS}...${NC}"
    ip netns exec "$NS" /bin/bash
}

start_dns_server() {
    echo -e "${GREEN}Starting DNS server...${NC}"
    ip netns exec ns-dns ./dns/dns &
    echo -e "${GREEN}DNS server started successfully${NC}"
}

start_webservers() {
    echo -e "${GREEN}Starting webservers...${NC}"
    ip netns exec ns-host-b ./web/webserver host-b 80 &
    ip netns exec ns-host-c ./web/webserver host-c 8080 &
    echo -e "${GREEN}Webservers started successfully${NC}"
}

start_client_application() {
    echo -e "${GREEN}Starting application...${NC}"
    mount -t bpf bpf /sys/fs/bpf
    cd .. && ./out/bin/hawkwing client --config ./test_assets/config.yaml
}

start_server_application() {
    echo -e "${GREEN}Starting application...${NC}"
    mount -t bpf bpf /sys/fs/bpf
    cd .. && ./out/bin/hawkwing server --config ./test_assets/config.yaml
}

# Parse arguments
while getopts ":hsci:n:pq" opt; do
    case $opt in
        s)
            start_network
            ;;
        c)
            clean_network
            ;;
        i)
            vpp_interaction $OPTARG
            ;;
        n)
            open_namespace_shell $OPTARG
            ;;
        p)
            start_client_application
            ;;
        q)
            start_server_application
            ;;
        h)
            usage
            ;;
        \?)
            echo -e "${RED}Invalid option:${NC} -$OPTARG" >&2
            usage
            ;;
    esac
done
