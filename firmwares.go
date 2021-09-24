package uClig

import (
	"fmt"
)

// The Firmwares object contains a list of the Firmware object.
type Firmwares struct {
	Entry []*Firmware `json:"firmwares"`
}

// GenerateList returns a list of each Firmware entry's Description.
func (fws *Firmwares) GenerateList() (list []string) {
	for i := 0; i < len(fws.Entry); i++ {
		desc := fws.Entry[i].GenerateDescription()
		list = append(list, desc)
	}
	return list
}

// The DeviceTypes list is a static record of the applicable device types
// as returned from the uc.FMS 'firmwares?deviceSet=true' endpoint.
var DeviceTypes = []string{
	"cig_wf160d",
	"cig_wf188",
	"cig_wf194c",
	"edgecore_eap101",
	"edgecore_eap102",
	"edgecore_ecs4100-12ph",
	"edgecore_ecw5211",
	"edgecore_ecw5410",
	"edgecore_oap100",
	"edgecore_spw2ac1200",
	"edgecore_ssw2ac2600",
	"hfcl_ion4.yml",
	"indio_um-305ac",
	"linksys_e8450-ubi",
	"linksys_ea6350",
	"linksys_ea8300",
	"mikrotik_nand",
	"mikrotik_nand-large",
	"tplink_cpe210_v3",
	"tplink_cpe510_v3",
	"tplink_eap225_outdoor_v1",
	"tplink_ec420",
	"tplink_ex227",
	"tplink_ex228",
	"tplink_ex447",
	"wallys_dr40x9",
}

// The Firmware object contains version control information including the download URI.
type Firmware struct {
	Created       int           `json:"created"`
	Description   string        `json:"description"`
	DeviceType    string        `json:"deviceType"`
	Digest        string        `json:"digest"`
	DownloadCount int           `json:"downloadCount"`
	FirmwareHash  string        `json:"firmwareHash"`
	ID            string        `json:"id"`
	Image         string        `json:"image"`
	ImageDate     int           `json:"imageDate"`
	Latest        bool          `json:"latest"`
	Location      string        `json:"location"`
	Notes         []interface{} `json:"notes"`
	Owner         string        `json:"owner"`
	Release       string        `json:"release"`
	Revision      string        `json:"revision"`
	Size          int           `json:"size"`
	Uploader      string        `json:"uploader"`
	URI           string        `json:"uri"`
}

// GenerateDescription returns a string of concatenated values describing the Firmware object.
func (fw *Firmware) GenerateDescription() string {
	desc := fmt.Sprintf("ID: %s", fw.ID)
	desc += fmt.Sprintf("Release: %s, ", fw.Release)
	desc += fmt.Sprintf("Revision: %s, ", fw.Revision)
	desc += fmt.Sprintf("Image Date: %d, Created: %d, ", fw.ImageDate, fw.Created)
	desc += fmt.Sprintf("URI: %s, ", fw.URI)

	return desc
}

// The Upgrade object represents the minimal information to upgrade a device.
type Upgrade struct {
	SerialNumber string `json:"serialNumber"`
	URI          string `json:"uri"`
}

// The Reboot object is used to marshal the reboot data
type Reboot struct {
	SerialNumber string `json:"serialNumber"`
}

// The Factory object is used to marshal the factory default data
type Factory struct {
	SerialNumber   string `json:"serialNumber"`
	KeepRedirector bool   `json:"keepRedirector"`
}

// The Notes object is used to marshal the notes data
type Notes struct {
	SerialNumber string  `json:"serialNumber"`
	Notes        []*Note `json:"notes"`
}

// The FirmwareDevice object represents the uc.FMS "Device" object and is used
// for checking status and firmware upgrade activities.
type FirmwareDevice struct {
	DeviceType   string `json:"deviceType"`
	EndPoint     string `json:"endPoint"`
	LastUpdate   int    `json:"lastUpdate"`
	Revision     string `json:"revision"`
	SerialNumber string `json:"serialNumber"`
	Status       string `json:"status"`
}

// IsConnected returns a bool interpreting the status of the FirmwareDevice.
func (fwd *FirmwareDevice) IsConnected() bool {
	return fwd.Status == "connected"
}
