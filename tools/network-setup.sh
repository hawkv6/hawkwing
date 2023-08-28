# Start VPP instances
$VPP_BINARY_PATH api-segment { prefix vpp1 } socksvr { socket-name /run/vpp/api-vpp1.sock } cpu {main-core 3} unix { cli-listen /run/vpp/cli.vpp1.sock  cli-prompt vpp1# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp1.sock}
$VPP_BINARY_PATH api-segment { prefix vpp2 } socksvr { socket-name /run/vpp/api-vpp2.sock } cpu {main-core 4} unix { cli-listen /run/vpp/cli.vpp2.sock  cli-prompt vpp2# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp2.sock}
$VPP_BINARY_PATH api-segment { prefix vpp3 } socksvr { socket-name /run/vpp/api-vpp3.sock } cpu {main-core 5} unix { cli-listen /run/vpp/cli.vpp3.sock  cli-prompt vpp3# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp3.sock}
$VPP_BINARY_PATH api-segment { prefix vpp4 } socksvr { socket-name /run/vpp/api-vpp4.sock } cpu {main-core 6} unix { cli-listen /run/vpp/cli.vpp4.sock  cli-prompt vpp4# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp4.sock}
$VPP_BINARY_PATH api-segment { prefix vpp5 } socksvr { socket-name /run/vpp/api-vpp5.sock } cpu {main-core 7} unix { cli-listen /run/vpp/cli.vpp5.sock  cli-prompt vpp5# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp5.sock}
$VPP_BINARY_PATH api-segment { prefix vpp6 } socksvr { socket-name /run/vpp/api-vpp6.sock } cpu {main-core 8} unix { cli-listen /run/vpp/cli.vpp6.sock  cli-prompt vpp6# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-vpp6.sock}
$VPP_BINARY_PATH api-segment { prefix site-a } socksvr { socket-name /run/vpp/api-site-a.sock } cpu {main-core 9} unix { cli-listen /run/vpp/cli.site-a.sock  cli-prompt site-a# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-site-a.sock}
$VPP_BINARY_PATH api-segment { prefix site-b } socksvr { socket-name /run/vpp/api-site-b.sock } cpu {main-core 10} unix { cli-listen /run/vpp/cli.site-b.sock  cli-prompt site-b# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-site-b.sock}
$VPP_BINARY_PATH api-segment { prefix site-c } socksvr { socket-name /run/vpp/api-site-c.sock } cpu {main-core 11} unix { cli-listen /run/vpp/cli.site-c.sock  cli-prompt site-c# } plugins { plugin default { disable } plugin af_packet_plugin.so { enable } plugin ping_plugin.so { enable } } statseg { socket-name /run/vpp/stats-site-c.sock}
sleep 5

# Create links between Ps and PEs
ip link add dev br-12 type bridge
ip link set dev br-12 up
ip link add dev br-13 type bridge
ip link set dev br-13 up
ip link add dev br-24 type bridge
ip link set dev br-24 up
ip link add dev br-35 type bridge
ip link set dev br-35 up
ip link add dev br-46 type bridge
ip link set dev br-46 up
ip link add dev br-56 type bridge
ip link set dev br-56 up
ip link add dev br-a1 type bridge
ip link set dev br-a1 up
ip link add dev br-b6 type bridge
ip link set dev br-b6 up
ip link add dev br-c6 type bridge
ip link set dev br-c6 up

# Create veth pair between site routers and hosts
ip link add host-a type veth peer name client-site-a
ip link set dev host-a up
ip link set dev client-site-a up
ip link add dns type veth peer name dns-site-a
ip link set dev dns up
ip link set dev dns-site-a up
ip link add host-b type veth peer name client-site-b
ip link set dev host-b up
ip link set dev client-site-b up
ip link add host-c type veth peer name client-site-c
ip link set dev host-c up
ip link set dev client-site-c up

# Create veth pair between vpp4 and beyond-ctrl
ip link add beyond-ctrl type veth peer name client-vpp4
ip link set dev beyond-ctrl up
ip link set dev client-vpp4 up

sleep 5

# Create host-a
ip netns add ns-host-a
ip link set host-a netns ns-host-a
ip netns exec ns-host-a ip link set host-a up
ip netns exec ns-host-a ip -6 address add 2001:cafe:a::a/64 dev host-a
ip netns exec ns-host-a ip -6 address add fcbb:cc00:1::a/48 dev host-a
ip netns exec ns-host-a ip -6 route add fcbb:bb00::/32 via 2001:cafe:a::1 dev host-a metric 1 # route to reach internal vpp network
ip netns exec ns-host-a ip -6 route add fcbb:aa00::/32 via 2001:cafe:a::1 dev host-a metric 1 # route to reach site-X routers
ip netns exec ns-host-a ip -6 route add fcbb:cc00:2::/48 via 2001:cafe:a::1 dev host-a metric 1 # route to reach host-b
ip netns exec ns-host-a ip -6 route add fcbb:cc00:3::/48 via 2001:cafe:a::1 dev host-a metric 1 # route to reach host-c
ip netns exec ns-host-a ip -6 route add fcbb:cc00:4::/48 via 2001:cafe:a::1 dev host-a metric 1 # route to reach dns
ip netns exec ns-host-a ip -6 route add fcbb:cc00:5::/48 via 2001:cafe:a::1 dev host-a metric 1 # route to reach beyond-ctrl

# Create host-b
ip netns add ns-host-b
ip link set host-b netns ns-host-b
ip netns exec ns-host-b ip link set host-b up
ip netns exec ns-host-b ip -6 address add 2001:cafe:b::b/64 dev host-b
ip netns exec ns-host-b ip -6 address add fcbb:cc00:2::a/48 dev host-b
ip netns exec ns-host-b ip -6 route add fcbb:bb00::/32 via 2001:cafe:b::1 dev host-b metric 1 # route to reach internal vpp network
ip netns exec ns-host-b ip -6 route add fcbb:aa00::/32 via 2001:cafe:b::1 dev host-b metric 1 # route to reach site-X routers
ip netns exec ns-host-b ip -6 route add fcbb:cc00:1::/48 via 2001:cafe:b::1 dev host-b metric 1 # route to reach host-a
ip netns exec ns-host-b ip -6 route add fcbb:cc00:3::/48 via 2001:cafe:b::1 dev host-b metric 1 # route to reach host-c
ip netns exec ns-host-b ip -6 route add fcbb:cc00:5::/48 via 2001:cafe:b::1 dev host-b metric 1 # route to reach beyond-ctrl

# Create host-c
ip netns add ns-host-c
ip link set host-c netns ns-host-c
ip netns exec ns-host-c ip link set host-c up
ip netns exec ns-host-c ip -6 address add 2001:cafe:c::c/64 dev host-c
ip netns exec ns-host-c ip -6 address add fcbb:cc00:3::a/48 dev host-c
ip netns exec ns-host-c ip -6 route add fcbb:bb00::/32 via 2001:cafe:c::1 dev host-c metric 1 # route to reach internal vpp network
ip netns exec ns-host-c ip -6 route add fcbb:aa00::/32 via 2001:cafe:c::1 dev host-c metric 1 # route to reach site-X routers
ip netns exec ns-host-c ip -6 route add fcbb:cc00:1::/48 via 2001:cafe:c::1 dev host-c metric 1 # route to reach host-a
ip netns exec ns-host-c ip -6 route add fcbb:cc00:2::/48 via 2001:cafe:c::1 dev host-c metric 1 # route to reach host-b
ip netns exec ns-host-c ip -6 route add fcbb:cc00:5::/48 via 2001:cafe:c::1 dev host-c metric 1 # route to reach beyond-ctrl

# Create dns
ip netns add ns-dns
ip link set dns netns ns-dns
ip netns exec ns-dns ip link set dns up
ip netns exec ns-dns ip -6 address add 2001:cafe:f::f/64 dev dns
ip netns exec ns-dns ip -6 address add fcbb:cc00:4::f/48 dev dns
ip netns exec ns-dns ip -6 route add fcbb:bb00::/32 via 2001:cafe:f::1 dev dns metric 1 # route to reach internal vpp network
ip netns exec ns-dns ip -6 route add fcbb:aa00:1::/48 via 2001:cafe:f::1 dev dns metric 1 # route to reach site-a
ip netns exec ns-dns ip -6 route add fcbb:cc00:1::/48 via 2001:cafe:f::1 dev dns metric 1 # route to reach host-a

# Create beyond-ctrl
ip netns add ns-beyond-ctrl
ip link set beyond-ctrl netns ns-beyond-ctrl
ip netns exec ns-beyond-ctrl ip link set beyond-ctrl up
ip netns exec ns-beyond-ctrl ip -6 address add 2001:db8:f4::f/64 dev beyond-ctrl
ip netns exec ns-beyond-ctrl ip -6 address add fcbb:cc00:5::f/48 dev beyond-ctrl
ip netns exec ns-beyond-ctrl ip -6 route add fcbb:bb00::/32 via 2001:db8:f4::4 dev beyond-ctrl metric 1 # route to reach internal vpp network
ip netns exec ns-beyond-ctrl ip -6 route add fcbb:aa00::/32 via 2001:db8:f4::4 dev beyond-ctrl metric 1 # route to reach site-X routers
ip netns exec ns-beyond-ctrl ip -6 route add fcbb:cc00:1::/48 via 2001:db8:f4::4 dev beyond-ctrl metric 1 # route to reach host-a
ip netns exec ns-beyond-ctrl ip -6 route add fcbb:cc00:2::/48 via 2001:db8:f4::4 dev beyond-ctrl metric 1 # route to reach host-b
ip netns exec ns-beyond-ctrl ip -6 route add fcbb:cc00:3::/48 via 2001:db8:f4::4 dev beyond-ctrl metric 1 # route to reach host-c

### Configure vpp1
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock loopback create-interface 
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface ip address loop0 fcbb:bb00:1::1/128
# tap10 br-12 (to vpp2)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock create tap id 10 host-bridge br-12
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface state tap10 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock enable ip6 interface tap10
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface ip address tap10 2001:db8:12::1/64
# tap11 br-13 (to vpp3)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock create tap id 11 host-bridge br-13
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface state tap11 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock enable ip6 interface tap11
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface ip address tap11 2001:db8:13::1/64
# tap12 br-a1 (to site-a)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock create tap id 12 host-bridge br-a1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface state tap12 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock enable ip6 interface tap12
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock set interface ip address tap12 2001:db8:a1::1/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:2::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:3::/48 via 2001:db8:13::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:4::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:5::/48 via 2001:db8:13::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:6::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:bb00:6::/48 via 2001:db8:13::3
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:aa00:1::/48 via 2001:db8:a1::a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:1::/48 via 2001:db8:a1::a
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:2::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:2::/48 via 2001:db8:13::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:aa00:2::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:aa00:2::/48 via 2001:db8:13::3
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:3::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:3::/48 via 2001:db8:13::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:aa00:3::/48 via 2001:db8:12::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:aa00:3::/48 via 2001:db8:13::3
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp1.sock ip route add fcbb:cc00:5::/48 via 2001:db8:12::2

### Configure vpp2
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface ip address loop0 fcbb:bb00:2::1/128
# tap20 br-12 (to vpp1)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock create tap id 20 host-bridge br-12
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface state tap20 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock enable ip6 interface tap20
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface ip address tap20 2001:db8:12::2/64
# tap21 br-24 (to vpp4)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock create tap id 21 host-bridge br-24
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface state tap21 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock enable ip6 interface tap21
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock set interface ip address tap21 2001:db8:24::2/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:1::/48 via 2001:db8:12::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:3::/48 via 2001:db8:12::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:4::/48 via 2001:db8:24::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:5::/48 via 2001:db8:12::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:5::/48 via 2001:db8:24::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:bb00:6::/48 via 2001:db8:24::4
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:cc00:1::/48 via 2001:db8:12::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:aa00:1::/48 via 2001:db8:12::1
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:cc00:2::/48 via 2001:db8:24::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:aa00:2::/48 via 2001:db8:24::4
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:cc00:3::/48 via 2001:db8:24::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:aa00:3::/48 via 2001:db8:24::4
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp2.sock ip route add fcbb:cc00:5::/48 via 2001:db8:24::4

### Configure vpp3
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface ip address loop0 fcbb:bb00:3::1/128
# tap30 br-13 (to vpp1)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock create tap id 30 host-bridge br-13
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface state tap30 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock enable ip6 interface tap30
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface ip address tap30 2001:db8:13::3/64
# tap31 br-35 (to vpp5)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock create tap id 31 host-bridge br-35
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface state tap31 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock enable ip6 interface tap31
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock set interface ip address tap31 2001:db8:35::3/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:1::/48 via 2001:db8:13::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:2::/48 via 2001:db8:13::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:4::/48 via 2001:db8:13::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:4::/48 via 2001:db8:35::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:5::/48 via 2001:db8:35::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:bb00:6::/48 via 2001:db8:35::5
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:cc00:1::/48 via 2001:db8:13::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:aa00:1::/48 via 2001:db8:13::1
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:cc00:2::/48 via 2001:db8:35::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:aa00:2::/48 via 2001:db8:35::5
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:cc00:3::/48 via 2001:db8:35::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:aa00:3::/48 via 2001:db8:35::5
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:cc00:5::/48 via 2001:db8:13::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp3.sock ip route add fcbb:cc00:5::/48 via 2001:db8:35::5

### Configure vpp4
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface ip address loop0 fcbb:bb00:4::1/128
# tap40 br-24 (to vpp2)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock create tap id 40 host-bridge br-24
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface state tap40 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock enable ip6 interface tap40
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface ip address tap40 2001:db8:24::4/64
# tap41 br-46 (to vpp6)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock create tap id 41 host-bridge br-46
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface state tap41 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock enable ip6 interface tap41
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface ip address tap41 2001:db8:46::4/64
# (to beyond-ctrl) # TODO check if that is correct
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock create host-interface name client-vpp4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface state host-client-vpp4 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock enable ip6 interface host-client-vpp4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock set interface ip address host-client-vpp4 2001:db8:f4::4/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:1::/48 via 2001:db8:24::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:2::/48 via 2001:db8:24::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:3::/48 via 2001:db8:24::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:3::/48 via 2001:db8:46::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:5::/48 via 2001:db8:46::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:bb00:6::/48 via 2001:db8:46::6
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:cc00:1::/48 via 2001:db8:24::2
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:aa00:1::/48 via 2001:db8:24::2
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:cc00:2::/48 via 2001:db8:46::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:aa00:2::/48 via 2001:db8:46::6
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:cc00:3::/48 via 2001:db8:46::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:aa00:3::/48 via 2001:db8:46::6
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp4.sock ip route add fcbb:cc00:5::/48 via 2001:db8:f4::f

### Configure vpp5
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface ip address loop0 fcbb:bb00:5::1/128
# tap50 br-35 (to vpp3)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock create tap id 50 host-bridge br-35
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface state tap50 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock enable ip6 interface tap50
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface ip address tap50 2001:db8:35::5/64
# tap51 br-56 (to vpp6)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock create tap id 51 host-bridge br-56
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface state tap51 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock enable ip6 interface tap51
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock set interface ip address tap51 2001:db8:56::5/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:1::/48 via 2001:db8:35::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:2::/48 via 2001:db8:35::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:2::/48 via 2001:db8:56::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:3::/48 via 2001:db8:35::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:4::/48 via 2001:db8:56::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:bb00:6::/48 via 2001:db8:56::6
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:cc00:1::/48 via 2001:db8:35::3
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:aa00:1::/48 via 2001:db8:35::3
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:cc00:2::/48 via 2001:db8:56::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:aa00:2::/48 via 2001:db8:56::6
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:cc00:3::/48 via 2001:db8:56::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:aa00:3::/48 via 2001:db8:56::6
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp5.sock ip route add fcbb:cc00:5::/48 via 2001:db8:56::6

### Configure vpp6
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface ip address loop0 fcbb:bb00:6::1/128
# tap60 br-46 (to vpp4)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock create tap id 60 host-bridge br-46
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface state tap60 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock enable ip6 interface tap60
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface ip address tap60 2001:db8:46::6/64
# tap61 br-56 (to vpp5)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock create tap id 61 host-bridge br-56
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface state tap61 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock enable ip6 interface tap61
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface ip address tap61 2001:db8:56::6/64
# tap62 br-b6 (to site-b)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock create tap id 62 host-bridge br-b6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface state tap62 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock enable ip6 interface tap62
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface ip address tap62 2001:db8:b6::6/64
# tap63 br-c6 (to site-c)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock create tap id 63 host-bridge br-c6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface state tap63 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock enable ip6 interface tap63
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock set interface ip address tap63 2001:db8:c6::6/64
# static routing vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:1::/48 via 2001:db8:46::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:1::/48 via 2001:db8:56::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:2::/48 via 2001:db8:46::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:3::/48 via 2001:db8:56::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:4::/48 via 2001:db8:46::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:bb00:5::/48 via 2001:db8:56::5
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:cc00:1::/48 via 2001:db8:46::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:cc00:1::/48 via 2001:db8:56::5
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:aa00:1::/48 via 2001:db8:46::4
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:aa00:1::/48 via 2001:db8:56::5
# static routing to site-b # TODO vrfs and routing with SRv6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:cc00:2::/48 via 2001:db8:b6::b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:aa00:2::/48 via 2001:db8:b6::b
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:cc00:3::/48 via 2001:db8:c6::c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:aa00:3::/48 via 2001:db8:c6::c
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.vpp6.sock ip route add fcbb:cc00:5::/48 via 2001:db8:46::4

### Configure site-a
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface ip address loop0 fcbb:aa00:1::1/128
# tap70 br-a1 (to vpp1)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock create tap id 70 host-bridge br-a1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface state tap70 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock enable ip6 interface tap70
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface ip address tap70 2001:db8:a1::a/64
# client-site-a (to host-a)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock create host-interface name client-site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface state host-client-site-a up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock enable ip6 interface host-client-site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface ip address host-client-site-a 2001:cafe:a::1/64
# dns-site-a (to dns)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock create host-interface name dns-site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface state host-dns-site-a up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock enable ip6 interface host-dns-site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock set interface ip address host-dns-site-a 2001:cafe:f::1/64
# static routing to dns
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:cc00:4::/48 via 2001:cafe:f::f
# static routing to host-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:cc00:1::/48 via 2001:cafe:a::a
# static routing to vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:bb00::/32 via 2001:db8:a1::1
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:aa00:2::/48 via 2001:db8:a1::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:cc00:2::/48 via 2001:db8:a1::1
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:aa00:3::/48 via 2001:db8:a1::1
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:cc00:3::/48 via 2001:db8:a1::1
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-a.sock ip route add fcbb:cc00:5::/48 via 2001:db8:a1::1

### Configure site-b
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface ip address loop0 fcbb:aa00:2::1/128
# tap80 br-b6 (to vpp6)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock create tap id 80 host-bridge br-b6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface state tap80 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock enable ip6 interface tap80
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface ip address tap80 2001:db8:b6::b/64
# client-site-b (to host-b)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock create host-interface name client-site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface state host-client-site-b up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock enable ip6 interface host-client-site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock set interface ip address host-client-site-b 2001:cafe:b::1/64
# static routing to vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:bb00::/32 via 2001:db8:b6::6
# static routing to host-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:cc00:2::/48 via 2001:cafe:b::b
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:cc00:1::/48 via 2001:db8:b6::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:aa00:1::/48 via 2001:db8:b6::6
# static routing to site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:cc00:3::/48 via 2001:db8:b6::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:aa00:3::/48 via 2001:db8:b6::6
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-b.sock ip route add fcbb:cc00:5::/48 via 2001:db8:b6::6

### Configure site-c
# loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock loopback create-interface
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface state loop0 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock enable ip6 interface loop0
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface ip address loop0 fcbb:aa00:3::1/128
# tap90 br-c6 (to vpp6)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock create tap id 90 host-bridge br-c6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface state tap90 up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock enable ip6 interface tap90
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface ip address tap90 2001:db8:c6::c/64
# client-site-c (to host-c)
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock create host-interface name client-site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface state host-client-site-c up
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock enable ip6 interface host-client-site-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock set interface ip address host-client-site-c 2001:cafe:c::1/64
# static routing to vpp network
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:bb00::/32 via 2001:db8:c6::6
# static routing to host-c
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:cc00:3::/48 via 2001:cafe:c::c
# static routing to site-a
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:cc00:1::/48 via 2001:db8:c6::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:aa00:1::/48 via 2001:db8:c6::6
# static routing to site-b
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:cc00:2::/48 via 2001:db8:c6::6
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:aa00:2::/48 via 2001:db8:c6::6
# static routing to beyond-ctrl
$VPPCTL_BINARY_PATH -s /run/vpp/cli.site-c.sock ip route add fcbb:cc00:5::/48 via 2001:db8:c6::6

sleep 10

# Ping to startup network & arp
ip netns exec ns-host-a ping -c 5 fcbb:cc00:2::a &
ip netns exec ns-host-a ping -c 5 fcbb:cc00:3::a &
ip netns exec ns-host-a ping -c 5 fcbb:cc00:4::f &
ip netns exec ns-host-a ping -c 5 fcbb:cc00:5::f 

sleep 5