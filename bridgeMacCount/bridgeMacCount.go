package bridgeMacCount

import "encoding/xml"

type RpcReply struct {
	XMLName          xml.Name `xml:"rpc-reply"`
	Text             string   `xml:",chardata"`
	Junos            string   `xml:"junos,attr"`
	L2aldRtbMacCount struct {
		Text                  string `xml:",chardata"`
		L2aldRtbMacCountEntry []struct {
			Text               string `xml:",chardata"`
			RtbMacCount        string `xml:"rtb-mac-count"`
			RtbName            string `xml:"rtb-name"`
			BdName             string `xml:"bd-name"`
			L2aldRtbIfMacCount struct {
				Text                    string `xml:",chardata"`
				L2aldRtbIfMacCountEntry []struct {
					Text          string `xml:",chardata"`
					InterfaceName string `xml:"interface-name"`
					MacCount      string `xml:"mac-count"`
				} `xml:"l2ald-rtb-if-mac-count-entry"`
			} `xml:"l2ald-rtb-if-mac-count"`
			L2aldRtbLearnVlanMacCount struct {
				Text                           string `xml:",chardata"`
				L2aldRtbLearnVlanMacCountEntry struct {
					Text           string `xml:",chardata"`
					LearnVlan      string `xml:"learn-vlan"`
					MacCount       string `xml:"mac-count"`
					StaticMacCount string `xml:"static-mac-count"`
				} `xml:"l2ald-rtb-learn-vlan-mac-count-entry"`
			} `xml:"l2ald-rtb-learn-vlan-mac-count"`
		} `xml:"l2ald-rtb-mac-count-entry"`
	} `xml:"l2ald-rtb-mac-count"`
	Cli struct {
		Text   string `xml:",chardata"`
		Banner string `xml:"banner"`
	} `xml:"cli"`
}
