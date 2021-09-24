package uClig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// The Auth object is used to supply credentials to the uc.SEC endpoint to retrieve
// an OAuth2 token
type Auth struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

// The OAuth2 object contains all of the information returned from the uc.SEC endpoint
// when a new token is generated. The AccessToken string is supplied as the "Bearer"
// header in subsequent requests to the UCentral endpoints
type OAuth2 struct {
	AccessToken string `json:"access_token"`
	ACLTemplate struct {
		Delete          bool `json:"Delete"`
		PortalLogin     bool `json:"PortalLogin"`
		Read            bool `json:"Read"`
		ReadWrite       bool `json:"ReadWrite"`
		ReadWriteCreate bool `json:"ReadWriteCreate"`
	} `json:"aclTemplate"`
	Created                int    `json:"created"`
	ErrorCode              int    `json:"errorCode"`
	ExpiresIn              int    `json:"expires_in"`
	IdleTimeout            int    `json:"idle_timeout"`
	RefreshToken           string `json:"refresh_token"`
	TokenType              string `json:"token_type"`
	UserMustChangePassword bool   `json:"userMustChangePassword"`
	Username               string `json:"username"`
}

// DisplayToken is a print function wrapped around the OAuth2 object.
func (o *OAuth2) DisplayToken() {
	fmt.Println(o.AccessToken)
}

// The UCentral object collects the required information to interact with the UCentral services.
type UCentral struct {
	SEC    string
	GW     string
	FMS    string
	Auth   *Auth
	OAuth2 *OAuth2
}

// The Endpoints object contains a list of the Endpoint object.
type Endpoints struct {
	Entry []*Endpoint `json:"endpoints"`
}

// The Endpoint object provides information about ways to access the UCentral services.
type Endpoint struct {
	AuthenticationType string `json:"authenticationType"`
	ID                 int64  `json:"id"`
	Type               string `json:"type"`
	URI                string `json:"uri"`
	Vendor             string `json:"vendor"`
}

// DisplayEntries is a print function wrapped around the Endpoints object.
func (ep *Endpoints) DisplayEntries() {
	for i := 0; i < len(ep.Entry); i++ {
		fmt.Printf("[%s] %s\n", ep.Entry[i].Type, ep.Entry[i].URI)
	}
}

// Login retrieves an OAuth2 token from the uc.SEC endpoint
// which can be used for access to the other endpoints.
func (uc *UCentral) Login() error {
	if uc.Auth.UserID == "" || uc.Auth.Password == "" {
		return errors.New("Missing Credentials")
	}
	if uc.SEC == "" {
		return errors.New("Missing Security Endpoint")
	}
	jsonAuth, err := json.Marshal(&uc.Auth)
	if err != nil {
		return err
	}
	resp, err := PostRequest(nil, uc.SEC, "oauth2", jsonAuth)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}
	return json.Unmarshal(body, &uc.OAuth2)
}

// Logout deletes the OAuth2 token from the uc.SEC endpoint
// and should be called after every completed session.
func (uc *UCentral) Logout() error {
	endpoint := fmt.Sprintf("oauth2/%s", uc.OAuth2.AccessToken)
	resp, err := DeleteRequest(uc.OAuth2, uc.SEC, endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}
	if resp.Status != "204 OK" {
		return errors.New(resp.Status)
	}
	return nil
}

// PopulateEndpoints asks the uc.SEC endpoint for other endpoints such as
// the GW and FMS, and populates the UCentral structure with them.
func (uc *UCentral) PopulateEndpoints() error {
	if uc.SEC == "" || uc.OAuth2.AccessToken == "" {
		return errors.New("Must authenticate first")
	}
	resp, err := GetRequest(uc.OAuth2, uc.SEC, "systemEndpoints")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	ep := &Endpoints{}
	err = json.Unmarshal(body, &ep)

	for i := 0; i < len(ep.Entry); i++ {
		switch {
		case strings.Contains(ep.Entry[i].Type, "gw"):
			tmp := strings.Split(ep.Entry[i].URI, "//")
			uc.GW = tmp[1]
		case strings.Contains(ep.Entry[i].Type, "fms"):
			tmp := strings.Split(ep.Entry[i].URI, "//")
			uc.FMS = tmp[1]
		default:
			fmt.Printf("%s :: %s\n", ep.Entry[i].Type, ep.Entry[i].URI)
		}
	}
	if uc.GW == "" || uc.FMS == "" {
		return errors.New("Did Not Find Desired Endpoints")
	} else {
		return nil
	}
}

