package tipWifi

import (
	"fmt"
)

// The Devices object contains a list of the GW's Device object.
type Devices struct {
	Entry []*Device `json:"devices"`
}

// The FirmwareDevices object contains a list of the FMS' FirmwareDevice object.
type FirmwareDevices struct {
	Entry []*FirmwareDevice `json:"devices"`
}

// GenerateStatusReport returns a list of Serial Numbers prefixed with
// whether they are "UP" or "DOWN".
func (fwds *FirmwareDevices) GenerateStatusReport() (list []string) {
	for i := 0; i < len(fwds.Entry); i++ {
		if fwds.Entry[i].IsConnected() {
			list = append(list, fmt.Sprintf("UP:%s", fwds.Entry[i].SerialNumber))
		} else {
			list = append(list, fmt.Sprintf("DOWN:%s", fwds.Entry[i].SerialNumber))
		}
	}
	return list
}

// The DeviceInfo list contains valid qualifiers when asking for specific information
// from a Device object.
var DeviceInfo = []string{
	"configuration", //
	"interfaces",
	"capabilities", //
	"status",       //
	"stats",        //
	"logs",         //
	"health",       //
}

// The DeviceInfo list contains valid commands that can be applied to a Device object
var DeviceActions = []string{
	"configure", //
	"request",   //
	"reboot",
	"factory", //
	"upgrade",
	"rtty",          //
	"getfile",       //
	"geteventqueue", //
	"deletelogs",    //
	"trace",         //
	"wifiscan",      //
	"toggleleds",    //
}

// GenerateDeviceReport is a print wrapper around the ListInfoBySn function.
func (devs *Devices) GenerateDeviceReport(info string) {
	sns := devs.ListSns()
	for n := 0; n < len(sns); n++ {
		infoList := devs.ListInfoBySn(sns[n], info)
		DisplayList(sns[n], infoList)
	}
}

// ListSns returns a list of SerialNumber strings configured on the uc.GW.
func (devs *Devices) ListSns() (list []string) {
	for i := 0; i < len(devs.Entry); i++ {
		list = append(list, devs.Entry[i].SerialNumber)
	}
	return list
}

// ListInfoBySn is a wrapper around the individual Device object's ListInfo
// function that only returns the Info list from a device matching the
// supplied Serial Number.
func (devs *Devices) ListInfoBySn(sn, info string) []string {
	for i := 0; i < len(devs.Entry); i++ {
		if sn == devs.Entry[i].SerialNumber {
			return devs.Entry[i].ListInfo(info)
		}
	}
	return []string{}
}

