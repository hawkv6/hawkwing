#!/bin/bash

# Define colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display help/usage message
usage() {
    echo -e "${YELLOW}Usage:${NC} $0 [-h] [-c] [-s] [-b] [-g] [-w]"
    echo "  -h                  Display this help message"
    echo "  -c <domain:port>    Start client"
    echo "  -s <server>         Start server"
    echo "  -b                  Start bpf client"
    echo "  -g                  Start controller"
    echo "  -w <bridge>         Start wireshark on bridge"
    exit 1
}

start_tcp_client() {
    local ARGS=(${OPTARG//:/ })
    local DOMAIN=${ARGS[0]}
    local PORT=${ARGS[1]}
    echo -e "${GREEN}Connecting to server ${DOMAIN} and port ${PORT}${NC}"
    ssh ins@172.16.16.28 "sudo ip netns exec ns-host-a /home/ins/hawkwing/tools/demo/demo client -H ${DOMAIN} -p ${PORT}"
}

start_wireshark() {
    local BRIDGE=$1
    echo -e "${GREEN}Starting wireshark on bridge ${BRIDGE}${NC}"
    ssh ins@172.16.16.28 "sudo tcpdump -U -nni br-${BRIDGE} -w -" | /Applications/Wireshark.app/Contents/MacOS/Wireshark -k -i -
}

start_bpf_client() {
    echo -e "${GREEN}Starting bpf client${NC}"
    ssh ins@172.16.16.28 "cd hawkwing && make start-client"
}

start_bpf_server() {
    local SERVER=$1
    echo -e "${GREEN}Starting bpf server${NC}"
    ssh ins@172.16.16.28 "cd hawkwing && make start-server_${SERVER}"
}

start_controller() {
    echo -e "${GREEN}Starting controller${NC}"
    ssh ins@172.16.16.28 "cd hawkwing && make start-controller"
}

# Parse arguments
while getopts ":hc:s:bgw:" opt; do
    case $opt in
        c)
            start_tcp_client $OPTARG
            ;;
        s)
            start_bpf_server $OPTARG
            ;;
        b)
            start_bpf_client
            ;;
        g) 
            start_controller
            ;;
        w)
            start_wireshark $OPTARG
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
