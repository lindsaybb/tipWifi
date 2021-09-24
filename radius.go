package uClig

type Radius struct {
	NasIdentifier    string `json:"nas-identifier"`     // "NAS-Identifier string for RADIUS messages. When used, this should be unique to the NAS within the scope of the RADIUS server.""
	ChargeableUserID int    `json:"chargeable-user-id"` // "This will enable support for Chargeable-User-Identity (RFC 4372)." def: false
	Local            struct {
		ServerIdentity string `json:"server-identity"` // "EAP methods that provide mechanism for authenticated server identity delivery use this value." default: uCentral
		Users          []struct {
			Mac      string `json:"mac"`
			UserName string `json:"user-name"` // min: 1
			Password string `json:"password"`  // min: 8, max: 63
			VlanID   int    `json:"vlan-id"`   // max: 4096
		} `json:"users,omitempty"`
	} `json:"local,omitempty"`
	Authentication struct {
		Host             string     `json:"host"`   // "The URI of our Radius server."
		Port             int        `json:"port"`   // "The network port of our Radius server." def: 1812
		Secret           string     `json:"secret"` // "The shared Radius authentication secret."
		RequestAttribute []struct { // [{"id": 27,"value": 900},{"id": 32,"value": "My NAS ID"},{"id": 56,"value": 1004},{"id": 126,"value": "Example Operator"}]
			ID    int         `json:"id"`
			Value interface{} `json:"value"`
		} `json:"request-attribute,omitempty"`
	} `json:"authentication,omitempty`
	Accounting struct {
		Host             string     `json:"host"`   // "The URI of our Radius server."
		Port             int        `json:"port"`   // "The network port of our Radius server." def: 1812
		Secret           string     `json:"secret"` // "The shared Radius authentication secret."
		RequestAttribute []struct { // [{"id": 27,"value": 900},{"id": 32,"value": "My NAS ID"},{"id": 56,"value": 1004},{"id": 126,"value": "Example Operator"}]
			ID    int         `json:"id"`
			Value interface{} `json:"value"`
		} `json:"request-attribute,omitempty"`
		Interval int `json:"interval,omitempty"` // min: 60, max: 600, def:60 "The interim accounting update interval. This value is defined in seconds."
	} `json:"accounting,omitempty"`
}
