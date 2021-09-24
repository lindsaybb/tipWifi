package tipWifi

import (
	"fmt"
)

// The Radio object represents a subset of the Device configuration related
// to the Radio, whether that is Wi-Fi or Other.
type Radio struct {
	Band         string      `json:"band"`                    // "Specifies the wireless band to configure the radio for. Available radio device phys on the target system are matched by the wireless band given here. If multiple radio phys support the same band, the settings specified here will be applied to all of them." ["2G","5G","5G-lower","5G-upper","6G"]
	Bandwidth    int         `json:"bandwidth,omitempty"`     // "Specifies a narrow channel width in MHz, possible values are 5, 10, 20."
	Channel      interface{} `json:"channel"`                 // "Specifies the wireless channel to use. A value of 'auto' starts the ACS algorithm."
	Country      string      `json:"country,omitempty"`       // min: 2, max: 2, "Specifies the country code, affects the available channels and transmission powers."
	ChannelMode  string      `json:"channel-mode,omitempty"`  // "Define the ideal channel mode that the radio shall use. This can be 802.11n, 802.11ac or 802.11ax. This is just a hint for the AP. If the requested value is not supported then the AP will use the highest common denominator." ["HT","VHT","HE"],"default": "HE"
	ChannelWidth int         `json:"channel-width,omitempty"` // "The channel width that the radio shall use. This is just a hint for the AP. If the requested value is not supported then the AP will use the highest common denominator." [20,40,80,160,8080], "default": 80
	RequireMode  string      `json:"require-mode,omitempty"`  // "Stations that do no fulfill these HT modes will be rejected." ["HT","VHT","HE"]
	Mimo         string      `json:"mimo,omitempty"`          // "This option allows configuring the antenna pairs that shall be used. This is just a hint for the AP. If the requested value is not supported then the AP will use the highest common denominator." ["1x1","2x2","3x3","4x4","5x5","6x6","7x7","8x8"]
	TxPower      int         `json:"tx-power,omitempty"`      // min: 0, max: 30, "This option specifies the transmission power in dBm"
	Rates        struct {
		Beacon    int `json:"beacon"`
		Multicast int `json:"multicast"`
	} `json:"rates,omitempty"`
	LegacyRates    int `json:"legacy-rates,omitempty"`    // "Allow legacy 802.11b data rates." def: false
	BeaconInterval int `json:"beacon-interval,omitempty"` // min: 15, max: 65535, def: 100, "Beacon interval in kus (1.024 ms)."
	DtimPeriod     int `json:"dtim-period,omitempty"`     // min: 1, max: 255, def: 2, "Set the DTIM (delivery traffic information message) period. There will be one DTIM per this many beacon frames. This may be set between 1 and 255. This option only has an effect on ap wifi-ifaces."
	MaximumClients int `json:"maximum-clients,omitempty"` // "Set the maximum number of clients that may connect to this radio. This value is accumulative for all attached VAP interfaces."
	HeSettings     struct {
		MultipleBssid int `json:"multiple-bssid"` // "Enabling this option will make the PHY broadcast its BSSs using the multiple BSSID beacon IE." def: false
		Ema           int `json:"ema"`            // "Enableing this option will make the PHY broadcast its multiple BSSID beacons using EMA." def: false
		BssColor      int `json:"bss-color"`      // def: 64, "This enables BSS Coloring on the PHY. setting it to 0 disables the feature 1-63 sets the color and 64 will make hostapd pick a random color."
	} `json:"he-settings,omitempty"`
	HostapdIfaceRaw []string `json:"hostapd-iface-raw,omitempty"` // "This array allows passing raw hostapd.conf lines." ["ap_table_expiration_time=3600","device_type=6-0050F204-1","ieee80211h=1","rssi_ignore_probe_request=-75","time_zone=EST5","uuid=12345678-9abc-def0-1234-56789abcdef0","venue_url=1:http://www.example.com/info-eng","wpa_deny_ptk0_rekey=0"]
}

// GenerateDescription returns a string of concatenated values describing the Radio object.
func (r *Radio) GenerateDescription() string {
	desc := fmt.Sprintf("Band: %s, ", r.Band)
	desc += fmt.Sprintf("Channel: %v, ", r.Channel)
	desc += fmt.Sprintf("Width: %d, ", r.ChannelWidth)
	desc += fmt.Sprintf("Mode: %s, ", r.ChannelMode)

	return desc
}
