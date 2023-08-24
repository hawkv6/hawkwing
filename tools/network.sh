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
    echo -e "${YELLOW}Usage:${NC} $0 [-h] [-s] [-c]"
    echo "  -h                  Display this help message"
    echo "  -s                  Start the network"
    echo "  -c                  Clean the network"
    echo "  -i [process]        Interact with VPP"
    echo "  -n [namespace]      Open a shell in the specified namespace"
    exit 1
}

# Function to start the network
start_network() {
    echo -e "${GREEN}Starting the network...${NC}"
    bash network-setup.sh
    echo -e "${GREEN}Network started successfully${NC}"
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
    $CMD
}

open_namespace_shell() {
    local NS=$1
    echo -e "${GREEN}Entering namespace ${NS}...${NC}"
    ip netns exec "$NS" /bin/bash
}

# Parse arguments
while getopts ":hsci:n:" opt; do
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
        h)
            usage
            ;;
        \?)
            echo -e "${RED}Invalid option:${NC} -$OPTARG" >&2
            usage
            ;;
    esac
done