// The Device object represents the complete uc.GW data model including configuration.
type Device struct {
	UUID          int    `json:"UUID"`
	Compatible    string `json:"compatible"`
	Configuration struct {
		UUID int      `json:"uuid"` // "The unique ID of the configuration. This is the unix timestamp of when the config was created."
		Unit struct { // "A device has certain properties that describe its identity and location. These properties are described inside this object."
			Name      string `json:"name,omitempty"`       // "This is a free text field, stating the administrative name of the device. It may contain spaces and special characters."
			Location  string `json:"location,omitempty"`   // "This is a free text field, stating the location of the device. It may contain spaces and special characters."
			Timezone  string `json:"timezone,omitempty"`   // "This allows you to change the TZ of the device." ["UTC","EST5","CET-1CEST,M3.5.0,M10.5.0/3"]
			LedActive int    `json:"led-active,omitempty"` // def: true, "This allows forcing all LEDs off."
		} `json:"unit,omitempty"`
		Globals struct {
			Ipv4Network string `json:"ipv4-network"` // "Define the IPv4 range that is delegatable to the downstream interfaces This is described as a CIDR block. (192.168.0.0/16, 172.16.128/17)"
			Ipv6Network string `json:"ipv6-network"` // "Define the IPv6 range that is delegatable to the downstream interfaces This is described as a CIDR block. (fdca:1234:4567::/48)"
		} `json:"globals,omitempty"`
		Radios     []*Radio     `json:"radios"`
		Interfaces []*Interface `json:"interfaces"`
		Services   struct {     // "This section describes all of the services that may be present on the AP. Each service is then referenced via its name inside an interface, ssid, ..."
			Lldp struct { //
				Describe string `json:"describe"` // def: "uCentral Access Point", "The LLDP description field. If set to \"auto\" it will be derived from unit.name."
				Location string `json:"location"` // def: "uCentral Network", "The LLDP location field. If set to \"auto\" it will be derived from unit.location."
			} `json:"lldp,omitempty"`
			SSH struct { // "This section can be used to setup a SSH server on the AP."
				Port                   int    `json:"port"`                    // def: 22, max: 65535, "This option defines which port the SSH server shall be available on."
				AuthoirzedKeys         string `json:"authorized-keys"`         // "This allows the upload of public ssh keys. Keys need to be seperated by a newline."
				PasswordAuthentication int    `json:"password-authentication"` // def: true, "This option defines if password authentication shall be enabled. If set to false, only ssh key based authentication is possible."
			} `json:"ssh,omitempty"`
			Ntp struct { // "This section can be used to setup the upstream NTP servers."
				Servers     []string `json:"servers"`      // "This is an array of URL/IP of the upstream NTP servers that the unit shall use to acquire its current time." ["0.openwrt.pool.ntp.org"]
				LocalServer int      `json:"local-server"` // def: true, "Start a NTP server that provides the time to local clients."
			} `json:"ntp,omitempty"`
			Mdns struct { // "This section can be used to configure the MDNS server."
				Enable int `json:"enable"` // def: false, "Enable this option if you would like to enable the MDNS server on the unit."
			} `json:"mdns,omitempty"`
			Rtty struct { // "This section can be used to setup a persistent connection to a rTTY server."
				Host  string `json:"host"`  // "The server that the device shall connect to."
				Port  int    `json:"port"`  // def: 5912, max: 65525, "This option defines the port that device shall connect to."
				Token string `json:"token"` // min: 32, max: 32, "The security token that shall be used to authenticate with the server."
			} `json:"rtty,omitempty"`
			Log struct { // "This section can be used to configure remote syslog support."
				Host  string `json:"host"`  // "IP address of a syslog server to which the log messages should be sent in addition to the local destination."
				Port  int    `json:"port"`  // min: 100, max: 65535, "IP address of a syslog server to which the log messages should be sent in addition to the local destination."
				Proto string `json:"proto"` // "Sets the protocol to use for the connection, either tcp or udp.", ["tcp","udp"],"default": "udp"
				Size  int    `json:"size"`  // min: 32, def: 1000, "Size of the file based log buffer in KiB. This value is used as the fallback value for log_buffer_size if the latter is not specified."
			} `json:"log,omitempty"`
			HTTP struct { // "Enable the webserver with the on-boarding webui"
				HTTPPort int `json:"http-port"` // min: 1, max: 65535, def: 80, "The port that the HTTP server should run on."
			} `json:"http,omitempty"`
			Igmp struct { // "This section allows enabling the IGMP/Multicast proxy"
				Enable int `json:"enable"` // def: false, "This option defines if the IGMP/Multicast proxy shall be enabled on the device."
			} `json:"igmp,omitempty"`
			Ieee8021X struct {
				CaCertificate        string     `json:"ca-certificate,omitempty"`     // "The local servers CA bundle."
				UseLocalCertificates int        `json:"use-local-certificates"`       // def: false, "The device will use its local certificate bundle for the Radius server and ignore all other certificate options in this section."
				ServerCertificate    string     `json:"server-certificate,omitempty"` // "The local servers certificate."
				PrivateKey           string     `json:"private-key,omitempty"`        // "The local servers private key"
				Users                []struct { // "Specifies a collection of local EAP user/psk/vid triplets."
					Mac      string `json:"mac"`       //
					UserName string `json:"user-name"` // min: 1,
					Password string `json:"password"`  // min: 8, max: 63
					VlanID   int    `json:"vlan-id"`   // max: 4096
				} `json:"users,omitempty"`
			} `json:"ieee8021x,omitempty"`
			RadiusProxy struct { // "This section can be used to setup a radius security proxy instance (radsecproxy)."
				Host   string `json:"host"`   // "The remote proxy server that the device shall connect to."
				Port   int    `json:"port"`   // def: 2083, max: 65535, "The remote proxy port that the device shall connect to."
				Secret string `json:"secret"` // "The radius secret that will be used for the connection."
			} `json:"radius-proxy,omitempty""`
			WifiSteering struct {
				Mode          string `json:"mode"`           // "Wifi sterring can happen either locally or via the backend gateway."
				AssocSteering bool   `json:"assoc-steering"` // "Allow rejecting assoc requests for steering purposes."
				//Network           string `json:"network"` // not in schema but pulled from device
				RequiredProbeSnr  int `json:"required-probe-snr"`  // "Minimum required signal level (dBm) for connected clients. If the client will be kicked if the SNR drops below this value."
				RequiredRoamSnr   int `json:"required-roam-snr"`   // "Minimum required signal level (dBm) to allow connections. If the SNR is below this value, probe requests will not be replied to."
				RequiredSnr       int `json:"required-snr"`        // "Minimum required signal level (dBm) before an attempt is made to roam the client to a better AP."
				LoadKickThreshold int `json:"load-kick-threshold"` // "Minimum channel load (%) before kicking clients"
			} `json:"wifi-steering,omitempty"`
		} `json:"services,omitempty""`
		Metrics struct {
			DhcpSnooping struct { // "DHCP snooping allows us to intercept DHCP packages on interface that are bridged, where DHCP is not offered as a service by the AP."
				Filters []string `json:"filters"` // "A list of the message types that shall be sent to the backend." ["ack","discover","offer","request","solicit","reply","renew"]
			} `json:"dhcp-snooping,omitempty"`
			Health struct { // "Health check gets executed periodically and will report a health value between 0-100 indicating how healthy the device thinks it is"
				Interval int `json:"interval"` // min: 60, "The reporting interval defined in seconds."
			} `json:"health,omitempty""`
			Statistics struct { // "Statistics are traffic counters, neighbor tables, ..."
				Interval int      `json:"interval"` // "The reporting interval defined in seconds."
				Types    []string `json:"types"`    // "A list of names of subsystems that shall be reported periodically." ["ssids","lldp","clients"]
			} `json:"statistics,omitempty""`
			WifiFrames struct { // "Define which types of ieee802.11 management frames shall be sent up to the controller."
				Filters []string `json:"filters"` // "A list of the management frames types that shall be sent to the backend." ["probe","auth","assoc","disassoc","deauth","local-deauth","inactive-deauth","key-mismatch","beacon-report","radar-detected"]
			} `json:"wifi-frames,omitempty"`
		} `json:"metrics,omitempty""`
		ConfigRaw [][]string `json:"config-raw,omitempty""` // "This object allows passing raw uci commands, that get applied after all the other configuration was ben generated." [["set","system.@system[0].timezone","GMT0"],["delete","firewall.@zone[0]"],["delete","dhcp.wan"],["add","dhcp","dhcp"],["add-list","system.ntp.server","0.pool.example.org"],["del-list","system.ntp.server","1.openwrt.pool.ntp.org"]]
	} `json:"configuration,omitempty"`
	CreatedTimestamp          int     `json:"createdTimestamp"`
	DevicePassword            string  `json:"devicePassword"`
	DeviceType                string  `json:"deviceType"`
	Firmware                  string  `json:"firmware"`
	FwUpdatePolicy            string  `json:"fwUpdatePolicy"`
	LastConfigurationChange   int     `json:"lastConfigurationChange"`
	LastConfigurationDownload int     `json:"lastConfigurationDownload"`
	LastFWUpdate              int     `json:"lastFWUpdate"`
	Location                  string  `json:"location"`
	MacAddress                string  `json:"macAddress"`
	Manufacturer              string  `json:"manufacturer"`
	Notes                     []*Note `json:"notes"`
	Owner                     string  `json:"owner"`
	SerialNumber              string  `json:"serialNumber"`
	Venue                     string  `json:"venue"`
}

