package tipWifi

type Passpoint struct { // "Enable Hotspot 2.0 support."
	VenueName  []string `json:"venue-name"`  // "This parameter can be used to configure one or more Venue Name Duples for Venue Name ANQP information."
	VenueGroup int      `json:"venue-group"` // max: 32, "The available values are defined in 802.11u."
	VenueType  int      `json:"venue-type"`  // max: 32, "The available values are defined in IEEE Std 802.11u-2011, 7.3.1.34"
	VenueURL   string   `json:"venue-url"`   // "This parameter can be used to configure one or more Venue URL Duples to provide additional information corresponding to Venue Name information."
	AuthType   struct { // "This parameter indicates what type of network authentication is used in the network."
		Type string `json:"type"` // "Specifies the specific network authentication type in use." ["terms-and-conditions","online-enrollment","http-redirection","dns-redirection"]
		URI  string `json:"uri"`  // "Specifies the redirect URL applicable to the indicated authentication type." ["https://operator.example.org/wireless-access/terms-and-conditions.html","http://www.example.com/redirect/me/here/"]
	} `json:"auth-type"`
	DomainName      string     `json:"domain-name"`        // "The IEEE 802.11u Domain Name."
	NaiRealm        []string   `json:"nai-realm"`          // "NAI Realm information"
	Osen            int        `json:"osen"`               // "OSU Server-Only Authenticated L2 Encryption Network;"
	AnqpDomain      int        `json:"anqp-domain"`        // min: 0, max: 65535, "ANQP Domain ID, An identifier for a set of APs in an ESS that share the same common ANQP information."
	Anqp3GppCellNet string     `json:"anqp-3gpp-cell-net"` // "The ANQP 3GPP Cellular Network information."
	FriendlyName    []string   `json:"friendly-name"`      // "This parameter can be used to configure one or more Operator Friendly Name Duples."
	Icon            []struct { // "The operator icons."
		Width    int    `json:"width"`    // "The width of the operator icon in pixel",
		Height   int    `json:"height"`   // "The height of the operator icon in pixel"
		Type     string `json:"type"`     // "The mimetype of the operator icon" ex: image/png
		URI      string `json:"uri"`      // "The URL the operator icon is available at"
		Language string `json:"language"` // "ISO 639-2 language code of the icon" ["eng","fre","ger","ita"]
	} `json:"icon"`
}
