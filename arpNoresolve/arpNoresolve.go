package arpNoresolve

// arp xml structure
import "encoding/xml"

type RpcReply struct {
	XMLName             xml.Name `xml:"rpc-reply"`
	Text                string   `xml:",chardata"`
	Junos               string   `xml:"junos,attr"`
	ArpTableInformation struct {
		Text          string `xml:",chardata"`
		Xmlns         string `xml:"xmlns,attr"`
		Style         string `xml:"style,attr"`
		ArpTableEntry []struct {
			Text               string `xml:",chardata"`
			MacAddress         string `xml:"mac-address"`
			IpAddress          string `xml:"ip-address"`
			InterfaceName      string `xml:"interface-name"`
			ArpTableEntryFlags struct {
				Text string `xml:",chardata"`
				None string `xml:"none"`
			} `xml:"arp-table-entry-flags"`
		} `xml:"arp-table-entry"`
		ArpEntryCount string `xml:"arp-entry-count"`
	} `xml:"arp-table-information"`
	Cli struct {
		Text   string `xml:",chardata"`
		Banner string `xml:"banner"`
	} `xml:"cli"`
}
