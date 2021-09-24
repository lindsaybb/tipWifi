package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/lindsaybb/tipWifi"
)

var (
	userFlag   = flag.String("un", "tip@ucentral.com", "uCentral Username")
	passFlag   = flag.String("pw", "openwifi", "uCentral Password")
	secUrlFlag = flag.String("sec", "lindsay.arilia.com:18001", "uCentral Security Endpoint")
	helpFlag   = flag.Bool("h", false, "Show this help")
)

var validArgs = []string{
	"listdevices",
	"getdevice",
	"getfirmware",
	"upgradefirmware",
	"reboot",
	"annotate",
	"factory",
}

var argFields = map[string][]string{
	validArgs[0]: []string{"<Info Type>"},
	validArgs[1]: []string{"Serial Number"},
	validArgs[2]: []string{"Device Type", "<all>"},
	validArgs[3]: []string{"Serial Number", "<url>"},
	validArgs[4]: []string{"Serial Number"},
	validArgs[5]: []string{"Serial Number", "[comma-separated notes in quotes]"},
	validArgs[6]: []string{"Serial Number", "<supply 'false' if don't want to keep Redirector>"},
}

/*
func usage() {
	fmt.Println(`Program Description:
	Call this program will appropriate arguments`)
	for _, v := range validArgs {
		fmt.Printf("%s: ", v)
		for _, o := range argFields[v] {
			fmt.Printf("%s ", o)
		}
		fmt.Println()
	}
}
*/
func init() {
	flag.Parse()
	if *helpFlag || flag.NArg() < 1 {
		//usage()
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	uc := &tipWifi.UCentral{
		SEC: *secUrlFlag,
		Auth: &tipWifi.Auth{
			UserID:   *userFlag,
			Password: *passFlag,
		},
	}
	err := uc.Login()
	//uc.OAuth2, err := uClig.LoginUCentral(un, pw, secUrl)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Logged In to uCentral")
	}
	defer logout(uc)
	//uc.OAuth2.DisplayToken()
	err = uc.PopulateEndpoints()
	if err != nil {
		log.Println(err)
	}

	// declaring them here to persist between loops... necessary?
	var devs *tipWifi.Devices
	var dev *tipWifi.Device

	var skip bool
	var skip2 bool
	for n := range flag.Args() {
		// Basic flow control allows for multi-command chains
		// and for variables to be supplied with commands.
		if skip2 {
			skip = true
			skip2 = false
			continue
		}
		if skip {
			skip = false
			continue
		}
		arg := parseArg(flag.Args()[n])
		switch arg {
		case 0:
			// "listdevices",
			var info string
			if flag.NArg() >= (n + 2) {
				next := strings.ToLower(flag.Args()[n+1])
				if existsInList(next, tipWifi.DeviceInfo) {
					info = next
					skip = true
				}
			}
			if devs == nil {
				devs, err = uc.ListDevices()
				if err != nil {
					log.Println(err)
					continue
				}
			}
			devs.GenerateDeviceReport(info)

		case 1:
			// "getdevice",
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[1], ":Must supply Device SN")
			}
			skip = true
			sn := strings.ToLower(flag.Args()[n+1])
			if len(sn) != 12 {
				log.Fatalln(sn, ":Incorrect Device SN Length")
			}
			var info string
			if flag.NArg() >= (n + 3) {
				next := strings.ToLower(flag.Args()[n+2])
				if existsInList(next, tipWifi.DeviceInfo) {
					info = next
					skip2 = true
				}
			}
			if dev == nil || dev.SerialNumber != sn {
				dev, err = uc.GetDevice(sn)
				if err != nil {
					log.Println(err)
				}
			}
			tipWifi.DisplayList(dev.SerialNumber, dev.ListInfo(info))

		case 2:
			// "getfirmware",
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[2], ":Must supply Device Type")
			}
			skip = true
			devType := strings.ToLower(flag.Args()[n+1])
			if !existsInList(devType, tipWifi.DeviceTypes) {
				log.Println("Retrieving Valid Device Types...")
				validDevices, err := uc.ListFirmwareDeviceTypes()
				if err != nil {
					log.Fatalln(err)
				}
				for _, v := range validDevices {
					log.Printf("\t%s\n", v)
				}
				log.Fatalln(devType, ":Invalid Device Type Supplied")
			}
			if flag.NArg() >= (n + 3) {
				latest := strings.ToLower(flag.Args()[n+2])
				if latest == "latest" {
					skip2 = true
					fw, err := uc.GetLatestFirmwareByDevice(devType)
					if err != nil {
						log.Fatalln(err)
					}
					tipWifi.DisplayList(devType, []string{fw.GenerateDescription()})
					continue
				}
			}
			fws, err := uc.GetFirmwareListByDevice(devType)
			if err != nil {
				log.Fatalln(err)
			}
			tipWifi.DisplayList(devType, fws.GenerateList())

		case 3:
			// "upgradefirmware",
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[3], ":Must supply Device SN")
			}
			skip = true
			sn := strings.ToLower(flag.Args()[n+1])
			if len(sn) != 12 {
				log.Fatalln(sn, ":Incorrect Device SN Length")
			}
			fwd, err := uc.GetFirmwareDevice(sn)
			if err != nil {
				log.Println(err)
				continue
			}
			if flag.NArg() >= (n + 3) {
				if strings.HasPrefix(flag.Args()[n+2], "http") {
					skip2 = true
					// optional arg of preferred URI supplied for upgrade
					err = uc.UpgradeDeviceFirmware(fwd.SerialNumber, flag.Args()[n+2])
					if err != nil {
						log.Println(err)
					}
					// if upgrade is unsuccessful, could roll back to latest
					// but rather exit and have the user call it again
					continue
				}
			}
			err = uc.UpgradeDeviceToLatest(fwd)
			if err != nil {
				log.Println(err)
			}
		case 4:
			// "reboot"
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[4], ":Must supply Device SN")
			}
			skip = true
			sn := strings.ToLower(flag.Args()[n+1])
			if len(sn) != 12 {
				log.Fatalln(sn, ":Incorrect Device SN Length")
			}
			err = uc.RebootDevice(sn)
			if err != nil {
				log.Println(err)
			}
		case 5:
			// "annotate"
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[5], ":Must supply Device SN")
			}
			skip = true
			sn := strings.ToLower(flag.Args()[n+1])
			if len(sn) != 12 {
				log.Fatalln(sn, ":Incorrect Device SN Length")
			}
			if flag.NArg() < (n + 3) {
				log.Fatalln(validArgs[5], ":Missing the Notes!")
			}
			skip2 = true
			notes := strings.Split(flag.Args()[n+2], ",")
			err = uc.AddNotesToDevice(sn, notes)
			if err != nil {
				log.Println(err)
			}
		case 6:
			// "factory"
			if flag.NArg() < (n + 1) {
				log.Fatalln(validArgs[5], ":Must supply Device SN")
			}
			skip = true
			sn := strings.ToLower(flag.Args()[n+1])
			if len(sn) != 12 {
				log.Fatalln(sn, ":Incorrect Device SN Length")
			}
			keep := true
			if flag.NArg() >= (n + 3) {
				if strings.ToLower(flag.Args()[n+2]) == "false" {
					keep = false
				}
			}
			err = uc.FactoryResetDevice(sn, keep)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Printf("Unknown arg: %s\n", arg)
		}

	}
	//fmt.Printf("API Gateway: %s\nFirmware Management System: %s\n", uc.GW, uc.FMS)
}

func logout(uc *tipWifi.UCentral) {
	err := uc.Logout()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Logged Out of uCentral")
	}
}

func existsInList(s string, l []string) bool {
	for _, v := range l {
		if strings.ToLower(s) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func parseArg(s string) int {
	s = strings.ToLower(s)
	for i := range validArgs {
		if s == validArgs[i] {
			return i
		}
	}
	return -1
}
