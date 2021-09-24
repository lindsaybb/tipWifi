package tipWifi

import (
	"fmt"
)

// The Interface object is a subset of the Device configuration related
// to its Interfaces, whether Ethernet or Other.
type Interface struct {
	Name         string   `json:"name,omitempty"`          // ex: LAN, "This is a free text field, stating the administrative name of the interface. It may contain spaces and special characters."
	Role         string   `json:"role,omitempty"`          // "The role defines if the interface is upstream or downstream facing." ["upstream","downstream"]
	IsolateHosts int      `json:"isolate-hosts,omitempty"` // def:?  "This option makes sure that any traffic leaving this interface is isolated and all local IP ranges are blocked. It essentially enforces \"guest network\" firewall settings."
	Metric       int      `json:"metric,omitempty"`        // min: 0, max: 4294967295 "The routing metric of this logical interface. Lower values have higher priority."
	Services     []string `json:"services,omitempty"`      // "The services that shall be offered on this logical interface. These are just strings such as \"ssh\", \"lldp\", \"mdns\""
	Vlan         struct { // "This section describes the vlan behaviour of a logical network interface."
		ID    int    `json:"id"`    // max: 4050, "This is the pvid of the vlan that shall be assigned to the interface. The individual physical network devices contained within the interface need to be told explicitly if egress traffic shall be tagged."
		Proto string `json:"proto"` // "The L2 vlan tag that shall be added (1q,1ad)", ["802.1ad","802.1q"],"default": "802.1q"
	} `json:"vlan,omitempty"`
	Bridge struct { // "This section describes the bridge behaviour of a logical network interface."
		Mtu          int `json:"mtu"`           // min: 256, max: 65535, "The MTU that shall be used by the network interface."
		TxQueueLen   int `json:"tx-queue-len"`  // "The Transmit Queue Length is a TCP/IP stack network interface value that sets the number of packets allowed per kernel transmit queue of a network interface device."
		IsolatePorts int `json:"isolate-ports"` // def: false, "Isolates the bridge ports from each other."
	} `json:"bridge,omitempty"`
	Ethernet []struct { // "This section defines the physical copper/fiber ports that are members of the interface. Network devices are referenced by their logical names."
		SelectPorts       []string `json:"select-ports,omitempty"`        // "The list of physical network devices that shall be added to the interface. The names are logical ones and wildcardable. \"WAN\" will use whatever the hardwares default upstream facing port is. \"LANx\" will use the \"x'th\" downstream facing ethernet port. LAN* will use all downstream ports." ["LAN1","LAN2","LAN3","LAN4","LAN*","WAN*","*"]
		Multicast         int      `json:"multicast,omitempty"`           // def: true, "Enable multicast support."
		Learning          int      `json:"learning,omitempty"`            // def: true, "Controls whether a given port will learn MAC addresses from received traffic or not. If learning if off, the bridge will end up flooding any traffic for which it has no FDB entry. By default this flag is on."
		Isolate           int      `json:"isolate,omitempty"`             // def: false, "Only allow communication with non-isolated bridge ports when enabled."
		Macaddr           string   `json:"macaddr,omitempty"`             // "Enforce a specific MAC to these ports."
		ReversePathFilter int      `json:"reverse-path-filter,omitempty"` // def: false, "Reverse Path filtering is a method used by the Linux Kernel to help prevent attacks used by Spoofing IP Addresses."
	} `json:"ethernet,omitempty"`
	Ipv4 struct { // "This section describes the IPv4 properties of a logical interface."
		Addressing   string   `json:"addressing"`        // "This option defines the method by which the IPv4 address of the interface is chosen." ["dynamic","static"]
		Subnet       string   `json:"subnet"`            // "This option defines the static IPv4 of the logical interface in CIDR notation. auto/24 can be used, causing the configuration layer to automatically use any address range from globals.ipv4-network."
		Gateway      string   `json:"gateway"`           // "This option defines the static IPv4 gateway of the logical interface."
		SendHostname int      `json:"send-hostname"`     // def: true, "include the devices hostname inside DHCP requests"
		UseDNS       []string `json:"use-dns,omitempty"` // "Define which DNS servers shall be used. This can either be a list of static IPv4 addresse or dhcp (use the server provided by the DHCP lease)", ["8.8.8.8","4.4.4.4"]
		Dhcp         struct { // "This section describes the DHCP server configuration"
			LeaseFirst      int    `json:"lease-first"`       // "The last octet of the first IPv4 address in this DHCP pool.", ex: 10
			LeaseCount      int    `json:"lease-count"`       // "The number of IPv4 addresses inside the DHCP pool.", ex: 100
			LeaseTime       string `json:"lease-time"`        // def: 6h, "How long the lease is valid before a RENEW must be issued."
			RelayServer     string `json:"relay-server"`      // "Start a L2 DHCP relay in this logical interface and use this IPv4 addr as the upstream server."
			CircuitIDFormat string `json:"circuit-id-format"` // "This option selects what info shall be contained within a relayed frames circuit ID. The string passed in has placeholders that are placed inside a bracket pair \"{}\". Any text not contained within brackets will be included as freetext. Valid placeholders are \"Name, Model, Location, Interface, VLAN-Id, SSID, Crypto, AP-MAC, AP-MAC-Hex, Client-MAC, Client-MAC-Hex\"", ["\\{Interface\\}:\\{VLAN-Id\\}:\\{SSID\\}:\\{Model\\}:\\{Name\\}:\\{AP-MAC\\}:\\{Location\\}","\\{AP-MAC\\};\\{SSID\\};\\{Crypto\\}","\\{Name\\} \\{ESSID\\}"]
			RemoteIDFormat  string `json:"remote-id-format"`  // "This option selects what info shall be contained within a relayed frames remote ID. The string passed in has placeholders that are placed inside a bracket pair \"{}\". Any text not contained within brackets will be included as freetext. Valid placeholders are \"VLAN-Id, SSID, AP-MAC, AP-MAC-Hex, Client-MAC, Client-MAC-Hex\"", ["\\{Client-MAC-hex\\} \\{SSID\\}","\\{AP-MAC-hex\\} \\{SSID\\}"]
		} `json:"dhcp,omitempty"`
		DhcpLeases []struct { // "This section describes the static DHCP leases of this logical interface."
			Macaddr           string `json:"macaddr"`             // "The MAC address of the host that this lease shall be used for."
			StaticLeaseOffset int    `json:"static-lease-offset"` // "The offset of the IP that shall be used in relation to the first IP in the available range."
			LeaseTime         string `json:"lease-time"`          // def: 6h, "How long the lease is valid before a RENEW muss ne issued."
			PublishHostname   int    `json:"publish-hostname"`    // def: true, "Shall the hosts hostname be made available locally via DNS.
		} `json:"dhcp-leases"`
	} `json:"ipv4,omitempty"`
	Ipv6 struct { // "This section describes the IPv6 properties of a logical interface."
		Addressing string `json:"addressing"`  // "This option defines the method by which the IPv6 subnet of the interface is acquired. In static addressing mode, the specified subnet and gateway, if any, are configured on the interface in a fixed manner. Also - if a prefix size hint is specified - a prefix of the given size is allocated from each upstream received prefix delegation pool and assigned to the interface. In dynamic addressing mode, a DHCPv6 client will be launched to obtain IPv6 prefixes for the interface itself and for downstream delegation. Note that dynamic addressing usually only ever makes sense on upstream interfaces." ["dynamic","static"]
		Subnet     string `json:"subnet"`      // "This option defines a static IPv6 prefix in CIDR notation to set on the logical interface. A special notation \"auto/64\" can be used, causing the configuration agent to automatically allocate a suitable prefix from the IPv6 address pool specified in globals.ipv6-network. This property only applies to static addressing mode. Note that this is usually not needed due to DHCPv6-PD assisted prefix assignment."
		Gateway    string `json:"gateway"`     // "This option defines the static IPv6 gateway of the logical interface. It only applies to static addressing mode. Note that this is usually not needed due to DHCPv6-PD assisted prefix assignment."
		PrefixSize int    `json:"prefix-size"` // min: 0, max: 64, "For dynamic addressing interfaces, this property specifies the prefix size to request from an upstream DHCPv6 server through prefix delegation. For static addressing interfaces, it specifies the size of the sub-prefix to allocate from the upstream-received delegation prefixes for assignment to the logical interface."
		Dhcpv6     struct {
			Mode         string   `json:"mode"`          // "Specifies the DHCPv6 server operation mode. When set to \"stateless\", the system will announce router advertisements only, without offering stateful DHCPv6 service. When set to \"stateful\", emitted router advertisements will instruct clients to obtain a DHCPv6 lease. When set to \"hybrid\", clients can freely chose whether to self-assign a random address through SLAAC, whether to request an address via DHCPv6, or both. For maximum compatibility with different clients, it is recommended to use the hybrid mode. The special mode \"relay\" will instruct the unit to act as DHCPv6 relay between this interface and any of the IPv6 interfaces in \"upstream\" mode.", ["hybrid","stateless","stateful","relay"]
			AnnounceDNS  []string `json:"announce-dns"`  // "Overrides the DNS server to announce in DHCPv6 and RA messages. By default, the device will announce its own local interface address as DNS server, essentially acting as proxy for downstream clients. By specifying a non-empty list of IPv6 addresses here, this default behaviour can be overridden."
			FilterPrefix string   `json:"filter-prefix"` // def: ::/0, "Selects a specific downstream prefix or a number of downstream prefix ranges to announce in DHCPv6 and RA messages. By default, all prefixes configured on a given downstream interface are advertised. By specifying an IPv6 prefix in CIDR notation here, only prefixes covered by this CIDR are selected."
		} `json:"dhcpv6"`
	} `json:"ipv6,omitempty"`
	Captive struct { // "This section can be used to setup a captive portal on the AP."
		GatewayName   string `json:"gateway-name"`   // def: "uCentral - Captive Portal", "This name will be presented to connecting users in on the splash page."
		GatewayFqdn   string `json:"gateway-fqdn"`   // def: "ucentral.splash", "The fqdn used for the captive portal IP."
		MaxClients    int    `json:"max-clients"`    // def: 32, "The maximum number of clients that shall be accept."
		UploadRate    int    `json:"upload-rate"`    // def: 0, "The maximum upload rate for a specific client."
		DownloadRate  int    `json:"download-rate"`  // def: 0, "The maximum download rate for a specific client."
		UploadQuota   int    `json:"upload-quota"`   // def: 0, "The maximum upload quota for a specific client."
		DownloadQuota int    `json:"download-quota"` // def: 0, "The maximum download quota for a specific client."
	} `json:"captive,omitempty"`
	Ssids  []*Ssid `json:"ssids"`
	Tunnel []struct {
		Proto       string `json:"proto,omitempty"`        // Schema represents three options: Mesh, VXLAN, GRE
		PeerAddress string `json:"peer-address,omitempty"` // "This is the IP address of the remote host, that the tunnel shall be established with."
		PeerPort    int    `json:"peer-port,omitempty"`    // "The network port that shall be used to establish the tunnel."
		VlanID      int    `json:"vlan-id,omitempty"`      // "This is the id of the vlan that shall be assigned to the interface."
	} `json:"tunnel,omitempty"`
}

// GenerateDescription returns a string of concatenated values describing the Interface object.
func (i *Interface) GenerateDescription() string {
	desc := fmt.Sprintf("Name: %s, ", i.Name)
	desc += fmt.Sprintf("Role: %s, ", i.Role)
	desc += fmt.Sprintf("IPv4 Addressing: %s, ", i.Ipv4.Addressing)
	desc += fmt.Sprintf("IPv6 Addressing: %s, ", i.Ipv6.Addressing)
	for c := 0; c < len(i.Ssids); c++ {
		desc += i.Ssids[c].GenerateDescription()
	}
	return desc
}

// GenerateIpDescription returns a string of concatenated values describing the Interface object's
// IP-related parameters.
// [++]
func (i *Interface) GenerateIpDescription() string {
	desc := fmt.Sprintf("IPv4:")

	return desc
}