// ListDevices returns the Devices object which is a list of the Device object.
// The Device object from the GW contains the complete Configuration and other
// detailed information on the device itself.
func (uc *UCentral) ListDevices() (*Devices, error) {
	if uc.GW == "" || uc.OAuth2.AccessToken == "" {
		return nil, errors.New("Must authenticate first")
	}
	resp, err := GetRequest(uc.OAuth2, uc.GW, "devices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	devs := &Devices{}
	err = json.Unmarshal(body, &devs)
	if err != nil {
		return nil, err
	}

	return devs, nil
}

// GetAllFirmwareDevices returns the FirmwareDevices object which is a list of
// the FirmwareDevice object which provides Status and FW tracking.
func (uc *UCentral) GetAllFirmwareDevices() (*FirmwareDevices, error) {
	resp, err := GetRequest(uc.OAuth2, uc.FMS, "connectedDevices")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	fwds := &FirmwareDevices{}
	err = json.Unmarshal(body, &fwds)
	if err != nil {
		return nil, err
	}

	return fwds, nil
}

// ListFirmwareDevices returns a list of valid DeviceTypes
func (uc *UCentral) ListFirmwareDeviceTypes() (list []string, err error) {
	resp, err := GetRequest(uc.OAuth2, uc.FMS, "firmwares?deviceSet=true")
	if err != nil {
		return list, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	err = json.Unmarshal(body, &list)
	return list, err
}

// GetDevice returns the Device object from the GW which includes the complete Configuration.
func (uc *UCentral) GetDevice(sn string) (*Device, error) {
	resp, err := GetRequest(uc.OAuth2, uc.GW, fmt.Sprintf("device/%s", sn))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	dev := &Device{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		return nil, err
	}
	return dev, nil
}

// GetFirmwareDevice returns the FirmwareDevice object when queried with a registered Serial number.
// This is a different data model than the general uc.GW Device object (devices.go),
// coming from the uc.SEC service to provides Status and FW tracking.
func (uc *UCentral) GetFirmwareDevice(sn string) (*FirmwareDevice, error) {
	fwds, err := uc.GetAllFirmwareDevices()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(fwds.Entry); i++ {
		if sn == fwds.Entry[i].SerialNumber {
			return fwds.Entry[i], nil
		}
	}
	return nil, errors.New("SN Not Found")
}

// GetFirmwareListByDevice returns a Firmwares object which contains a list of the Firmware object.
// The Firmware object contains version control information including the download URI.
func (uc *UCentral) GetFirmwareListByDevice(dev string) (*Firmwares, error) {
	resp, err := GetRequest(uc.OAuth2, uc.FMS, fmt.Sprintf("firmwares?deviceType=%s", dev))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	fws := &Firmwares{}
	err = json.Unmarshal(body, &fws)
	if err != nil {
		return nil, err
	}
	return fws, nil
}

// GetLatestFirmwareByDevice queries the FMS version control registry by Device Type.
// A single Firmware object representing the latest available image for the device is returned.
func (uc *UCentral) GetLatestFirmwareByDevice(dev string) (*Firmware, error) {
	resp, err := GetRequest(uc.OAuth2, uc.FMS, fmt.Sprintf("firmwares?latestOnly=true&deviceType=%s", dev))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	fw := &Firmware{}
	err = json.Unmarshal(body, &fw)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

// UpgradeDeviceToLatest takes a FirmwareDevice as input wrapper around the
// UpgradeDeviceFirmware function to control the input variables.
func (uc *UCentral) UpgradeDeviceToLatest(dev *FirmwareDevice) error {
	fw, err := uc.GetLatestFirmwareByDevice(dev.DeviceType)
	if err != nil {
		return err
	}
	// if current fw and latest are the same, don't upgrade
	if dev.Revision == fw.Revision {
		return errors.New("Latest Revision same as Current Version!")
	}
	fmt.Println(dev.SerialNumber, fw.URI)
	return uc.UpgradeDeviceFirmware(dev.SerialNumber, fw.URI)
}

// UpgradeDeviceFirmware takes a SerialNumber and URI (link to get new fw) as input
// and applies the upgrade to the device, returning any error.
func (uc *UCentral) UpgradeDeviceFirmware(sn, uri string) error {
	upg := &Upgrade{
		SerialNumber: sn,
		URI:          uri,
	}
	jsonData, err := json.Marshal(&upg)
	if err != nil {
		return err
	}
	resp, err := PostRequest(uc.OAuth2, uc.GW, fmt.Sprintf("device/%s/upgrade", sn), jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}

	return json.Unmarshal(body, &uc.OAuth2)
}

// RebootDevice takes a SerialNumber as input and reboots the device
func (uc *UCentral) RebootDevice(sn string) error {
	r := &Reboot{
		SerialNumber: sn,
	}
	jsonData, err := json.Marshal(&r)
	if err != nil {
		return err
	}
	resp, err := PostRequest(uc.OAuth2, uc.GW, fmt.Sprintf("device/%s/reboot", sn), jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}
	if resp.Status != "200 OK" {
		return errors.New(resp.Status)
	}
	return nil
}

// Factory reset takes a SerialNumber and bool as input and factory resets the device.
// The bool's effect is whether the "Redirector", that is the established means by
// which the device connects to UCentral, is kept or also wiped.
func (uc *UCentral) FactoryResetDevice(sn string, keepRedirector bool) error {
	f := &Factory{
		SerialNumber:   sn,
		KeepRedirector: keepRedirector,
	}
	jsonData, err := json.Marshal(&f)
	if err != nil {
		return err
	}
	resp, err := PostRequest(uc.OAuth2, uc.GW, fmt.Sprintf("device/%s/factory", sn), jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}
	if resp.Status != "200 OK" {
		return errors.New(resp.Status)
	}
	return nil
}

// AddNoteToDevice access a SerialNumber and slice of strings as input, and applied
// each line of the slice to the notes section of the device.
func (uc *UCentral) AddNotesToDevice(sn string, notes []string) error {
	n := &Notes{
		SerialNumber: sn,
	}
	t := time.Now().Unix()
	for _, note := range notes {
		no := &Note{
			Created:   int(t),
			CreatedBy: "LindsayBB",
			Note:      note,
		}
		n.Notes = append(n.Notes, no)
	}
	jsonData, err := json.Marshal(&n)
	if err != nil {
		return err
	}
	resp, err := PutRequest(uc.OAuth2, uc.GW, fmt.Sprintf("device/%s", sn), jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if Debug {
		fmt.Printf("|+| %s |+|\n", resp.Status)
	}
	if resp.Status != "200 OK" {
		return errors.New(resp.Status)
	}
	return nil
}
