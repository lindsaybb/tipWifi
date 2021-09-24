package uClig

import (
	"fmt"
)

// The Ssid object represents a subset of the Device configuration related
// to the SSID configuration for Wi-Fi.
type Ssid struct {
	Purpose           string   `json:"purpose"`                      // ["user-defined","onboarding-ap","onboarding-sta"]
	Name              string   `json:"name"`                         // maxLength: 32, minLength: 1
	WifiBands         []string `json:"wifi-bands"`                   // ["2G","5G","5G-lower","5G-upper","6G"]
	BssMode           string   `json:"bss-mode"`                     // ["ap","sta","mesh","wds-ap","wds-sta","wds-repeater"], "default": "ap"
	Bssid             string   `json:"bssid,omitempty"`              // "Override the BSSID of the network, only applicable in adhoc or sta mode."
	HiddenSsid        int      `json:"hidden-ssid,omitempty"`        // "Disables the broadcasting of beacon frames if set to 1 and,in doing so, hides the ESSID."
	IsolateClients    int      `json:"isolate-clients,omitempty"`    // "Isolates wireless clients from each other on this BSS."
	PowerSave         int      `json:"power-save,omitempty"`         // "Unscheduled Automatic Power Save Delivery."
	RtsThreshold      int      `json:"rts-threshold,omitempty`       // min: 1, max: 65535, "Set the RTS/CTS threshold of the BSS."
	BroadcastTime     int      `json:"broadcast-time,omitempty"`     // "This option will make the unit broadcast the time inside its beacons."
	UnicastConversion int      `json:"unicast-conversion,omitempty"` // "Convert multicast traffic to unicast on this BSS."
	Services          []string `json:"services,omitempty"`           // "The services that shall be offered on this logical interface. These are just strings such as \"wifi-steering\""
	MaximumClients    int      `json:"maximum-clients,omitempty"`    // "Set the maximum number of clients that may connect to this VAP."
	ProxyArp          int      `json:"proxy-arp,omitempty"`          // "Proxy ARP is the technique in which the host router, answers ARP requests intended for another machine."
	VendorElements    string   `json:"vendor-elements,omitempty"`    // "This option allows embedding custom vendor specific IEs inside the beacons of a BSS in AP mode."
	Encryption        struct {
		Proto      string `json:"proto"`      // "The wireless encryption protocol that shall be used for this BSS", ["none","psk","psk2","psk-mixed","wpa","wpa2","wpa-mixed","sae","sae-mixed","wpa3","wpa3-mixed"],
		Key        string `json:"key"`        // min:8, max: 63, "The Pre Shared Key (PSK) that is used for encryption on the BSS when using any of the WPA-PSK modes."
		Ieee80211W string `json:"ieee80211w"` // "Enable 802.11w Management Frame Protection (MFP) for this BSS." ["disabled","optional","required"]
	} `json:"encryption,omitempty"`
	MultiPsk struct { // "A SSID can have multiple PSK/VID mappings. Each one of them can be bound to a specific MAC or be a wildcard."
		Mac    string `json:"mac"`     //
		Key    string `json:"key"`     // min: 8, max: 63, "The Pre Shared Key (PSK) that is used for encryption on the BSS when using any of the WPA-PSK modes."
		VlanID int    `json:"vlan-id"` // max: 4096
	} `json:"multi-psk,omitempty"`
	Rrm struct { // "Enable 802.11k Radio Resource Management (RRM) for this BSS."
		NeighborReporting int    `json:"neighbor-reporting"` // "Enable neighbor report via radio measurements (802.11k)."
		Lci               string `json:"lci"`                // "The content of a LCI measurement subelement"
		CivicLocation     string `json:"civic-location"`     // "The content of a location civic measurement subelement"
		FtmResponder      int    `json:"ftm-responder"`      // "Publish fine timing measurement (FTM) responder functionality on this BSS."
		StationaryAp      int    `json:"stationary-ap"`      // "Stationary AP config indicates that the AP doesn't move."
	} `json:"rrm,omitempty"`
	Rates struct { // "The rate configuration of this BSS."
		Beacon    int `json:"beacon"`    // "The beacon rate that shall be used by the BSS. Values are in Mbps.", [0,1000,2000,5500,6000,9000,11000,12000,18000,24000,36000,48000,54000]
		Multicast int `json:"multicast"` // "The multicast rate that shall be used by the BSS. Values are in Mbps." [0,1000,2000,5500,6000,9000,11000,12000,18000,24000,36000,48000,54000]
	} `json:"rates,omitempty"`
	RateLimit struct { // "The UE rate-limiting configuration of this BSS."
		IngressRate int `json:"ingress-rate"` // "The ingress rate to which hosts will be shaped. Values are in Mbps"
		EgressRate  int `json:"egress-rate"`  // "The egress rate to which hosts will be shaped. Values are in Mbps"
	} `json:"rate-limit,omitempty"`
	Roaming struct { // "Enable 802.11r Fast Roaming for this BSS."
		MessageExchange  string `json:"message-exchange"`  // "Shall the pre authenticated message exchange happen over the air or distribution system." ["air","ds"],"default": "ds"
		GeneratePsk      int    `json:"generate-psk"`      // "Whether to generate FT response locally for PSK networks. This avoids use of PMK-R1 push/pull from other APs with FT-PSK networks." def: true (1)
		DomainIdentifier string `json:"domain-identifier"` // min:4, max:4, "Mobility Domain identifier (dot11FTMobilityDomainID, MDID)."
		PmkR0KeyHolder   string `json:"pmk-r0-key-holder"` // "The pairwise master key R0. This is unique to the mobility domain and is required for fast roaming over the air. If the field is left empty a deterministic key is generated."
		PmkR1KeyHolder   string `json:"pmk-r1-key-holder"` // "The pairwise master key R1. This is unique to the mobility domain and is required for fast roaming over the air. If the field is left empty a deterministic key is generated."
	} `json:"roaming,omitempty"`
	Radius       *Radius `json:"radius,omitempty"`
	Certificates struct {
		UseLocalCertificates int    `json:"use-local-certificates"` // "The device will use its local certificate bundle for the TLS setup and ignores all other certificate options in this section." def: false
		CaCertificate        string `json:"ca-certificate"`         // "The local servers CA bundle."
		Certificate          string `json:"certificate"`            // "The local servers certificate."
		PrivateKey           string `json:"private-key"`            // "The local servers private key"
		PrivateKeyPassword   string `json:"private-key-password"`   // "The password required to read the private key."
	} `json:"certificates,omitempty"`
	PassPoint         *Passpoint `json:"pass-point,omitempty"`
	QualityThresholds struct {
		ProbeRequestRssi       int `json:"probe-request-rssi"`       // "Probe requests will be ignored if the rssi is below this threshold.
		AssociationRequestRssi int `json:"association-request-rssi"` // "Association requests will be denied if the rssi is below this threshold."
	} `json:"quality-thresholds,omitempty"`
	HostapdBssRaw string `json:"hostapd-bss-raw,omitempty"` // "This array allows passing raw hostapd.conf lines." ["ap_table_expiration_time=3600","device_type=6-0050F204-1","ieee80211h=1","rssi_ignore_probe_request=-75","time_zone=EST5","uuid=12345678-9abc-def0-1234-56789abcdef0","venue_url=1:http://www.example.com/info-eng","wpa_deny_ptk0_rekey=0"]
}