type Note struct {
	Created   int    `json:"created"`
	CreatedBy string `json:"createdBy"`
	Note      string `json:"note"`
}

// ListInfo is a handler which checks the supplied string against the
// DeviceInfo list to return a subset of the requested information.
// An empty string can be supplied to receive a 'generic' description.
func (dev *Device) ListInfo(info string) []string {
	switch {
	case info == "interfaces":
		return dev.GenerateInterfaceReport()
	case info == "configuration":
		//
	case info == "capabilities":
		//
	case info == "status":
		//
	case info == "stats":
		//
	case info == "logs":
		//
	case info == "health":
		//
	default: // == ""
		return dev.GenerateDescription()
	}
	return []string{}
}

// GenerateDescription returns a 'generic' list of important elements of
// the Device object in a printable form.
func (dev *Device) GenerateDescription() (list []string) {
	desc := fmt.Sprintf("Manufacturer: %s, ", dev.Manufacturer)
	desc += fmt.Sprintf("Type: %s, ", dev.DeviceType)
	desc += fmt.Sprintf("MAC Address: %s, ", dev.MacAddress)
	desc += fmt.Sprintf("Firmware: %s, ", dev.Firmware)
	desc += fmt.Sprintf("UUID: %d, ", dev.UUID)
	list = append(list, desc)
	return list
}

// GenerateConfigReport returns a list of Configuration details that are
// critical or interesting.
// [+]
func (dev *Device) GenerateConfigReport() (list []string) {
	desc := fmt.Sprintf("Location: %s, ", dev.Configuration.Unit.Location)
	desc += fmt.Sprintf("Timezone: %s, ", dev.Configuration.Unit.Timezone)
	list = append(list, desc)
	return list
}

// GenerateInterfaceReport returns a list of each Interface's descriptions.
func (dev *Device) GenerateInterfaceReport() (list []string) {
	for n := 0; n < len(dev.Configuration.Interfaces); n++ {
		list = append(list, dev.Configuration.Interfaces[n].GenerateDescription())
	}
	return list
}

// GenerateRadioReport returns a list of each Radio's descriptions.
func (dev *Device) GenerateRadioReport() (list []string) {
	for n := 0; n < len(dev.Configuration.Radios); n++ {
		list = append(list, dev.Configuration.Radios[n].GenerateDescription())
	}
	return list
}
