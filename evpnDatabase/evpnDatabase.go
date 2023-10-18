package evpnDatabase

import (
	"encoding/xml"
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"strings"
)

type EVPNDatabase struct {
	XMLName xml.Name `xml:"rpc-reply"`
	Junos   string   `xml:"xmlns:junos,attr"`
	EvpnDB  struct {
		XMLName         xml.Name `xml:"evpn-database-information"`
		RoutingXMLNS    string   `xml:"xmlns,attr"`
		EvpnDBInstances []struct {
			XMLName      xml.Name `xml:"evpn-database-instance"`
			JunosStyle   string   `xml:"junos:style,attr"`
			InstanceName string   `xml:"instance-name"`
			MacEntries   []struct {
				XMLName          xml.Name `xml:"mac-entry"`
				JunosStyle       string   `xml:"junos:style,attr"`
				VlanID           int      `xml:"vlan-id"`
				MacAddress       string   `xml:"mac-address"`
				ActiveSource     string   `xml:"active-source"`
				ActiveSourceTime string   `xml:"active-source-timestamp"`
				IPAddress        []string `xml:"ip-address,omitempty"`
			} `xml:"mac-entry"`
		} `xml:"evpn-database-instance"`
	} `xml:"evpn-database-information"`
	CLI struct {
		XMLName xml.Name `xml:"cli"`
		Banner  string   `xml:"banner"`
	} `xml:"cli"`
}

