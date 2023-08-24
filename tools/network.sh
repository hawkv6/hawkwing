#!/bin/bash

# Define colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check for sudo privileges
if [[ $EUID -ne 0 ]]; then
    echo -e "${RED}This script must be run with sudo privileges. Please use 'sudo' to execute the script.${NC}"
    exit 1
fi

# Function to display help/usage message
usage() {
    echo -e "${YELLOW}Usage:${NC} $0 [-h] [-s] [-c]"
    echo "  -h        Display this help message"
    echo "  -s        Start the network"
    echo "  -c        Clean the network"
    exit 1
}

# Function to start the network
start_network() {
    echo -e "${GREEN}Starting the network...${NC}"
    # Add your commands to start the network here
    # For example:
    # docker-compose up -d
}

# Function to clean the network
clean_network() {
    echo -e "${RED}Cleaning the network...${NC}"
    # Add your commands to clean the network here
    # For example:
    # docker-compose down

    echo -e "Killing all VPP instances"
    # kill $(pidof vpp)

    echo -e "Delete host namespaces"
    # ip netns del host1

    echo -e "Delete bridges"
    # ip link del br1

    echo -e "Delete veth pairs"
    # ip link del veth1

    echo -e "${GREEN}Network cleaned successfully${NC}"
}

# Parse arguments
while getopts ":hsc" opt; do
    case $opt in
        s)
            start_network
            ;;
        c)
            clean_network
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