// GenerateDescription returns a string of concatenated values describing the Ssid object.
func (s *Ssid) GenerateDescription() string {
	desc := fmt.Sprintf("Name: %s, ", s.Name)
	desc += fmt.Sprintf("Mode: %s, ", s.BssMode)
	desc += fmt.Sprintf("Bands: ")
	if len(s.WifiBands) < 1 {
		desc += fmt.Sprintf(", ")
	} else {
		for i := 0; i < len(s.WifiBands); i++ {
			desc += fmt.Sprintf("%s, ", s.WifiBands[i])
		}
	}
	desc += fmt.Sprintf("Services: ")
	if len(s.Services) < 1 {
		desc += fmt.Sprintf(", ")
	} else {
		for i := 0; i < len(s.Services); i++ {
			desc += fmt.Sprintf("%s, ", s.Services[i])
		}
	}
	if s.HiddenSsid == 1 {
		desc += fmt.Sprintf("Hidden, ")
	}
	if s.IsolateClients == 1 {
		desc += fmt.Sprintf("Client Isolation, ")
	}
	desc += fmt.Sprintf("Encyption: %s, ", s.Encryption.Proto)
	desc += fmt.Sprintf("Probe Req RSSI: %d, Assoc Req RSSI: %d, ", s.QualityThresholds.ProbeRequestRssi, s.QualityThresholds.AssociationRequestRssi)

	return desc
}