func EvpnDatabase(pairs [][]string) {
	for _, pair := range pairs {

		// extract node name from the pair
		path := pair[0] // extract the first element of the pair
		// split the path by backslashes and get the last component (file name)
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		// split the file name by hyphens and get the second component (node name)
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/evpn database:", nodeName)
		pterm.Info.Printf("Comparing EVPN Database between %s and %s\n", pair[1], pair[0])
		file1Content, err := os.Open(pair[1])
		if err != nil {
			// Handle error
			continue
		}
		defer file1Content.Close()
		//Read XML Data from file
		xmlData1, err := io.ReadAll(file1Content)
		if err != nil {
			log.Fatal(err)
		}

		file2Content, err := os.Open(pair[0])
		if err != nil {
			// Handle error
			continue
		}
		defer file2Content.Close()
		//Read XML Data from file
		xmlData2, err := io.ReadAll(file2Content)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the XML data into a Go structs
		var instanceInfo1 EVPNDatabase
		err = xml.Unmarshal(xmlData1, &instanceInfo1)
		if err != nil {
			log.Fatal(err)
		}

		var instanceInfo2 EVPNDatabase
		err = xml.Unmarshal(xmlData2, &instanceInfo2)
		if err != nil {
			log.Fatal(err)
		}

		// Compare the contents of the two files
		var differences []string
		for i, instance1 := range instanceInfo1.EvpnDB.EvpnDBInstances {
			instance2 := instanceInfo2.EvpnDB.EvpnDBInstances[i]
			if instance1.InstanceName == instance2.InstanceName {
				differences = append(differences, pterm.FgLightWhite.Sprintf("EVPN Instance Name: %s", instance1.InstanceName))

				//Define a IP/ESI to Name map
				bpeMap := map[string]string{
					"10.17.2.22":                    "mcld5-bpe-1a",
					"10.17.2.23":                    "mcld5-bpe-1b",
					"10.17.2.24":                    "mcmlw-bpe-1a",
					"10.17.2.25":                    "mcmlw-bpe-1b",
					"00:24:00:00:45:00:00:00:00:00": "LD5-esi-ae45",
					"00:18:00:00:45:00:00:00:00:00": "MLW-esi-ae45",
				}
				// Missing MAC Entry
				for _, macEntry1 := range instance1.MacEntries {
					found := false
					for _, macEntry2 := range instance2.MacEntries {
						if macEntry1.MacAddress == macEntry2.MacAddress && macEntry1.VlanID == macEntry2.VlanID {
							found = true
							break
						}
					}
					if !found {
						//check if the ip or esi in the map
						if bpeName, ok := bpeMap[macEntry1.ActiveSource]; ok {
							macEntry1.ActiveSource = bpeName
						}
						// Append the difference
						differences = append(differences, pterm.FgLightRed.Sprintf("- Missing MAC Entry -> VLAN %d, MAC %s, IP %s, Active Source: %s", macEntry1.VlanID, macEntry1.MacAddress, macEntry1.IPAddress, macEntry1.ActiveSource))
					}
				}
				// Added MAC Entry
				for _, macEntry2 := range instance2.MacEntries {
					found := false
					for _, macEntry1 := range instance1.MacEntries {
						if macEntry2.MacAddress == macEntry1.MacAddress {
							found = true
							break
						}
					}
					if !found {
						if bpeName, ok := bpeMap[macEntry2.ActiveSource]; ok {
							macEntry2.ActiveSource = bpeName
						}
						// Append the difference
						differences = append(differences, pterm.FgLightGreen.Sprintf("+ Added MAC Entry -> VLAN %d, MAC %s, IP %s, Active Source: %s", macEntry2.VlanID, macEntry2.MacAddress, macEntry2.IPAddress, macEntry2.ActiveSource))
					}
				}

				// Active Source Changed
				for _, macEntry1 := range instance1.MacEntries {
					for _, macEntry2 := range instance2.MacEntries {
						if macEntry1.MacAddress == macEntry2.MacAddress && macEntry1.VlanID == macEntry2.VlanID && macEntry1.ActiveSource != macEntry2.ActiveSource {
							if bpeName, ok := bpeMap[macEntry1.ActiveSource]; ok {
								macEntry1.ActiveSource = bpeName
							}
							if bpeName, ok := bpeMap[macEntry2.ActiveSource]; ok {
								macEntry2.ActiveSource = bpeName
							}
							// Append the difference
							differences = append(differences, pterm.FgLightMagenta.Sprintf("*Active Source has changed for VLAN %d -- %s %s FROM %s -> TO %s", macEntry1.VlanID, macEntry1.MacAddress, macEntry1.IPAddress, macEntry1.ActiveSource, macEntry2.ActiveSource))

						}
					}
				}
				// Missing IP
				for _, macEntry1 := range instance1.MacEntries {
					for _, macEntry2 := range instance2.MacEntries {
						if macEntry1.MacAddress == macEntry2.MacAddress && macEntry1.VlanID == macEntry2.VlanID {
							// Compare each IP in the two slices
							for _, ipEntry1 := range macEntry1.IPAddress {
								ipFound := false
								for _, ipEntry2 := range macEntry2.IPAddress {
									if ipEntry1 == ipEntry2 {
										ipFound = true
										break
									}
								}
								if !ipFound {
									differences = append(differences, pterm.FgLightRed.Sprintf("- Missing IP in VLAN %d / MAC %s -> %s", macEntry1.VlanID, macEntry1.MacAddress, ipEntry1))
								}
							}
						}
					}

				}
				// Added IP Addresses
				for _, macEntry2 := range instance2.MacEntries {
					for _, macEntry1 := range instance1.MacEntries {
						if macEntry2.MacAddress == macEntry1.MacAddress && macEntry2.VlanID == macEntry1.VlanID {
							// Compare each IP in the two slices
							for _, ipEntry2 := range macEntry2.IPAddress {
								ipFound := false
								for _, ipEntry1 := range macEntry1.IPAddress {
									if ipEntry2 == ipEntry1 {
										ipFound = true
										break
									}
								}
								if !ipFound {
									differences = append(differences, pterm.FgLightGreen.Sprintf("+ Added IP in VLAN %d / MAC %s -> %s", macEntry2.VlanID, macEntry2.MacAddress, ipEntry2))
								}
							}
						}
					}

				}
			}
		}

		if len(differences) == 0 {
			pterm.Success.Println("No differences found.\n")
		} else {
			pterm.Warning.Println("Differences found. Please check the following Instances:\n")
			for _, diff := range differences {
				pterm.Println(diff)
			}
		}
	}
}
